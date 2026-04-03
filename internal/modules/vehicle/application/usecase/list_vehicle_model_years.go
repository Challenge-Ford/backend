package vehicleusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	"torque/internal/core/pagination"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type ListVehicleModelYearsUseCase struct {
	repo vehicledomain.ModelRepository
}

func NewListVehicleModelYears(repo vehicledomain.ModelRepository) *ListVehicleModelYearsUseCase {
	return &ListVehicleModelYearsUseCase{repo: repo}
}

func (uc *ListVehicleModelYearsUseCase) Execute(ctx context.Context, modelID uuid.UUID) (*pagination.Result[*vehicledto.VehicleModelYearOutput], error) {
	model, err := uc.repo.GetModelByID(ctx, vehicledomain.VehicleModelID(modelID))
	if err != nil {
		return nil, apperr.Internal("failed to get model", err)
	}
	if model == nil {
		return nil, apperr.NotFound("vehicle model")
	}

	years, err := uc.repo.ListModelYears(ctx, vehicledomain.VehicleModelID(modelID))
	if err != nil {
		return nil, apperr.Internal("failed to list model years", err)
	}

	out := make([]*vehicledto.VehicleModelYearOutput, len(years))
	for i, y := range years {
		out[i] = vehicledto.ToVehicleModelYearOutput(y)
	}

	perPage := len(years)
	if perPage < 1 {
		perPage = 1
	}
	result := pagination.NewResult(out, pagination.Page{Page: 1, PerPage: perPage}, len(years))
	return &result, nil
}
