package vehicleusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type GetVehicleUseCase struct {
	repo              vehicledomain.Repository
	telemetryResolver vehicledomain.TelemetryResolver
}

func NewGetVehicle(repo vehicledomain.Repository, telemetryResolver vehicledomain.TelemetryResolver) *GetVehicleUseCase {
	return &GetVehicleUseCase{repo: repo, telemetryResolver: telemetryResolver}
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

	// dtcMap returns true only for vehicles with at least one open DTC.
	// Missing vehicle IDs in the map implicitly mean no active DTCs (false).
	dtcMap, err := uc.telemetryResolver.HasActiveDTCs(ctx, []uuid.UUID{uuid.UUID(vehicle.ID)})
	if err != nil {
		return nil, apperr.Internal("failed to check active dtcs", err)
	}
	_, hasActive := dtcMap[uuid.UUID(vehicle.ID)]
	out.HasActiveDTCs = hasActive

	return out, nil
}
