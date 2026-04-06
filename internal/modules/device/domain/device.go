package devicedomain

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"torque/internal/core/db"
	"torque/internal/core/id"
)

type DeviceID uuid.UUID

func NewDeviceID() DeviceID {
	return DeviceID(id.New())
}

func (d DeviceID) String() string {
	return uuid.UUID(d).String()
}

func (d DeviceID) Value() (driver.Value, error) {
	return uuid.UUID(d).String(), nil
}

func (d *DeviceID) Scan(src any) error {
	switch val := src.(type) {
	case string:
		parsed, err := uuid.Parse(val)
		if err != nil {
			return err
		}
		*d = DeviceID(parsed)
		return nil
	case []byte:
		parsed, err := uuid.Parse(string(val))
		if err != nil {
			return err
		}
		*d = DeviceID(parsed)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
}

type Device struct {
	ID            DeviceID   `gorm:"type:uuid;primaryKey"`
	Name          string     `gorm:"not null;uniqueIndex"`
	VehicleID     *uuid.UUID `gorm:"type:uuid;index"`
	CertificateCN string     `gorm:"not null;uniqueIndex"`
	CertificateSN string     `gorm:"not null"`
	db.AuditableModel

	// Populated via JOIN in List, not persisted
	VehicleVIN      *string `gorm:"->"`
	VehiclePlate    *string `gorm:"->"`
	VehicleColor    *string `gorm:"->"`
	VehicleModelName *string `gorm:"->"`
	VehicleYear     *int    `gorm:"->"`
	VehicleModelURL *string `gorm:"->"`
}

func (Device) TableName() string {
	return "device.devices"
}

func (d *Device) Delete(byUser uuid.UUID) {
	now := time.Now()
	d.DeletedAt = db.SoftDeletedAt{Time: now, Valid: true}
	d.DeletedBy = &byUser
}
