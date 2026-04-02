package vehicleusecase

import (
	"context"

	"torque/internal/core/apperr"
	"torque/internal/core/appctx"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type UpdateVehicleUseCase struct {
	repo vehicledomain.Repository
}

func NewUpdateVehicle(repo vehicledomain.Repository) *UpdateVehicleUseCase {
	return &UpdateVehicleUseCase{repo: repo}
}

func (uc *UpdateVehicleUseCase) Execute(ctx context.Context, id vehicledomain.VehicleID, input vehicledto.UpdateVehicleInput) (*vehicledto.VehicleOutput, error) {
	auth := appctx.MustGetAuth(ctx)

	vehicle, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle", err)
	}
	if vehicle == nil {
		return nil, apperr.NotFound("vehicle")
	}

	if input.Plate != "" {
		plate := vehicledomain.Plate(input.Plate)
		if err := plate.Validate(); err != nil {
			return nil, apperr.BadRequest(err.Error())
		}
		vehicle.Plate = plate
	}

	if input.Model != "" {
		vehicle.Model = input.Model
	}

	if input.Year != 0 {
		vehicle.Year = input.Year
	}

	if input.Color != "" {
		color := vehicledomain.Color(input.Color)
		if err := color.Validate(); err != nil {
			return nil, apperr.BadRequest(err.Error())
		}
		vehicle.Color = color
	}

	vehicle.UpdatedBy = auth.UserID

	if err := uc.repo.Save(ctx, vehicle); err != nil {
		return nil, apperr.Internal("failed to update vehicle", err)
	}

	return vehicledto.ToVehicleOutput(vehicle), nil
}
