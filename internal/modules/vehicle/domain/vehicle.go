package vehicledomain

import (
	"github.com/google/uuid"
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

type Vehicle struct {
	ID         VehicleID `gorm:"type:uuid;primaryKey"`
	CustomerID uuid.UUID `gorm:"type:uuid;not null;index"`
	VIN        VIN       `gorm:"uniqueIndex;not null"`
	Plate      Plate     `gorm:"uniqueIndex;not null"`
	Model      string    `gorm:"not null"`
	Year       int       `gorm:"not null"`
	Color      Color     `gorm:"not null"`
	db.AuditableModel
}
