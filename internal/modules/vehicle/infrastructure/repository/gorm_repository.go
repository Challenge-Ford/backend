package vehiclerepository

import (
	"context"

	"gorm.io/gorm"
	"torque/internal/core/pagination"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Save(ctx context.Context, vehicle *vehicledomain.Vehicle) error {
	return r.db.WithContext(ctx).Save(vehicle).Error
}

func (r *GormRepository) GetByID(ctx context.Context, id vehicledomain.VehicleID) (*vehicledomain.Vehicle, error) {
	var vehicle vehicledomain.Vehicle
	err := r.db.WithContext(ctx).Preload("ModelYear.Model").First(&vehicle, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *GormRepository) GetByVIN(ctx context.Context, vin vehicledomain.VIN) (*vehicledomain.Vehicle, error) {
	var vehicle vehicledomain.Vehicle
	err := r.db.WithContext(ctx).First(&vehicle, "vin = ?", vin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *GormRepository) GetByPlate(ctx context.Context, plate vehicledomain.Plate) (*vehicledomain.Vehicle, error) {
	var vehicle vehicledomain.Vehicle
	err := r.db.WithContext(ctx).First(&vehicle, "plate = ?", plate).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *GormRepository) List(ctx context.Context, page pagination.Page) ([]*vehicledomain.Vehicle, int, error) {
	var vehicles []*vehicledomain.Vehicle
	var total int64

	if err := r.db.WithContext(ctx).Model(&vehicledomain.Vehicle{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("ModelYear.Model").
		Offset(page.Offset()).
		Limit(page.PerPage).
		Find(&vehicles).Error; err != nil {
		return nil, 0, err
	}

	return vehicles, int(total), nil
}
