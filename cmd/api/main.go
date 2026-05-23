package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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
	dtcCatalogRepo := telemetryrepository.NewPgxDTCatalogRepository(pool)

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
		telemetryusecase.NewListActiveDTCs(dtcRepo, dtcCatalogRepo, vehicleResolver),
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
	docs := handler.NewDocsHandler()
	authConfig, err := middleware.NewAuthConfigFromEnv()
	if err != nil {
		log.Fatal("invalid auth bypass configuration", zap.Error(err))
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger(log))

	r.Get("/docs", docs.Redoc)
	r.Get("/openapi.yaml", docs.OpenAPI)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authConfig))

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
	})

	port := mustEnv("PORT")
	log.Info("starting HTTP server", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
