package telemetryrepository

import (
	"context"
	"time"

	"gorm.io/gorm"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Insert(ctx context.Context, entry *telemetrydomain.TelemetryEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *GormRepository) Latest(ctx context.Context, vin string) (*telemetrydomain.TelemetryEntry, error) {
	var entry telemetrydomain.TelemetryEntry
	err := r.db.WithContext(ctx).
		Where("vin = ?", vin).
		Order("time DESC").
		First(&entry).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &entry, err
}

func (r *GormRepository) List(ctx context.Context, vin string, from, to time.Time, limit int, after *time.Time) ([]*telemetrydomain.TelemetryEntry, error) {
	q := r.db.WithContext(ctx).
		Where("vin = ? AND time >= ? AND time <= ?", vin, from, to).
		Order("time ASC").
		Limit(limit)
	if after != nil {
		q = q.Where("time > ?", after)
	}
	var entries []*telemetrydomain.TelemetryEntry
	return entries, q.Find(&entries).Error
}

func (r *GormRepository) Summary(ctx context.Context, vin string, from, to time.Time, bucket string) ([]*telemetrydomain.TelemetrySummary, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
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
		FROM telemetry.entries
		WHERE vin = $2 AND time >= $3 AND time <= $4
		GROUP BY bucket
		ORDER BY bucket ASC
	`, bucket, vin, from, to).Rows()
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
