package vehicledomain

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"torque/internal/core/db"
	"torque/internal/core/id"
)

type VehicleID uuid.UUID

func NewVehicleID() VehicleID {
	return VehicleID(id.New())
}

func (v VehicleID) String() string {
	return uuid.UUID(v).String()
}

func (v VehicleID) Value() (driver.Value, error) {
	return uuid.UUID(v).String(), nil
}

func (v *VehicleID) Scan(src any) error {
	switch val := src.(type) {
	case string:
		parsed, err := uuid.Parse(val)
		if err != nil {
			return err
		}
		*v = VehicleID(parsed)
		return nil
	case []byte:
		parsed, err := uuid.Parse(string(val))
		if err != nil {
			return err
		}
		*v = VehicleID(parsed)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
}

type Vehicle struct {
	ID           VehicleID          `gorm:"type:uuid;primaryKey"`
	CustomerID   *uuid.UUID         `gorm:"type:uuid;index"`
	ModelYearID  VehicleModelYearID `gorm:"type:uuid;not null;index"`
	ModelYear    *VehicleModelYear  `gorm:"foreignKey:ModelYearID"`
	VIN          VIN                `gorm:"not null"`
	Plate        Plate              `gorm:"not null"`
	Color        Color              `gorm:"not null"`
	db.AuditableModel
}

func (Vehicle) TableName() string {
	return "vehicle.vehicles"
}

func (v *Vehicle) Delete(byUser uuid.UUID) {
	now := time.Now()
	v.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	v.DeletedBy = &byUser
}
