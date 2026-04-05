package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"torque/internal/core/db"
	"torque/internal/core/logger"
	telemetrydomain "torque/internal/modules/telemetry/domain"
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

	// Primary DB — device/vehicle lookups
	mainConn, err := db.Connect(mustEnv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close(mainConn)

	// TimescaleDB — telemetry persistence
	tsConn, err := db.Connect(mustEnv("TIMESERIES_DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to timescaledb", zap.Error(err))
	}
	defer db.Close(tsConn)

	if err := migrate(tsConn); err != nil {
		log.Fatal("migration failed", zap.Error(err))
	}

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

	telemetryRepo := telemetryrepository.NewGormRepository(tsConn)
	dtcRepo := telemetryrepository.NewGormDTCRepository(tsConn)

	go consume(log, ch, queueTelemetry, func(body []byte) error {
		return handleTelemetry(body, mainConn, telemetryRepo)
	})

	go consume(log, ch, queueDTC, func(body []byte) error {
		return handleDTC(body, mainConn, dtcRepo)
	})

	log.Info("worker started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("shutting down")
}

func migrate(tsConn *gorm.DB) error {
	if err := tsConn.Exec("CREATE SCHEMA IF NOT EXISTS telemetry").Error; err != nil {
		return err
	}

	if err := tsConn.Exec(`
		CREATE TABLE IF NOT EXISTS telemetry.entries (
			time            TIMESTAMPTZ      NOT NULL,
			device_id       UUID             NOT NULL,
			vin             TEXT             NOT NULL,
			lat             DOUBLE PRECISION,
			lng             DOUBLE PRECISION,
			alt             DOUBLE PRECISION,
			gps_speed       DOUBLE PRECISION,
			heading         DOUBLE PRECISION,
			hdop            DOUBLE PRECISION,
			rpm             INTEGER,
			speed           INTEGER,
			coolant_temp    DOUBLE PRECISION,
			intake_temp     DOUBLE PRECISION,
			engine_load     DOUBLE PRECISION,
			throttle_pos    DOUBLE PRECISION,
			fuel_level      DOUBLE PRECISION,
			fuel_trim_short DOUBLE PRECISION,
			fuel_trim_long  DOUBLE PRECISION,
			maf             DOUBLE PRECISION,
			battery_voltage DOUBLE PRECISION,
			PRIMARY KEY (time, device_id)
		)
	`).Error; err != nil {
		return err
	}

	if err := tsConn.Exec(`
		SELECT create_hypertable('telemetry.entries', 'time', if_not_exists => true)
	`).Error; err != nil {
		return err
	}

	return db.Migrate(tsConn, &telemetrydomain.ActiveDTC{})
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

type telemetryMsg struct {
	Time           *time.Time `json:"time"`
	VIN            string     `json:"vin"`
	Lat            *float64   `json:"lat"`
	Lng            *float64   `json:"lng"`
	Alt            *float64   `json:"alt"`
	GPSSpeed       *float64   `json:"gps_speed"`
	Heading        *float64   `json:"heading"`
	HDOP           *float64   `json:"hdop"`
	RPM            *int       `json:"rpm"`
	Speed          *int       `json:"speed"`
	CoolantTemp    *float64   `json:"coolant_temp"`
	IntakeTemp     *float64   `json:"intake_temp"`
	EngineLoad     *float64   `json:"engine_load"`
	ThrottlePos    *float64   `json:"throttle_pos"`
	FuelLevel      *float64   `json:"fuel_level"`
	FuelTrimShort  *float64   `json:"fuel_trim_short"`
	FuelTrimLong   *float64   `json:"fuel_trim_long"`
	MAF            *float64   `json:"maf"`
	BatteryVoltage *float64   `json:"battery_voltage"`
}

type dtcMsg struct {
	VIN    string `json:"vin"`
	Code   string `json:"code"`
	Status string `json:"status"`
}

func handleTelemetry(body []byte, mainConn *gorm.DB, repo *telemetryrepository.GormRepository) error {
	var msg telemetryMsg
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("unmarshal telemetry: %w", err)
	}
	if msg.VIN == "" {
		return fmt.Errorf("missing vin")
	}

	deviceID, err := lookupDeviceByVIN(mainConn, msg.VIN)
	if err != nil {
		return fmt.Errorf("lookup device for vin %s: %w", msg.VIN, err)
	}

	t := time.Now().UTC()
	if msg.Time != nil {
		t = msg.Time.UTC()
	}

	return repo.Insert(context.Background(), &telemetrydomain.TelemetryEntry{
		Time:           t,
		DeviceID:       deviceID,
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

func handleDTC(body []byte, mainConn *gorm.DB, repo *telemetryrepository.GormDTCRepository) error {
	var msg dtcMsg
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("unmarshal dtc: %w", err)
	}
	if msg.VIN == "" || msg.Code == "" {
		return fmt.Errorf("missing vin or code")
	}
	if msg.Status != "opened" && msg.Status != "closed" {
		return fmt.Errorf("invalid status: %s", msg.Status)
	}

	deviceID, err := lookupDeviceByVIN(mainConn, msg.VIN)
	if err != nil {
		return fmt.Errorf("lookup device for vin %s: %w", msg.VIN, err)
	}

	ctx := context.Background()
	if msg.Status == "opened" {
		return repo.SetActive(ctx, deviceID, msg.VIN, msg.Code, time.Now().UTC())
	}
	return repo.SetInactive(ctx, deviceID, msg.Code)
}

func lookupDeviceByVIN(mainConn *gorm.DB, vin string) (uuid.UUID, error) {
	var id string
	err := mainConn.Raw(`
		SELECT d.id
		FROM device.devices d
		JOIN vehicle.vehicles v ON v.id = d.vehicle_id AND v.deleted_at IS NULL
		WHERE v.vin = $1
		  AND d.vehicle_id IS NOT NULL
		  AND d.deleted_at IS NULL
		LIMIT 1
	`, vin).Scan(&id).Error
	if err != nil {
		return uuid.Nil, err
	}
	if id == "" {
		return uuid.Nil, fmt.Errorf("no commissioned device found for vin %s", vin)
	}
	return uuid.Parse(id)
}
