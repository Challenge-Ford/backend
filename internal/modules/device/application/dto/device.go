package devicedto

import (
	devicedomain "torque/internal/modules/device/domain"
)

type CreateDeviceInput struct {
	Name string `json:"name" validate:"required,device_name"`
}

type CommissionDeviceInput struct {
	VehicleID string `json:"vehicleId" validate:"required,uuid"`
}

type DeviceVehicleOutput struct {
	ID        string  `json:"id"`
	VIN       string  `json:"vin"`
	Plate     string  `json:"plate"`
	Color     string  `json:"color"`
	ModelName string  `json:"modelName"`
	Year      int     `json:"year"`
	ModelURL  *string `json:"modelUrl"`
}

type DeviceOutput struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Vehicle   *DeviceVehicleOutput `json:"vehicle"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
}

type CreateDeviceOutput struct {
	DeviceOutput
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

func ToDeviceOutput(d *devicedomain.Device) *DeviceOutput {
	var vehicle *DeviceVehicleOutput
	if d.VehicleID != nil && d.VehicleVIN != nil {
		vehicle = &DeviceVehicleOutput{
			ID:        d.VehicleID.String(),
			VIN:       *d.VehicleVIN,
			Plate:     strVal(d.VehiclePlate),
			Color:     strVal(d.VehicleColor),
			ModelName: strVal(d.VehicleModelName),
			ModelURL:  d.VehicleModelURL,
		}
		if d.VehicleYear != nil {
			vehicle.Year = *d.VehicleYear
		}
	}

	return &DeviceOutput{
		ID:        d.ID.String(),
		Name:      d.Name,
		Vehicle:   vehicle,
		CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
