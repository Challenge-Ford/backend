package vehicleusecase

import (
	"context"

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

	// dtcMap returns true only for VINs with at least one open DTC.
	// Missing VINs in the map implicitly mean no active DTCs (false).
	dtcMap, err := uc.telemetryResolver.HasActiveDTCs(ctx, []string{string(vehicle.VIN)})
	if err != nil {
		return nil, apperr.Internal("failed to check active dtcs", err)
	}
	_, hasActive := dtcMap[string(vehicle.VIN)]
	out.HasActiveDTCs = hasActive

	return out, nil
}
