package vehicledto

import vehicledomain "torque/internal/modules/vehicle/domain"

type CreateVehicleModelInput struct {
	Name string `json:"name" validate:"required"`
}

type CreateVehicleModelYearInput struct {
	Year int `json:"year" validate:"required,min=1886"`
}

type VehicleModelOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type VehicleModelYearOutput struct {
	ID      string `json:"id"`
	ModelID string `json:"modelId"`
	Year    int    `json:"year"`
}

func ToVehicleModelOutput(m *vehicledomain.VehicleModel) *VehicleModelOutput {
	return &VehicleModelOutput{
		ID:   m.ID.String(),
		Name: m.Name,
		Type: m.Type,
	}
}

func ToVehicleModelYearOutput(y *vehicledomain.VehicleModelYear) *VehicleModelYearOutput {
	return &VehicleModelYearOutput{
		ID:      y.ID.String(),
		ModelID: y.ModelID.String(),
		Year:    y.Year,
	}
}
