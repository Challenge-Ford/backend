package telemetryrepository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type GormDTCRepository struct {
	db *gorm.DB
}

func NewGormDTCRepository(db *gorm.DB) *GormDTCRepository {
	return &GormDTCRepository{db: db}
}

func (r *GormDTCRepository) SetActive(ctx context.Context, deviceID uuid.UUID, vin, code string, at time.Time) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}, {Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"vin", "detected_at"}),
	}).Create(&telemetrydomain.ActiveDTC{
		DeviceID:   deviceID,
		VIN:        vin,
		Code:       code,
		DetectedAt: at,
	}).Error
}

func (r *GormDTCRepository) SetInactive(ctx context.Context, deviceID uuid.UUID, code string) error {
	return r.db.WithContext(ctx).
		Where("device_id = ? AND code = ?", deviceID, code).
		Delete(&telemetrydomain.ActiveDTC{}).Error
}

func (r *GormDTCRepository) ListActive(ctx context.Context, vin string) ([]*telemetrydomain.ActiveDTC, error) {
	var dtcs []*telemetrydomain.ActiveDTC
	return dtcs, r.db.WithContext(ctx).
		Where("vin = ?", vin).
		Order("detected_at DESC").
		Find(&dtcs).Error
}
