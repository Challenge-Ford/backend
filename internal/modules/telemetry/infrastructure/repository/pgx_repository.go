package telemetryrepository

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

const selectFields = `
	time, device_id, vin,
	lat, lng, alt, gps_speed, heading, hdop,
	rpm, speed, coolant_temp, intake_temp, engine_load,
	throttle_pos, fuel_level, fuel_trim_short, fuel_trim_long, maf, battery_voltage`

type PgxRepository struct {
	pool *pgxpool.Pool
}

func NewPgxRepository(pool *pgxpool.Pool) *PgxRepository {
	return &PgxRepository{pool: pool}
}

func (r *PgxRepository) Insert(ctx context.Context, e *telemetrydomain.TelemetryEntry) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO telemetry_entries (
			time, device_id, vin,
			lat, lng, alt, gps_speed, heading, hdop,
			rpm, speed, coolant_temp, intake_temp, engine_load,
			throttle_pos, fuel_level, fuel_trim_short, fuel_trim_long, maf, battery_voltage
		) VALUES (
			$1, $2, $3,
			$4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20
		) ON CONFLICT (time, device_id) DO NOTHING`,
		e.Time, e.DeviceID, e.VIN,
		e.Lat, e.Lng, e.Alt, e.GPSSpeed, e.Heading, e.HDOP,
		e.RPM, e.Speed, e.CoolantTemp, e.IntakeTemp, e.EngineLoad,
		e.ThrottlePos, e.FuelLevel, e.FuelTrimShort, e.FuelTrimLong, e.MAF, e.BatteryVoltage,
	)
	return err
}

func (r *PgxRepository) Latest(ctx context.Context, vin string) (*telemetrydomain.TelemetryEntry, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT `+selectFields+`
		FROM telemetry_entries
		WHERE vin = $1
		ORDER BY time DESC
		LIMIT 1`, vin)

	e, err := scanEntry(row)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *PgxRepository) List(ctx context.Context, vin string, from, to *time.Time, limit int, after *time.Time) ([]*telemetrydomain.TelemetryEntry, error) {
	args := []any{vin}
	q := `
		SELECT ` + selectFields + `
		FROM telemetry_entries
		WHERE vin = $1`
	argIdx := 2
	if from != nil {
		q += ` AND time >= $` + strconv.Itoa(argIdx)
		args = append(args, *from)
		argIdx++
	}
	if to != nil {
		q += ` AND time <= $` + strconv.Itoa(argIdx)
		args = append(args, *to)
		argIdx++
	}
	if after != nil {
		q += ` AND time > $` + strconv.Itoa(argIdx)
		args = append(args, *after)
		argIdx++
	}
	q += ` ORDER BY time ASC LIMIT $` + strconv.Itoa(argIdx)
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*telemetrydomain.TelemetryEntry
	for rows.Next() {
		e, err := scanEntry(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *PgxRepository) Summary(ctx context.Context, vin string, from, to time.Time, bucket string) ([]*telemetrydomain.TelemetrySummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			time_bucket($1::interval, time) AS bucket,
			AVG(rpm)::float             AS avg_rpm,
			MAX(rpm)::float             AS max_rpm,
			AVG(speed)::float           AS avg_speed,
			MAX(speed)::float           AS max_speed,
			AVG(coolant_temp)           AS avg_coolant_temp,
			MAX(coolant_temp)           AS max_coolant_temp,
			AVG(engine_load)            AS avg_engine_load,
			AVG(maf)                    AS avg_maf,
			AVG(battery_voltage)        AS avg_battery_voltage
		FROM telemetry_entries
		WHERE vin = $2 AND time >= $3 AND time <= $4
		GROUP BY bucket
		ORDER BY bucket ASC`,
		bucket, vin, from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*telemetrydomain.TelemetrySummary
	for rows.Next() {
		s := &telemetrydomain.TelemetrySummary{}
		if err := rows.Scan(
			&s.Bucket,
			&s.AvgRPM, &s.MaxRPM,
			&s.AvgSpeed, &s.MaxSpeed,
			&s.AvgCoolantTemp, &s.MaxCoolantTemp,
			&s.AvgEngineLoad,
			&s.AvgMAF,
			&s.AvgBatteryVoltage,
		); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}

// scanner is satisfied by both pgx.Row and pgx.Rows
type scanner interface {
	Scan(dest ...any) error
}

func scanEntry(s scanner) (*telemetrydomain.TelemetryEntry, error) {
	e := &telemetrydomain.TelemetryEntry{}
	err := s.Scan(
		&e.Time, &e.DeviceID, &e.VIN,
		&e.Lat, &e.Lng, &e.Alt, &e.GPSSpeed, &e.Heading, &e.HDOP,
		&e.RPM, &e.Speed, &e.CoolantTemp, &e.IntakeTemp, &e.EngineLoad,
		&e.ThrottlePos, &e.FuelLevel, &e.FuelTrimShort, &e.FuelTrimLong, &e.MAF, &e.BatteryVoltage,
	)
	if err != nil {
		return nil, err
	}
	return e, nil
}
