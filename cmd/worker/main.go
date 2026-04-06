package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"torque/internal/core/db"
	"torque/internal/core/logger"
	"torque/internal/infrastructure/adapters"
	"torque/internal/infrastructure/messaging"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	devicerepository "torque/internal/modules/device/infrastructure/repository"
	deviceusecase "torque/internal/modules/device/application/usecase"
	telemetryrepository "torque/internal/modules/telemetry/infrastructure/repository"
)

const (
	queueTelemetry = "torque.telemetry"
	queueDTC       = "torque.dtc"
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

	amqpConn, err := amqp.Dial(mustEnv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal("failed to connect to rabbitmq", zap.Error(err))
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		log.Fatal("failed to open rabbitmq channel", zap.Error(err))
	}
	defer ch.Close()

	for _, q := range []string{queueTelemetry, queueDTC} {
		if _, err := ch.QueueDeclare(q, true, false, false, false, nil); err != nil {
			log.Fatal("failed to declare queue", zap.String("queue", q), zap.Error(err))
		}
	}

	deviceRepo := devicerepository.NewRepository(pool)
	telemetryRepo := telemetryrepository.NewPgxRepository(tsPool)
	dtcRepo := telemetryrepository.NewPgxDTCRepository(tsPool)
	findCommissionedByVIN := deviceusecase.NewFindCommissionedByVIN(deviceRepo)
	findDeviceByVehicle := deviceusecase.NewFindDeviceByVehicle(deviceRepo)
	deviceResolver := adapters.NewDeviceResolver(findCommissionedByVIN, findDeviceByVehicle)

	recordTelemetry := telemetryusecase.NewRecordTelemetry(telemetryRepo, deviceResolver)
	recordDTC := telemetryusecase.NewRecordDTC(dtcRepo, deviceResolver)

	go consume(log, ch, queueTelemetry, func(body []byte) error {
		return handleTelemetry(body, recordTelemetry)
	})

	go consume(log, ch, queueDTC, func(body []byte) error {
		return handleDTC(body, recordDTC)
	})

	log.Info("worker started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("shutting down")
}

func consume(log *zap.Logger, ch *amqp.Channel, queue string, handle func([]byte) error) {
	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to start consumer", zap.String("queue", queue), zap.Error(err))
	}
	for msg := range msgs {
		if err := handle(msg.Body); err != nil {
			log.Error("failed to process message",
				zap.String("queue", queue),
				zap.Error(err),
				zap.ByteString("body", msg.Body),
			)
			msg.Nack(false, false)
		} else {
			msg.Ack(false)
		}
	}
}

func handleTelemetry(body []byte, uc *telemetryusecase.RecordTelemetryUseCase) error {
	var msg messaging.TelemetryMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("unmarshal telemetry: %w", err)
	}
	if msg.VIN == "" {
		return fmt.Errorf("missing vin")
	}
	return uc.Execute(context.Background(), telemetrydto.RecordTelemetryInput{
		Time:           msg.Time,
		VIN:            msg.VIN,
		Lat:            msg.Lat,
		Lng:            msg.Lng,
		Alt:            msg.Alt,
		GPSSpeed:       msg.GPSSpeed,
		Heading:        msg.Heading,
		HDOP:           msg.HDOP,
		RPM:            msg.RPM,
		Speed:          msg.Speed,
		CoolantTemp:    msg.CoolantTemp,
		IntakeTemp:     msg.IntakeTemp,
		EngineLoad:     msg.EngineLoad,
		ThrottlePos:    msg.ThrottlePos,
		FuelLevel:      msg.FuelLevel,
		FuelTrimShort:  msg.FuelTrimShort,
		FuelTrimLong:   msg.FuelTrimLong,
		MAF:            msg.MAF,
		BatteryVoltage: msg.BatteryVoltage,
	})
}

func handleDTC(body []byte, uc *telemetryusecase.RecordDTCUseCase) error {
	var msg messaging.DTCMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("unmarshal dtc: %w", err)
	}
	if msg.VIN == "" || msg.Code == "" {
		return fmt.Errorf("missing vin or code")
	}
	return uc.Execute(context.Background(), telemetrydto.RecordDTCInput{
		VIN:    msg.VIN,
		Code:   msg.Code,
		Status: msg.Status,
		Time:   msg.Time,
	})
}
