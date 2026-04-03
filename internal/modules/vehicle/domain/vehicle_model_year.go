package vehicledomain

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"torque/internal/core/id"
)

type VehicleModelYearID uuid.UUID

func NewVehicleModelYearID() VehicleModelYearID {
	return VehicleModelYearID(id.New())
}

func (v VehicleModelYearID) String() string {
	return uuid.UUID(v).String()
}

func (v VehicleModelYearID) Value() (driver.Value, error) {
	return uuid.UUID(v).String(), nil
}

func (v *VehicleModelYearID) Scan(src any) error {
	switch val := src.(type) {
	case string:
		parsed, err := uuid.Parse(val)
		if err != nil {
			return err
		}
		*v = VehicleModelYearID(parsed)
		return nil
	case []byte:
		parsed, err := uuid.Parse(string(val))
		if err != nil {
			return err
		}
		*v = VehicleModelYearID(parsed)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", val)
	}
}

type VehicleModelYear struct {
	ID      VehicleModelYearID `gorm:"type:uuid;primaryKey"`
	ModelID VehicleModelID     `gorm:"type:uuid;not null;index"`
	Year    int                `gorm:"not null"`
	Model   *VehicleModel      `gorm:"foreignKey:ModelID"`
}

func (VehicleModelYear) TableName() string {
	return "vehicle.vehicle_model_years"
}
