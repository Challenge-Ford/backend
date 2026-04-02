package vehicledto

import (
	"github.com/google/uuid"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type CreateVehicleInput struct {
	CustomerID uuid.UUID `json:"customerId" validate:"required"`
	VIN        string    `json:"vin"        validate:"required,vin"`
	Plate      string    `json:"plate"      validate:"required,plate"`
	Model      string    `json:"model"      validate:"required"`
	Year       int       `json:"year"       validate:"required,min=1886"`
	Color      string    `json:"color"      validate:"required,hexcolor"`
}

type UpdateVehicleInput struct {
	Plate string `json:"plate" validate:"omitempty,plate"`
	Model string `json:"model" validate:"omitempty"`
	Year  int    `json:"year"  validate:"omitempty,min=1886"`
	Color string `json:"color" validate:"omitempty,hexcolor"`
}

type VehicleOutput struct {
	ID         string `json:"id"`
	CustomerID string `json:"customerId"`
	VIN        string `json:"vin"`
	Plate      string `json:"plate"`
	Model      string `json:"model"`
	Year       int    `json:"year"`
	Color      string `json:"color"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func ToVehicleOutput(v *vehicledomain.Vehicle) *VehicleOutput {
	return &VehicleOutput{
		ID:         v.ID.String(),
		CustomerID: v.CustomerID.String(),
		VIN:        string(v.VIN),
		Plate:      string(v.Plate),
		Model:      v.Model,
		Year:       v.Year,
		Color:      string(v.Color),
		CreatedAt:  v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  v.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
