package vehiclerepository

import (
	"context"

	"gorm.io/gorm"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type GormModelRepository struct {
	db *gorm.DB
}

func NewGormModelRepository(db *gorm.DB) *GormModelRepository {
	return &GormModelRepository{db: db}
}

func (r *GormModelRepository) GetModelByID(ctx context.Context, id vehicledomain.VehicleModelID) (*vehicledomain.VehicleModel, error) {
	var model vehicledomain.VehicleModel
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *GormModelRepository) ListModels(ctx context.Context) ([]*vehicledomain.VehicleModel, error) {
	var models []*vehicledomain.VehicleModel
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *GormModelRepository) GetModelYearByID(ctx context.Context, id vehicledomain.VehicleModelYearID) (*vehicledomain.VehicleModelYear, error) {
	var year vehicledomain.VehicleModelYear
	err := r.db.WithContext(ctx).Preload("Model").First(&year, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &year, nil
}

func (r *GormModelRepository) ListModelYears(ctx context.Context, modelID vehicledomain.VehicleModelID) ([]*vehicledomain.VehicleModelYear, error) {
	var years []*vehicledomain.VehicleModelYear
	if err := r.db.WithContext(ctx).Where("model_id = ?", modelID).Order("year ASC").Find(&years).Error; err != nil {
		return nil, err
	}
	return years, nil
}
