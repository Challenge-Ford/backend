package telemetryusecase

import (
	"context"

	"github.com/google/uuid"
	"torque/internal/core/apperr"
	telemetrydto "torque/internal/modules/telemetry/application/dto"
	telemetrydomain "torque/internal/modules/telemetry/domain"
)

type ListActiveDTCsUseCase struct {
	dtcRepo         telemetrydomain.DTCRepository
	catalogRepo     telemetrydomain.DTCCatalogRepository
	vehicleResolver telemetrydomain.VehicleResolver
}

func NewListActiveDTCs(dtcRepo telemetrydomain.DTCRepository, catalogRepo telemetrydomain.DTCCatalogRepository, vehicleResolver telemetrydomain.VehicleResolver) *ListActiveDTCsUseCase {
	return &ListActiveDTCsUseCase{dtcRepo: dtcRepo, catalogRepo: catalogRepo, vehicleResolver: vehicleResolver}
}

func (uc *ListActiveDTCsUseCase) Execute(ctx context.Context, vehicleID uuid.UUID) (*telemetrydto.DTCListOutput, error) {
	vin, err := uc.vehicleResolver.GetVINByID(ctx, vehicleID)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle", err)
	}
	if vin == "" {
		return nil, apperr.NotFound("vehicle")
	}

	modelYearID, err := uc.vehicleResolver.GetModelYearIDByVehicleID(ctx, vehicleID)
	if err != nil {
		return nil, apperr.Internal("failed to get vehicle model year", err)
	}

	active, err := uc.dtcRepo.ListActive(ctx, vin)
	if err != nil {
		return nil, apperr.Internal("failed to list active dtcs", err)
	}

	out := make([]*telemetrydto.DTCOutput, len(active))
	for i, d := range active {
		o := &telemetrydto.DTCOutput{
			Code: d.Code,
			Time: d.Time,
		}
		if cat, _ := uc.catalogRepo.GetWithEstimates(ctx, d.Code, modelYearID); cat != nil {
			o.Description = cat.Description
			o.System = cat.System
			o.Severity = cat.Severity
			o.RequiresStop = cat.RequiresStop
			o.CostMinCents = cat.CostMinCents
			o.CostMaxCents = cat.CostMaxCents
			o.TimeMin = cat.TimeMin
			o.TimeMax = cat.TimeMax
		}
		out[i] = o
	}
	return &telemetrydto.DTCListOutput{Data: out}, nil
}
