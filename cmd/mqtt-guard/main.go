package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"torque/cmd/mqtt-guard/handler"
	guardmiddleware "torque/cmd/mqtt-guard/middleware"
	coremiddleware "torque/internal/core/middleware"
	"torque/internal/core/db"
	"torque/internal/core/logger"
	deviceusecase "torque/internal/modules/device/application/usecase"
	devicerepository "torque/internal/modules/device/infrastructure/repository"
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

	deviceRepo := devicerepository.NewGormRepository(conn)

	authHandler := handler.NewAuthHandler(deviceusecase.NewAuthenticateDevice(deviceRepo))
	aclHandler := handler.NewACLHandler(deviceusecase.NewAuthorizeDevice(deviceRepo))

	r := chi.NewRouter()
	r.Use(coremiddleware.Logger(log))
	r.Use(guardmiddleware.SharedSecret(mustEnv("MQTT_GUARD_SECRET")))

	r.Post("/mqtt/auth", authHandler.Auth)
	r.Post("/mqtt/acl", aclHandler.ACL)

	port := mustEnv("PORT")
	log.Info("starting mqtt-guard", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
