package vehicleusecase

import (
	"context"

	"torque/internal/core/apperr"
	"torque/internal/core/pagination"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type ListVehiclesUseCase struct {
	repo vehicledomain.Repository
}

func NewListVehicles(repo vehicledomain.Repository) *ListVehiclesUseCase {
	return &ListVehiclesUseCase{repo: repo}
}

func (uc *ListVehiclesUseCase) Execute(ctx context.Context, page pagination.Page) (*pagination.Result[*vehicledto.VehicleOutput], error) {
	page.Normalize(pagination.DefaultConfig)

	vehicles, total, err := uc.repo.List(ctx, page)
	if err != nil {
		return nil, apperr.Internal("failed to list vehicles", err)
	}

	output := make([]*vehicledto.VehicleOutput, len(vehicles))
	for i, v := range vehicles {
		output[i] = vehicledto.ToVehicleOutput(v)
	}

	result := pagination.NewResult(output, page, total)
	return &result, nil
}
