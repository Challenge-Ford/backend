package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"torque/internal/infrastructure/messaging"
)

const (
	queueTelemetry = "torque.telemetry"
	queueDTC       = "torque.dtc"
)


func pf(v float64) *float64 { return &v }
func pi(v int) *int         { return &v }
func pt(v time.Time) *time.Time { return &v }

func clamp(v, lo, hi float64) float64 {
	return math.Max(lo, math.Min(hi, v))
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func noise(scale float64) float64 {
	return (rand.Float64()*2 - 1) * scale
}

func publish(ch *amqp.Channel, queue string, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ch.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func dial() (*amqp.Connection, *amqp.Channel, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://torque:torque@localhost:5672/"
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
	for _, q := range []string{queueTelemetry, queueDTC} {
		if _, err := ch.QueueDeclare(q, true, false, false, false, nil); err != nil {
			ch.Close()
			conn.Close()
			return nil, nil, fmt.Errorf("declare queue %s: %w", q, err)
		}
	}
	return conn, ch, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: cli <command> [flags]

commands:
  telemetry   publish a single telemetry event
  dtc         publish a DTC event
  simulate    publish a simulated drive session

Run "cli <command> -h" for flags.`)
	os.Exit(1)
}

// ── telemetry ──────────────────────────────────────────────────────────────

func runTelemetry(args []string) {
	fs := flag.NewFlagSet("telemetry", flag.ExitOnError)

	vin := fs.String("vin", "", "vehicle VIN (required)")
	rpm := fs.Int("rpm", 0, "engine RPM")
	speed := fs.Int("speed", 0, "vehicle speed km/h")
	coolant := fs.Float64("coolant-temp", 0, "coolant temperature °C")
	intake := fs.Float64("intake-temp", 0, "intake air temperature °C")
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

	fs.Parse(args)

	if *vin == "" {
		fmt.Fprintln(os.Stderr, "error: --vin is required")
		fs.Usage()
		os.Exit(1)
	}

	msg := messaging.TelemetryMessage{VIN: *vin}
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "rpm":           msg.RPM = pi(*rpm)
		case "speed":         msg.Speed = pi(*speed)
		case "coolant-temp":  msg.CoolantTemp = pf(*coolant)
		case "intake-temp":   msg.IntakeTemp = pf(*intake)
		case "engine-load":   msg.EngineLoad = pf(*load)
		case "throttle-pos":  msg.ThrottlePos = pf(*throttle)
		case "fuel-level":    msg.FuelLevel = pf(*fuel)
		case "fuel-trim-short": msg.FuelTrimShort = pf(*trimShort)
		case "fuel-trim-long":  msg.FuelTrimLong = pf(*trimLong)
		case "maf":           msg.MAF = pf(*maf)
		case "battery-voltage": msg.BatteryVoltage = pf(*battery)
		case "lat":           msg.Lat = pf(*lat)
		case "lng":           msg.Lng = pf(*lng)
		case "alt":           msg.Alt = pf(*alt)
		case "gps-speed":     msg.GPSSpeed = pf(*gpsSpeed)
		case "heading":       msg.Heading = pf(*heading)
		case "hdop":          msg.HDOP = pf(*hdop)
		}
	})

	conn, ch, err := dial()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	defer ch.Close()

	if err := publish(ch, queueTelemetry, msg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("published telemetry")
}

// ── dtc ───────────────────────────────────────────────────────────────────

func runDTC(args []string) {
	fs := flag.NewFlagSet("dtc", flag.ExitOnError)

	vin := fs.String("vin", "", "vehicle VIN (required)")
	code := fs.String("code", "", "DTC code e.g. P0300 (required)")
	status := fs.String("status", "", "opened or closed (required)")

	fs.Parse(args)

	if *vin == "" || *code == "" || *status == "" {
		fmt.Fprintln(os.Stderr, "error: --vin, --code and --status are required")
		fs.Usage()
		os.Exit(1)
	}
	if *status != "opened" && *status != "closed" {
		fmt.Fprintln(os.Stderr, "error: --status must be 'opened' or 'closed'")
		os.Exit(1)
	}

	conn, ch, err := dial()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	defer ch.Close()

	if err := publish(ch, queueDTC, messaging.DTCMessage{VIN: *vin, Code: *code, Status: *status}); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("published dtc %s %s\n", *code, *status)
}

// ── simulate ───────────────────────────────────────────────────────────────

// GPS waypoints simulating a São Paulo commute (Centro → Guarulhos direction)
var waypoints = [][2]float64{
	{-23.5505, -46.6333},
	{-23.5450, -46.6100},
	{-23.5420, -46.5800},
	{-23.5200, -46.5600},
	{-23.4900, -46.5200},
	{-23.4600, -46.4800},
	{-23.4300, -46.4400},
	{-23.4300, -46.4400}, // parked
	{-23.4300, -46.4400},
	{-23.4600, -46.4800},
	{-23.4900, -46.5200},
	{-23.5200, -46.5600},
	{-23.5420, -46.5800},
	{-23.5450, -46.6100},
	{-23.5505, -46.6333},
}

func gpsAt(i, n int) (lat, lng float64) {
	t := float64(i) / float64(n-1) * float64(len(waypoints)-1)
	idx := int(t)
	frac := t - float64(idx)
	if idx >= len(waypoints)-1 {
		return waypoints[len(waypoints)-1][0], waypoints[len(waypoints)-1][1]
	}
	a, b := waypoints[idx], waypoints[idx+1]
	return lerp(a[0], b[0], frac) + noise(0.001), lerp(a[1], b[1], frac) + noise(0.001)
}

func simulatePoint(i, n int, fuelBase float64) messaging.TelemetryMessage {
	progress := float64(i) / float64(n-1)

	var (
		rpm        float64
		speed      float64
		coolant    float64
		load       float64
		throttle   float64
		fuel       float64
		maf        float64
		battery    float64
		gpsSpeed   float64
		heading    float64
	)

	switch {
	case progress < 0.05: // cold start
		p := progress / 0.05
		rpm = clamp(1200-p*300+noise(50), 900, 1300)
		speed = 0
		coolant = lerp(22, 60, p)
		load = clamp(25+noise(5), 18, 35)
		throttle = clamp(8+noise(2), 5, 12)
		fuel = fuelBase - progress*float64(n)*0.05
		maf = clamp(3.5+noise(0.5), 2.5, 5.0)
		battery = clamp(13.9+noise(0.2), 13.5, 14.5)
		gpsSpeed = 0
		heading = 45

	case progress < 0.20: // city outbound
		p := (progress - 0.05) / 0.15
		rpm = clamp(1800+noise(400), 1200, 3200)
		speed = clamp(35+noise(20), 0, 60)
		coolant = lerp(60, 92, p) + noise(1)
		load = clamp(45+noise(15), 25, 65)
		throttle = clamp(22+noise(10), 8, 45)
		fuel = fuelBase - float64(i)*0.12
		maf = clamp(8+noise(3), 4, 16)
		battery = clamp(14.1+noise(0.15), 13.8, 14.5)
		gpsSpeed = clamp(speed+noise(3), 0, 65)
		heading = clamp(45+noise(15), 0, 359)

	case progress < 0.40: // highway outbound
		rpm = clamp(2400+noise(200), 2000, 3000)
		speed = clamp(100+noise(15), 80, 130)
		coolant = clamp(92+noise(2), 88, 98)
		load = clamp(55+noise(10), 40, 70)
		throttle = clamp(35+noise(8), 20, 55)
		fuel = fuelBase - float64(i)*0.18
		maf = clamp(16+noise(3), 11, 22)
		battery = clamp(14.2+noise(0.1), 14.0, 14.5)
		gpsSpeed = clamp(speed+noise(2), 75, 135)
		heading = clamp(45+noise(10), 0, 359)

	case progress < 0.50: // parked
		p := (progress - 0.40) / 0.10
		rpm = 0
		speed = 0
		coolant = clamp(92-p*40+noise(1), 50, 95)
		load = 0
		throttle = 0
		fuel = fuelBase - float64(i)*0.18
		maf = 0
		battery = clamp(12.6+noise(0.1), 12.3, 12.9)
		gpsSpeed = 0
		heading = 225

	case progress < 0.55: // warm restart
		p := (progress - 0.50) / 0.05
		rpm = clamp(900-p*150+noise(30), 750, 1000)
		speed = 0
		coolant = clamp(lerp(52, 80, p)+noise(2), 50, 85)
		load = clamp(20+noise(4), 15, 28)
		throttle = clamp(7+noise(2), 5, 10)
		fuel = fuelBase - float64(i)*0.04
		maf = clamp(2.8+noise(0.4), 2.0, 4.0)
		battery = clamp(13.8+noise(0.2), 13.5, 14.2)
		gpsSpeed = 0
		heading = 225

	case progress < 0.70: // city return
		rpm = clamp(1900+noise(400), 1200, 3200)
		speed = clamp(38+noise(18), 0, 60)
		coolant = clamp(90+noise(2), 86, 96)
		load = clamp(47+noise(15), 25, 65)
		throttle = clamp(24+noise(10), 8, 45)
		fuel = fuelBase - float64(i)*0.12
		maf = clamp(9+noise(3), 4, 16)
		battery = clamp(14.1+noise(0.15), 13.8, 14.5)
		gpsSpeed = clamp(speed+noise(3), 0, 65)
		heading = clamp(225+noise(15), 0, 359)

	case progress < 0.90: // highway return
		rpm = clamp(2350+noise(200), 2000, 2900)
		speed = clamp(98+noise(15), 80, 130)
		coolant = clamp(91+noise(2), 87, 97)
		load = clamp(53+noise(10), 38, 68)
		throttle = clamp(33+noise(8), 18, 52)
		fuel = fuelBase - float64(i)*0.18
		maf = clamp(15.5+noise(3), 10, 22)
		battery = clamp(14.2+noise(0.1), 14.0, 14.5)
		gpsSpeed = clamp(speed+noise(2), 75, 130)
		heading = clamp(225+noise(10), 0, 359)

	default: // arrival
		p := (progress - 0.90) / 0.10
		if p < 0.8 {
			rpm = clamp(lerp(1600, 750, p)+noise(100), 0, 2000)
		}
		speed = clamp(lerp(40, 0, p), 0, 45)
		coolant = clamp(90-p*5+noise(1), 70, 93)
		load = clamp(lerp(40, 0, p)+noise(5), 0, 50)
		throttle = clamp(lerp(20, 0, p)+noise(3), 0, 30)
		fuel = fuelBase - float64(i)*0.08
		maf = clamp(lerp(7, 0, p)+noise(1), 0, 12)
		battery = clamp(lerp(14.1, 12.5, p)+noise(0.1), 12.3, 14.5)
		gpsSpeed = speed
		heading = clamp(225+noise(15), 0, 359)
	}

	lat, lng := gpsAt(i, n)
	alt := clamp(760+noise(5), 748, 775)
	intake := clamp(28+noise(4), 20, 45)
	fuelLevel := clamp(fuel, 0, 100)
	trimShort := noise(3)
	trimLong := noise(1.5)
	hdop := clamp(1.0+noise(0.3), 0.7, 2.0)

	msg := messaging.TelemetryMessage{
		VIN:            "",
		Lat:            pf(math.Round(lat*1e6) / 1e6),
		Lng:            pf(math.Round(lng*1e6) / 1e6),
		Alt:            pf(math.Round(alt*10) / 10),
		GPSSpeed:       pf(math.Round(gpsSpeed*10) / 10),
		Heading:        pf(math.Round(heading*10) / 10),
		HDOP:           pf(math.Round(hdop*10) / 10),
		CoolantTemp:    pf(math.Round(coolant*10) / 10),
		IntakeTemp:     pf(math.Round(intake*10) / 10),
		EngineLoad:     pf(math.Round(load*10) / 10),
		ThrottlePos:    pf(math.Round(throttle*10) / 10),
		FuelLevel:      pf(math.Round(fuelLevel*10) / 10),
		FuelTrimShort:  pf(math.Round(trimShort*100) / 100),
		FuelTrimLong:   pf(math.Round(trimLong*100) / 100),
		MAF:            pf(math.Round(maf*100) / 100),
		BatteryVoltage: pf(math.Round(battery*100) / 100),
	}

	if rpm > 0 {
		msg.RPM = pi(int(rpm))
	}
	if speed > 0 {
		msg.Speed = pi(int(speed))
	}

	return msg
}

func runSimulate(args []string) {
	fs := flag.NewFlagSet("simulate", flag.ExitOnError)

	vin := fs.String("vin", "", "vehicle VIN (required)")
	points := fs.Int("points", 100, "number of data points")
	intervalMin := fs.Int("interval", 5, "minutes between data points")
	seed := fs.Int64("seed", time.Now().UnixNano(), "random seed")

	fs.Parse(args)

	if *vin == "" {
		fmt.Fprintln(os.Stderr, "error: --vin is required")
		fs.Usage()
		os.Exit(1)
	}

	rand.New(rand.NewSource(*seed))

	conn, ch, err := dial()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	defer ch.Close()

	interval := time.Duration(*intervalMin) * time.Minute
	start := time.Now().Add(-interval * time.Duration(*points))
	fuelBase := 78.0

	ok, fail := 0, 0
	for i := 0; i < *points; i++ {
		t := start.Add(interval * time.Duration(i))
		msg := simulatePoint(i, *points, fuelBase)
		msg.VIN = *vin
		msg.Time = pt(t.UTC())

		if err := publish(ch, queueTelemetry, msg); err != nil {
			fmt.Fprintf(os.Stderr, "point %d: %v\n", i, err)
			fail++
		} else {
			ok++
		}
	}

	fmt.Printf("simulate done — published: %d, failed: %d\n", ok, fail)
}

func main() {
	godotenv.Load()

	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "telemetry":
		runTelemetry(os.Args[2:])
	case "dtc":
		runDTC(os.Args[2:])
	case "simulate":
		runSimulate(os.Args[2:])
	default:
		usage()
	}
}
