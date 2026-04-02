package vehicledto

import (
	"github.com/google/uuid"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type CreateVehicleInput struct {
	CustomerID uuid.UUID `validate:"required"`
	VIN        string    `validate:"required,vin"`
	Plate      string    `validate:"required,plate"`
	Model      string    `validate:"required"`
	Year       int       `validate:"required,min=1886"`
	Color      string    `validate:"required,hexcolor"`
}

type UpdateVehicleInput struct {
	Plate string `validate:"omitempty,plate"`
	Model string `validate:"omitempty"`
	Year  int    `validate:"omitempty,min=1886"`
	Color string `validate:"omitempty,hexcolor"`
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
