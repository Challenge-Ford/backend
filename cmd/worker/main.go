package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"torque/internal/core/db"
	"torque/internal/core/logger"
	"torque/internal/infrastructure/adapters"
	"torque/internal/infrastructure/messaging"
	deviceusecase "torque/internal/modules/device/application/usecase"
	devicerepository "torque/internal/modules/device/infrastructure/repository"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetryusecase "torque/internal/modules/telemetry/application/usecase"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	telemetryrepository "torque/internal/modules/telemetry/infrastructure/repository"
)

const (
	queueVehicleStateObserved = "torque.vehicle.state.observed"
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

	if _, err := ch.QueueDeclare(queueVehicleStateObserved, true, false, false, false, nil); err != nil {
		log.Fatal("failed to declare queue", zap.String("queue", queueVehicleStateObserved), zap.Error(err))
	}

	deviceRepo := devicerepository.NewRepository(pool)
	stateRepo := telemetryrepository.NewPgxStateObservationRepository(tsPool)
	findCommissionedByVIN := deviceusecase.NewFindCommissionedByVIN(deviceRepo)
	findDeviceByVehicle := deviceusecase.NewFindDeviceByVehicle(deviceRepo)
	deviceResolver := adapters.NewDeviceResolver(findCommissionedByVIN, findDeviceByVehicle)

	recordVehicleState := telemetryusecase.NewRecordVehicleState(stateRepo, deviceResolver)

	go consume(log, ch, queueVehicleStateObserved, func(body []byte) error {
		return handleVehicleStateObserved(body, recordVehicleState)
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

func handleVehicleStateObserved(body []byte, uc *telemetryusecase.RecordVehicleStateUseCase) error {
	var msg messaging.VehicleStateObservedMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("unmarshal vehicle state observed: %w", err)
	}
	messageID, err := uuid.Parse(msg.MessageID)
	if err != nil {
		return fmt.Errorf("invalid message_id: %w", err)
	}
	deviceID, err := uuid.Parse(msg.DeviceID)
	if err != nil {
		return fmt.Errorf("invalid device_id: %w", err)
	}
	vehicleID, err := uuid.Parse(msg.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle_id: %w", err)
	}

	return uc.Execute(context.Background(), telemetrydto.RecordVehicleStateInput{
		SchemaVersion: msg.SchemaVersion,
		MessageID:     messageID,
		DeviceID:      deviceID,
		VehicleID:     vehicleID,
		ObservedAt:    msg.ObservedAt,
		State: telemetrydomain.VehicleState{
			Position:    toDomainPosition(msg.State.Position),
			Powertrain:  toDomainPowertrain(msg.State.Powertrain),
			Fuel:        toDomainFuel(msg.State.Fuel),
			Electrical:  toDomainElectrical(msg.State.Electrical),
			Diagnostics: toDomainDiagnostics(msg.State.Diagnostics),
		},
		Observation: toDomainObservation(msg.Observation),
		RawPayload:  append([]byte(nil), body...),
	})
}

func toDomainPosition(in *messaging.PositionState) *telemetrydomain.PositionState {
	if in == nil {
		return nil
	}
	return &telemetrydomain.PositionState{
		Source: in.Source, Lat: in.Lat, Lng: in.Lng, Alt: in.Alt,
		Speed: in.Speed, Heading: in.Heading, HDOP: in.HDOP,
	}
}

func toDomainPowertrain(in *messaging.PowertrainState) *telemetrydomain.PowertrainState {
	if in == nil {
		return nil
	}
	return &telemetrydomain.PowertrainState{
		RPM: in.RPM, Speed: in.Speed, EngineLoad: in.EngineLoad,
		ThrottlePos: in.ThrottlePos, CoolantTemp: in.CoolantTemp,
		IntakeTemp: in.IntakeTemp, MAF: in.MAF,
	}
}

func toDomainFuel(in *messaging.FuelState) *telemetrydomain.FuelState {
	if in == nil {
		return nil
	}
	return &telemetrydomain.FuelState{Level: in.Level, TrimShort: in.TrimShort, TrimLong: in.TrimLong}
}

func toDomainElectrical(in *messaging.ElectricalState) *telemetrydomain.ElectricalState {
	if in == nil {
		return nil
	}
	return &telemetrydomain.ElectricalState{BatteryVoltage: in.BatteryVoltage}
}

func toDomainDiagnostics(in *messaging.DiagnosticsState) *telemetrydomain.DiagnosticsState {
	if in == nil {
		return nil
	}
	return &telemetrydomain.DiagnosticsState{OpenDTCs: in.OpenDTCs}
}

func toDomainObservation(in messaging.ObservationMetadata) telemetrydomain.ObservationMetadata {
	errors := make([]telemetrydomain.ObservationError, len(in.Errors))
	for i, e := range in.Errors {
		errors[i] = telemetrydomain.ObservationError{Block: e.Block, Code: e.Code, Message: e.Message}
	}
	return telemetrydomain.ObservationMetadata{Errors: errors}
}
