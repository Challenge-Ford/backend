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

type VehicleModelYearColorOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

type VehicleModelYearOutput struct {
	ID       string                        `json:"id"`
	ModelID  string                        `json:"modelId"`
	Year     int                           `json:"year"`
	ModelURL *string                       `json:"modelUrl"`
	Colors   []VehicleModelYearColorOutput `json:"colors"`
}

func ToVehicleModelOutput(m *vehicledomain.VehicleModel) *VehicleModelOutput {
	return &VehicleModelOutput{
		ID:   m.ID.String(),
		Name: m.Name,
		Type: m.Type,
	}
}

func ToVehicleModelYearOutput(y *vehicledomain.VehicleModelYear) *VehicleModelYearOutput {
	colors := make([]VehicleModelYearColorOutput, len(y.Colors))
	for i, c := range y.Colors {
		colors[i] = VehicleModelYearColorOutput{
			ID:   c.ID.String(),
			Name: c.Name,
			Hex:  c.Hex,
		}
	}
	return &VehicleModelYearOutput{
		ID:       y.ID.String(),
		ModelID:  y.ModelID.String(),
		Year:     y.Year,
		ModelURL: y.ModelURL,
		Colors:   colors,
	}
}
