package vehicleusecase

import (
	"context"

	"torque/internal/core/appctx"
	"torque/internal/core/apperr"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"

	"github.com/go-playground/validator/v10"
)

type CreateVehicleUseCase struct {
	repo     vehicledomain.Repository
	validate *validator.Validate
}

func NewCreateVehicle(repo vehicledomain.Repository, validate *validator.Validate) *CreateVehicleUseCase {
	return &CreateVehicleUseCase{repo: repo, validate: validate}
}

func (uc *CreateVehicleUseCase) Execute(ctx context.Context, input vehicledto.CreateVehicleInput) (*vehicledto.VehicleOutput, error) {
	auth := appctx.MustGetAuth(ctx)

	if err := uc.validate.Struct(input); err != nil {
		return nil, apperr.FromValidatorErr(err)
	}

	if existing, err := uc.repo.GetByVIN(ctx, vehicledomain.VIN(input.VIN)); err != nil {
		return nil, apperr.Internal("failed to check VIN", err)
	} else if existing != nil {
		return nil, apperr.Conflict("vehicle with this VIN already exists")
	}

	if existing, err := uc.repo.GetByPlate(ctx, vehicledomain.Plate(input.Plate)); err != nil {
		return nil, apperr.Internal("failed to check plate", err)
	} else if existing != nil {
		return nil, apperr.Conflict("vehicle with this plate already exists")
	}

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
