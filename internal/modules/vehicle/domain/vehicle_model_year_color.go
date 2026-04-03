package vehicledomain

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"torque/internal/core/id"
)

type VehicleModelYearColorID uuid.UUID

func NewVehicleModelYearColorID() VehicleModelYearColorID {
	return VehicleModelYearColorID(id.New())
}

func (v VehicleModelYearColorID) String() string {
	return uuid.UUID(v).String()
}

func (v VehicleModelYearColorID) Value() (driver.Value, error) {
	return uuid.UUID(v).String(), nil
}

func (v *VehicleModelYearColorID) Scan(src any) error {
	switch val := src.(type) {
	case string:
		parsed, err := uuid.Parse(val)
		if err != nil {
			return err
		}
		*v = VehicleModelYearColorID(parsed)
		return nil
	case []byte:
		parsed, err := uuid.Parse(string(val))
		if err != nil {
			return err
		}
		*v = VehicleModelYearColorID(parsed)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", val)
	}
}

type VehicleModelYearColor struct {
	ID          VehicleModelYearColorID `gorm:"type:uuid;primaryKey"`
	ModelYearID VehicleModelYearID      `gorm:"type:uuid;not null;index"`
	Name        string                  `gorm:"type:varchar(100);not null"`
	Hex         string                  `gorm:"type:varchar(7);not null"`
}

func (VehicleModelYearColor) TableName() string {
	return "vehicle.vehicle_model_year_colors"
}
