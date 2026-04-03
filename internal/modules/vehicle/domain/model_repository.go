package vehicledomain

import "context"

type ModelRepository interface {
	GetModelByID(ctx context.Context, id VehicleModelID) (*VehicleModel, error)
	ListModels(ctx context.Context) ([]*VehicleModel, error)

	GetModelYearByID(ctx context.Context, id VehicleModelYearID) (*VehicleModelYear, error)
	ListModelYears(ctx context.Context, modelID VehicleModelID) ([]*VehicleModelYear, error)
}
