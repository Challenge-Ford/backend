package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"torque/internal/infrastructure/messaging"
)

const queueVehicleStateObserved = "torque.vehicle.state.observed"

func pf(v float64) *float64 { return &v }
func pi(v int) *int         { return &v }
func ps(v string) *string   { return &v }

func clamp(v, lo, hi float64) float64 {
	return math.Max(lo, math.Min(hi, v))
}

func noise(scale float64) float64 {
	return (rand.Float64()*2 - 1) * scale
}

func publish(ch *amqp.Channel, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ch.Publish("", queueVehicleStateObserved, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func dial() (*amqp.Connection, *amqp.Channel, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		fmt.Fprintln(os.Stderr, "error: RABBITMQ_URL is not set")
		os.Exit(1)
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("connect rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("open channel: %w", err)
	}
	if _, err := ch.QueueDeclare(queueVehicleStateObserved, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, fmt.Errorf("declare queue %s: %w", queueVehicleStateObserved, err)
	}
	return conn, ch, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: cli <command> [flags]

commands:
  state       publish a vehicle state snapshot
  simulate    publish a simulated drive session
  reset       truncate all vehicle state observations

Run "cli <command> -h" for flags.`)
	os.Exit(1)
}

func runReset(_ []string) {
	dsn := os.Getenv("TIMESERIES_DATABASE_URL")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "error: TIMESERIES_DATABASE_URL is not set")
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to connect to database:", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	if _, err := conn.Exec(context.Background(), "TRUNCATE TABLE vehicle_state_observations, vehicle_state_message_ids"); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println("vehicle state observations cleared")
}

func runState(args []string) {
	fs := flag.NewFlagSet("state", flag.ExitOnError)

	deviceID := fs.String("device-id", "", "device UUID (required)")
	vehicleID := fs.String("vehicle-id", "", "vehicle UUID (required)")
	rpm := fs.Int("rpm", 0, "engine RPM")
	speed := fs.Int("speed", 0, "vehicle speed km/h")
	coolant := fs.Float64("coolant-temp", 0, "coolant temperature C")
	intake := fs.Float64("intake-temp", 0, "intake air temperature C")
	load := fs.Float64("engine-load", 0, "engine load %")
	throttle := fs.Float64("throttle-pos", 0, "throttle position %")
	fuel := fs.Float64("fuel-level", 0, "fuel level %")
	trimShort := fs.Float64("fuel-trim-short", 0, "short-term fuel trim %")
	trimLong := fs.Float64("fuel-trim-long", 0, "long-term fuel trim %")
	maf := fs.Float64("maf", 0, "MAF air flow g/s")
	battery := fs.Float64("battery-voltage", 0, "battery voltage V")
	lat := fs.Float64("lat", 0, "GPS latitude")
	lng := fs.Float64("lng", 0, "GPS longitude")
	alt := fs.Float64("alt", 0, "GPS altitude m")
	gpsSpeed := fs.Float64("gps-speed", 0, "GPS speed km/h")
	heading := fs.Float64("heading", 0, "GPS heading degrees")
	hdop := fs.Float64("hdop", 0, "GPS HDOP")
	dtcs := fs.String("dtcs", "", "comma-separated open DTC codes")

	fs.Parse(args)

	if *deviceID == "" || *vehicleID == "" {
		fmt.Fprintln(os.Stderr, "error: --device-id and --vehicle-id are required")
		fs.Usage()
		os.Exit(1)
	}

	msg := baseMessage(*deviceID, *vehicleID)
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "rpm":
			ensurePowertrain(&msg).RPM = pi(*rpm)
		case "speed":
			ensurePowertrain(&msg).Speed = pi(*speed)
		case "coolant-temp":
			ensurePowertrain(&msg).CoolantTemp = pf(*coolant)
		case "intake-temp":
			ensurePowertrain(&msg).IntakeTemp = pf(*intake)
		case "engine-load":
			ensurePowertrain(&msg).EngineLoad = pf(*load)
		case "throttle-pos":
			ensurePowertrain(&msg).ThrottlePos = pf(*throttle)
		case "maf":
			ensurePowertrain(&msg).MAF = pf(*maf)
		case "fuel-level":
			ensureFuel(&msg).Level = pf(*fuel)
		case "fuel-trim-short":
			ensureFuel(&msg).TrimShort = pf(*trimShort)
		case "fuel-trim-long":
			ensureFuel(&msg).TrimLong = pf(*trimLong)
		case "battery-voltage":
			ensureElectrical(&msg).BatteryVoltage = pf(*battery)
		case "lat":
			ensurePosition(&msg).Lat = pf(*lat)
		case "lng":
			ensurePosition(&msg).Lng = pf(*lng)
		case "alt":
			ensurePosition(&msg).Alt = pf(*alt)
		case "gps-speed":
			ensurePosition(&msg).Speed = pf(*gpsSpeed)
		case "heading":
			ensurePosition(&msg).Heading = pf(*heading)
		case "hdop":
			ensurePosition(&msg).HDOP = pf(*hdop)
		case "dtcs":
			msg.State.Diagnostics = &messaging.DiagnosticsState{OpenDTCs: splitCSV(*dtcs)}
		}
	})

	conn, ch, err := dial()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	defer ch.Close()

	if err := publish(ch, msg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("published vehicle state")
}

var waypoints = [][2]float64{
	{-23.5505, -46.6333},
	{-23.5450, -46.6100},
	{-23.5420, -46.5800},
	{-23.5200, -46.5600},
	{-23.4900, -46.5200},
	{-23.4600, -46.4800},
	{-23.4300, -46.4400},
}

func runSimulate(args []string) {
	fs := flag.NewFlagSet("simulate", flag.ExitOnError)
	deviceID := fs.String("device-id", "", "device UUID (required)")
	vehicleID := fs.String("vehicle-id", "", "vehicle UUID (required)")
	count := fs.Int("count", 60, "number of snapshots")
	interval := fs.Duration("interval", time.Second, "interval between snapshots")
	fs.Parse(args)

	if *deviceID == "" || *vehicleID == "" {
		fmt.Fprintln(os.Stderr, "error: --device-id and --vehicle-id are required")
		fs.Usage()
		os.Exit(1)
	}

	conn, ch, err := dial()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	defer ch.Close()

	for i := 0; i < *count; i++ {
		msg := simulatedMessage(*deviceID, *vehicleID, i, *count)
		if err := publish(ch, msg); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		time.Sleep(*interval)
	}
	fmt.Println("simulation complete")
}

func baseMessage(deviceID, vehicleID string) messaging.VehicleStateObservedMessage {
	return messaging.VehicleStateObservedMessage{
		SchemaVersion: 1,
		MessageID:     uuid.NewString(),
		DeviceID:      deviceID,
		VehicleID:     vehicleID,
		ObservedAt:    time.Now().UTC(),
		State:         messaging.VehicleState{},
		Observation:   messaging.ObservationMetadata{Errors: []messaging.ObservationError{}},
	}
}

func simulatedMessage(deviceID, vehicleID string, i, n int) messaging.VehicleStateObservedMessage {
	progress := float64(i) / math.Max(1, float64(n-1))
	wp := waypoints[int(progress*float64(len(waypoints)-1))]
	speed := int(clamp(65+noise(25), 0, 120))
	rpm := int(clamp(float64(900+speed*32)+noise(300), 700, 4200))
	fuel := clamp(78-progress*4+noise(0.2), 0, 100)
	msg := baseMessage(deviceID, vehicleID)
	msg.State.Position = &messaging.PositionState{
		Source:  ps("gps"),
		Lat:     pf(wp[0] + noise(0.001)),
		Lng:     pf(wp[1] + noise(0.001)),
		Alt:     pf(760 + noise(12)),
		Speed:   pf(float64(speed)),
		Heading: pf(clamp(120+noise(20), 0, 360)),
		HDOP:    pf(clamp(0.8+noise(0.2), 0.4, 2.0)),
	}
	msg.State.Powertrain = &messaging.PowertrainState{
		RPM:         pi(rpm),
		Speed:       pi(speed),
		EngineLoad:  pf(clamp(35+noise(20), 0, 100)),
		ThrottlePos: pf(clamp(18+noise(12), 0, 100)),
		CoolantTemp: pf(clamp(88+noise(4), 70, 105)),
		IntakeTemp:  pf(clamp(34+noise(5), 10, 60)),
		MAF:         pf(clamp(10+float64(speed)*0.08+noise(2), 0, 80)),
	}
	msg.State.Fuel = &messaging.FuelState{Level: pf(fuel), TrimShort: pf(noise(3)), TrimLong: pf(noise(2))}
	msg.State.Electrical = &messaging.ElectricalState{BatteryVoltage: pf(clamp(13.8+noise(0.25), 11.8, 14.8))}
	msg.State.Diagnostics = &messaging.DiagnosticsState{OpenDTCs: []string{}}
	return msg
}

func ensurePosition(msg *messaging.VehicleStateObservedMessage) *messaging.PositionState {
	if msg.State.Position == nil {
		msg.State.Position = &messaging.PositionState{Source: ps("gps")}
	}
	return msg.State.Position
}

func ensurePowertrain(msg *messaging.VehicleStateObservedMessage) *messaging.PowertrainState {
	if msg.State.Powertrain == nil {
		msg.State.Powertrain = &messaging.PowertrainState{}
	}
	return msg.State.Powertrain
}

func ensureFuel(msg *messaging.VehicleStateObservedMessage) *messaging.FuelState {
	if msg.State.Fuel == nil {
		msg.State.Fuel = &messaging.FuelState{}
	}
	return msg.State.Fuel
}

func ensureElectrical(msg *messaging.VehicleStateObservedMessage) *messaging.ElectricalState {
	if msg.State.Electrical == nil {
		msg.State.Electrical = &messaging.ElectricalState{}
	}
	return msg.State.Electrical
}

func splitCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	return out
}

func main() {
	godotenv.Load()
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "state":
		runState(os.Args[2:])
	case "simulate":
		runSimulate(os.Args[2:])
	case "reset":
		runReset(os.Args[2:])
	default:
		usage()
	}
}
