package main

import (
	"fmt"
	"net"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	vehiclev1 "torque/gen/proto/vehicle/v1"
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

	vehicleHandler := handler.NewVehicleHandler(
		vehicleusecase.NewCreateVehicle(repo, validate),
		vehicleusecase.NewGetVehicle(repo),
		vehicleusecase.NewListVehicles(repo),
		vehicleusecase.NewUpdateVehicle(repo),
		vehicleusecase.NewDeleteVehicle(repo),
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.Auth),
	)
	vehiclev1.RegisterVehicleServiceServer(grpcServer, vehicleHandler)

	if os.Getenv("GRPC_REFLECTION") == "true" {
		reflection.Register(grpcServer)
	}

	lis, err := net.Listen("tcp", ":"+mustEnv("PORT"))
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	log.Info("starting gRPC server", zap.String("port", mustEnv("PORT")))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
