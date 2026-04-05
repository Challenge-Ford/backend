package vehicledto

import (
	"github.com/google/uuid"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type CreateVehicleInput struct {
	CustomerID  *uuid.UUID `json:"customerId"  validate:"omitempty"`
	ModelYearID uuid.UUID  `json:"modelYearId" validate:"required"`
	VIN         string     `json:"vin"         validate:"required,vin"`
	Plate       string     `json:"plate"       validate:"required,plate"`
	Color       string     `json:"color"       validate:"required,hexcolor"`
}

type UpdateVehicleInput struct {
	ModelYearID *uuid.UUID `json:"modelYearId" validate:"omitempty"`
	Plate       string     `json:"plate"       validate:"omitempty,plate"`
	Color       string     `json:"color"       validate:"omitempty,hexcolor"`
}

type VehicleOutput struct {
	ID             string  `json:"id"`
	CustomerID     *string `json:"customerId"`
	ModelID        string  `json:"modelId"`
	ModelYearID    string  `json:"modelYearId"`
	ModelName      string  `json:"modelName"`
	Year           int     `json:"year"`
	ModelURL       *string `json:"modelUrl"`
	VIN            string  `json:"vin"`
	Plate          string  `json:"plate"`
	Color          string  `json:"color"`
	HasActiveDTCs  bool    `json:"hasActiveDtcs"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

func ToVehicleOutput(v *vehicledomain.Vehicle) *VehicleOutput {
	var customerID *string
	if v.CustomerID != nil {
		s := v.CustomerID.String()
		customerID = &s
	}

	var modelID, modelYearID, modelName string
	var year int
	var modelURL *string
	if v.ModelYear != nil {
		modelYearID = v.ModelYear.ID.String()
		year = v.ModelYear.Year
		modelID = v.ModelYear.ModelID.String()
		modelURL = v.ModelYear.ModelURL
		if v.ModelYear.Model != nil {
			modelName = v.ModelYear.Model.Name
		}
	}

	return &VehicleOutput{
		ID:          v.ID.String(),
		CustomerID:  customerID,
		ModelID:     modelID,
		ModelYearID: modelYearID,
		ModelName:   modelName,
		Year:        year,
		ModelURL:    modelURL,
		VIN:        string(v.VIN),
		Plate:      string(v.Plate),
		Color:      string(v.Color),
		CreatedAt:  v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  v.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
