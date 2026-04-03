package vehicleusecase

import (
	"context"

	"torque/internal/core/apperr"
	"torque/internal/core/pagination"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type ListVehicleModelsUseCase struct {
	repo vehicledomain.ModelRepository
}

func NewListVehicleModels(repo vehicledomain.ModelRepository) *ListVehicleModelsUseCase {
	return &ListVehicleModelsUseCase{repo: repo}
}

func (uc *ListVehicleModelsUseCase) Execute(ctx context.Context) (*pagination.Result[*vehicledto.VehicleModelOutput], error) {
	models, err := uc.repo.ListModels(ctx)
	if err != nil {
		return nil, apperr.Internal("failed to list vehicle models", err)
	}

	out := make([]*vehicledto.VehicleModelOutput, len(models))
	for i, m := range models {
		out[i] = vehicledto.ToVehicleModelOutput(m)
	}

	perPage := len(models)
	if perPage < 1 {
		perPage = 1
	}
	result := pagination.NewResult(out, pagination.Page{Page: 1, PerPage: perPage}, len(models))
	return &result, nil
}
