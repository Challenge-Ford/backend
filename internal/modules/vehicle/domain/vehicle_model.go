package vehicledomain

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"torque/internal/core/id"
)

type VehicleModelID uuid.UUID

func NewVehicleModelID() VehicleModelID {
	return VehicleModelID(id.New())
}

func (v VehicleModelID) String() string {
	return uuid.UUID(v).String()
}

func (v VehicleModelID) Value() (driver.Value, error) {
	return uuid.UUID(v).String(), nil
}

func (v *VehicleModelID) Scan(src any) error {
	switch val := src.(type) {
	case string:
		parsed, err := uuid.Parse(val)
		if err != nil {
			return err
		}
		*v = VehicleModelID(parsed)
		return nil
	case []byte:
		parsed, err := uuid.Parse(string(val))
		if err != nil {
			return err
		}
		*v = VehicleModelID(parsed)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", val)
	}
}

type VehicleModel struct {
	ID   VehicleModelID `gorm:"type:uuid;primaryKey"`
	Name string         `gorm:"uniqueIndex;not null"`
	Type string         `gorm:"not null"`
}

func (VehicleModel) TableName() string {
	return "vehicle.vehicle_models"
}
