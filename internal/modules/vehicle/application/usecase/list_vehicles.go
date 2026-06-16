package vehicleusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	"torque/internal/core/pagination"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type ListVehiclesUseCase struct {
	repo              vehicledomain.Repository
	telemetryResolver vehicledomain.TelemetryResolver
}

func NewListVehicles(repo vehicledomain.Repository, telemetryResolver vehicledomain.TelemetryResolver) *ListVehiclesUseCase {
	return &ListVehiclesUseCase{repo: repo, telemetryResolver: telemetryResolver}
}

func (uc *ListVehiclesUseCase) Execute(ctx context.Context, page pagination.Page) (*pagination.Result[*vehicledto.VehicleOutput], error) {
	page.Normalize(pagination.DefaultConfig)

	vehicles, total, err := uc.repo.List(ctx, page)
	if err != nil {
		return nil, apperr.Internal("failed to list vehicles", err)
	}

	vehicleIDs := make([]uuid.UUID, len(vehicles))
	for i, v := range vehicles {
		vehicleIDs[i] = uuid.UUID(v.ID)
	}

	dtcMap, err := uc.telemetryResolver.HasActiveDTCs(ctx, vehicleIDs)
	if err != nil {
		return nil, apperr.Internal("failed to check active dtcs", err)
	}

	output := make([]*vehicledto.VehicleOutput, len(vehicles))
	for i, v := range vehicles {
		out := vehicledto.ToVehicleOutput(v)
		out.HasActiveDTCs = dtcMap[uuid.UUID(v.ID)]
		output[i] = out
	}

	result := pagination.NewResult(output, page, total)
	return &result, nil
}
