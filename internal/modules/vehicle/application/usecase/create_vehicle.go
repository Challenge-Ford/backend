package vehicleusecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"torque/internal/core/appctx"
	"torque/internal/core/apperr"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type CreateVehicleUseCase struct {
	repo      vehicledomain.Repository
	modelRepo vehicledomain.ModelRepository
	validate  *validator.Validate
}

func NewCreateVehicle(repo vehicledomain.Repository, modelRepo vehicledomain.ModelRepository, validate *validator.Validate) *CreateVehicleUseCase {
	return &CreateVehicleUseCase{repo: repo, modelRepo: modelRepo, validate: validate}
}

func (uc *CreateVehicleUseCase) Execute(ctx context.Context, input vehicledto.CreateVehicleInput) (*vehicledto.VehicleOutput, error) {
	auth := appctx.MustGetAuth(ctx)

	if err := uc.validate.Struct(input); err != nil {
		return nil, apperr.FromValidatorErr(err)
	}

	modelYear, err := uc.modelRepo.GetModelYearByID(ctx, vehicledomain.VehicleModelYearID(input.ModelYearID))
	if err != nil {
		return nil, apperr.Internal("failed to get model year", err)
	}
	if modelYear == nil {
		return nil, apperr.NotFound("vehicle model year")
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
		ID:          vehicledomain.NewVehicleID(),
		CustomerID:  input.CustomerID,
		ModelYearID: vehicledomain.VehicleModelYearID(input.ModelYearID),
		ModelYear:   modelYear,
		VIN:         vehicledomain.VIN(input.VIN),
		Plate:       vehicledomain.Plate(input.Plate),
		Color:       vehicledomain.Color(input.Color),
	}
	vehicle.CreatedBy = auth.UserID
	vehicle.UpdatedBy = auth.UserID

	if err := uc.repo.Save(ctx, vehicle); err != nil {
		return nil, apperr.Internal("failed to create vehicle", err)
	}

	return vehicledto.ToVehicleOutput(vehicle), nil
}
