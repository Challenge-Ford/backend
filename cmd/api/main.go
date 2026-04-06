package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"torque/cmd/api/handler"
	"torque/cmd/api/middleware"
	"torque/internal/core/db"
	"torque/internal/core/logger"
	"torque/internal/core/pki"
	"torque/internal/infrastructure/adapters"
	deviceusecase "torque/internal/modules/device/application/usecase"
	devicedomain "torque/internal/modules/device/domain"
	devicerepository "torque/internal/modules/device/infrastructure/repository"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	telemetryrepository "torque/internal/modules/telemetry/infrastructure/repository"
	vehicleusecase "torque/internal/modules/vehicle/application/usecase"
	vehicledomain "torque/internal/modules/vehicle/domain"
	vehiclerepository "torque/internal/modules/vehicle/infrastructure/repository"
)

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "missing required environment variable: %s\n", key)
		os.Exit(1)
	}
	return v
}

func main() {
	godotenv.Load()

	log, err := logger.New(os.Getenv("LOG_JSON") == "true")
	if err != nil {
		fmt.Println("failed to init logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	ctx := context.Background()

	pool, err := db.ConnectPgx(ctx, mustEnv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	tsPool, err := db.ConnectPgx(ctx, mustEnv("TIMESERIES_DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to timescaledb", zap.Error(err))
	}
	defer tsPool.Close()

	if err := runMigrations(ctx, pool); err != nil {
		log.Fatal("migration failed", zap.Error(err))
	}

	stepPKI, err := pki.NewStepCAClient(
		mustEnv("STEP_CA_URL"),
		mustEnv("STEP_CA_PROVISIONER"),
		mustEnv("STEP_CA_PROVISIONER_PASSWORD"),
		mustEnv("STEP_CA_ROOT_CERT"),
	)
	if err != nil {
		log.Fatal("failed to init step-ca pki client", zap.Error(err))
	}

	repo := vehiclerepository.NewRepository(pool)
	modelRepo := vehiclerepository.NewModelRepository(pool)
	deviceRepo := devicerepository.NewRepository(pool)
	telemetryRepo := telemetryrepository.NewPgxRepository(tsPool)
	dtcRepo := telemetryrepository.NewPgxDTCRepository(tsPool)

	validate := validator.New()
	validate.RegisterValidation("vin", func(fl validator.FieldLevel) bool {
		return vehicledomain.VIN(fl.Field().String()).Validate() == nil
	})
	validate.RegisterValidation("plate", func(fl validator.FieldLevel) bool {
		return vehicledomain.Plate(fl.Field().String()).Validate() == nil
	})
	validate.RegisterValidation("device_name", func(fl validator.FieldLevel) bool {
		return devicedomain.DeviceName(fl.Field().String()).Validate() == nil
	})

	findVehicle := vehicleusecase.NewFindVehicle(repo)
	existsVehicle := vehicleusecase.NewExistsVehicle(repo)
	findCommissionedByVIN := deviceusecase.NewFindCommissionedByVIN(deviceRepo)
	findDeviceByVehicle := deviceusecase.NewFindDeviceByVehicle(deviceRepo)
	checkActiveDTCs := telemetryusecase.NewCheckActiveDTCs(dtcRepo)

	vehicleResolver := adapters.NewVehicleResolver(findVehicle, existsVehicle)
	deviceResolver := adapters.NewDeviceResolver(findCommissionedByVIN, findDeviceByVehicle)
	telemetryResolver := adapters.NewTelemetryResolver(checkActiveDTCs)

	devices := handler.NewDeviceHandler(
		deviceusecase.NewListDevices(deviceRepo),
		deviceusecase.NewCreateDevice(deviceRepo, stepPKI, validate),
		deviceusecase.NewCommissionDevice(deviceRepo, vehicleResolver, validate),
		deviceusecase.NewDecommissionDevice(deviceRepo),
	)

	telemetry := handler.NewTelemetryHandler(
		telemetryusecase.NewListTelemetry(telemetryRepo, vehicleResolver),
		telemetryusecase.NewListActiveDTCs(dtcRepo, vehicleResolver),
	)

	vehicles := handler.NewVehicleHandler(
		vehicleusecase.NewCreateVehicle(repo, modelRepo, validate),
		vehicleusecase.NewGetVehicle(repo, telemetryResolver),
		vehicleusecase.NewListVehicles(repo, telemetryResolver),
		vehicleusecase.NewUpdateVehicle(repo, modelRepo),
		vehicleusecase.NewDeleteVehicle(repo, deviceResolver),
	)

	vehicleModels := handler.NewVehicleModelHandler(
		vehicleusecase.NewListVehicleModels(modelRepo),
		vehicleusecase.NewListVehicleModelYears(modelRepo),
	)

	health := handler.NewHealthHandler(pool, tsPool)

	r := chi.NewRouter()
	r.Use(middleware.Logger(log))
	r.Use(middleware.Auth)

	r.Get("/healthz", health.Liveness)
	r.Get("/readyz", health.Readiness)

	r.Route("/devices", func(r chi.Router) {
		r.Get("/", devices.List)
		r.Post("/", devices.Create)
		r.Post("/{id}/commission", devices.Commission)
		r.Post("/{id}/decommission", devices.Decommission)
	})

	r.Route("/vehicles", func(r chi.Router) {
		r.Get("/", vehicles.List)
		r.Post("/", vehicles.Create)
		r.Get("/{id}", vehicles.Get)
		r.Patch("/{id}", vehicles.Update)
		r.Delete("/{id}", vehicles.Delete)
		r.Get("/{id}/telemetry", telemetry.ListTelemetry)
		r.Get("/{id}/dtcs", telemetry.ListDTCs)
	})

	r.Route("/vehicle-models", func(r chi.Router) {
		r.Get("/", vehicleModels.List)
		r.Get("/{id}/years", vehicleModels.ListYears)
	})

	port := mustEnv("PORT")
	log.Info("starting HTTP server", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	stmts := []string{
		// Schemas
		`CREATE SCHEMA IF NOT EXISTS vehicle`,
		`CREATE SCHEMA IF NOT EXISTS device`,

		// vehicle_models
		`CREATE TABLE IF NOT EXISTS vehicle.vehicle_models (
			id   UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			type VARCHAR(255) NOT NULL
		)`,

		// vehicle_model_years
		`CREATE TABLE IF NOT EXISTS vehicle.vehicle_model_years (
			id         UUID PRIMARY KEY,
			model_id   UUID NOT NULL REFERENCES vehicle.vehicle_models(id),
			year       INT NOT NULL,
			model_url  TEXT
		)`,

		// vehicle_model_year_colors
		`CREATE TABLE IF NOT EXISTS vehicle.vehicle_model_year_colors (
			id            UUID PRIMARY KEY,
			model_year_id UUID NOT NULL REFERENCES vehicle.vehicle_model_years(id),
			name          VARCHAR(100) NOT NULL,
			hex           VARCHAR(7) NOT NULL
		)`,

		// vehicles
		`CREATE TABLE IF NOT EXISTS vehicle.vehicles (
			id             UUID PRIMARY KEY,
			customer_id    UUID,
			model_year_id  UUID NOT NULL REFERENCES vehicle.vehicle_model_years(id),
			vin            VARCHAR(17) NOT NULL,
			plate          VARCHAR(7) NOT NULL,
			color          VARCHAR(7) NOT NULL,
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_by     UUID NOT NULL,
			updated_by     UUID NOT NULL,
			deleted_at     TIMESTAMPTZ,
			deleted_by     UUID
		)`,

		// devices
		`CREATE TABLE IF NOT EXISTS device.devices (
			id              UUID PRIMARY KEY,
			name            VARCHAR(255) NOT NULL,
			vehicle_id      UUID REFERENCES vehicle.vehicles(id),
			certificate_cn  VARCHAR(255) NOT NULL,
			certificate_sn  VARCHAR(255) NOT NULL,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_by      UUID NOT NULL,
			updated_by      UUID NOT NULL,
			deleted_at      TIMESTAMPTZ,
			deleted_by      UUID
		)`,

		// Indexes
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicle_model_years_unique
			ON vehicle.vehicle_model_years (model_id, year)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicles_vin_active
			ON vehicle.vehicles (vin) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicles_plate_active
			ON vehicle.vehicles (plate) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_name
			ON device.devices (name)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_certificate_cn
			ON device.devices (certificate_cn)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_certificate_sn
			ON device.devices (certificate_sn)`,
	}

	for _, sql := range stmts {
		if _, err := tx.Exec(ctx, sql); err != nil {
			return fmt.Errorf("migration: %s: %w", sql, err)
		}
	}
	return tx.Commit(ctx)
}
