package main

import (
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

	conn, err := db.Connect(mustEnv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close(conn)

	if err := conn.Exec("CREATE SCHEMA IF NOT EXISTS vehicle").Error; err != nil {
		log.Fatal("failed to create schema", zap.Error(err))
	}

	if err := db.Migrate(conn, &vehicledomain.Vehicle{}); err != nil {
		log.Fatal("migration failed", zap.Error(err))
	}

	repo := vehiclerepository.NewGormRepository(conn)

	validate := validator.New()
	validate.RegisterValidation("vin", func(fl validator.FieldLevel) bool {
		return vehicledomain.VIN(fl.Field().String()).Validate() == nil
	})
	validate.RegisterValidation("plate", func(fl validator.FieldLevel) bool {
		return vehicledomain.Plate(fl.Field().String()).Validate() == nil
	})

	vehicles := handler.NewVehicleHandler(
		vehicleusecase.NewCreateVehicle(repo, validate),
		vehicleusecase.NewGetVehicle(repo),
		vehicleusecase.NewListVehicles(repo),
		vehicleusecase.NewUpdateVehicle(repo),
		vehicleusecase.NewDeleteVehicle(repo),
	)

	r := chi.NewRouter()
	r.Use(middleware.Auth)

	r.Route("/vehicles", func(r chi.Router) {
		r.Get("/", vehicles.List)
		r.Post("/", vehicles.Create)
		r.Get("/{id}", vehicles.Get)
		r.Patch("/{id}", vehicles.Update)
		r.Delete("/{id}", vehicles.Delete)
	})

	port := mustEnv("PORT")
	log.Info("starting HTTP server", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
