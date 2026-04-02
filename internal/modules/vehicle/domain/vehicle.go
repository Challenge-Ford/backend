package vehicledomain

import (
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

func (v *Vehicle) Delete(byUser uuid.UUID) {
	now := time.Now()
	v.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	v.DeletedBy = &byUser
}
