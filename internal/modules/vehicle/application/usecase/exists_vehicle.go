package vehicleusecase

import (
	"context"

	vehicledomain "torque/internal/modules/vehicle/domain"
)

type ExistsVehicleUseCase struct {
	repo vehicledomain.Repository
}

func NewExistsVehicle(repo vehicledomain.Repository) *ExistsVehicleUseCase {
	return &ExistsVehicleUseCase{repo: repo}
}

func (uc *ExistsVehicleUseCase) Execute(ctx context.Context, id vehicledomain.VehicleID) (bool, error) {
	return uc.repo.Exists(ctx, id)
}
