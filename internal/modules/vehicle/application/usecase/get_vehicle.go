package vehicleusecase

import (
	"context"

	"torque/internal/core/apperr"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type GetVehicleUseCase struct {
	repo    vehicledomain.Repository
	dtcRepo telemetrydomain.DTCRepository
}

func NewGetVehicle(repo vehicledomain.Repository, dtcRepo telemetrydomain.DTCRepository) *GetVehicleUseCase {
	return &GetVehicleUseCase{repo: repo, dtcRepo: dtcRepo}
}

func (uc *GetVehicleUseCase) Execute(ctx context.Context, id vehicledomain.VehicleID) (*vehicledto.VehicleOutput, error) {
	vehicle, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle", err)
	}
	if vehicle == nil {
		return nil, apperr.NotFound("vehicle")
	}

	out := vehicledto.ToVehicleOutput(vehicle)

	dtcMap, err := uc.dtcRepo.HasActiveDTCs(ctx, []string{string(vehicle.VIN)})
	if err != nil {
		return nil, apperr.Internal("failed to check active dtcs", err)
	}
	out.HasActiveDTCs = dtcMap[string(vehicle.VIN)]

	return out, nil
}
