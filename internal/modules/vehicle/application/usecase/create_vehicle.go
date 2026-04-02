package vehicleusecase

import (
	"context"

	"torque/internal/core/apperr"
	"torque/internal/core/appctx"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type CreateVehicleUseCase struct {
	repo vehicledomain.Repository
}

func NewCreateVehicle(repo vehicledomain.Repository) *CreateVehicleUseCase {
	return &CreateVehicleUseCase{repo: repo}
}

func (uc *CreateVehicleUseCase) Execute(ctx context.Context, input vehicledto.CreateVehicleInput) (*vehicledto.VehicleOutput, error) {
	auth := appctx.MustGetAuth(ctx)

	vehicle := &vehicledomain.Vehicle{
		ID:         vehicledomain.NewVehicleID(),
		CustomerID: input.CustomerID,
		VIN:        vehicledomain.VIN(input.VIN),
		Plate:      vehicledomain.Plate(input.Plate),
		Model:      input.Model,
		Year:       input.Year,
		Color:      vehicledomain.Color(input.Color),
	}
	vehicle.CreatedBy = auth.UserID
	vehicle.UpdatedBy = auth.UserID

	if err := uc.repo.Save(ctx, vehicle); err != nil {
		return nil, apperr.Internal("failed to create vehicle", err)
	}

	return vehicledto.ToVehicleOutput(vehicle), nil
}
