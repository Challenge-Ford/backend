package devicerepository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"torque/internal/core/pagination"
	devicedomain "torque/internal/modules/device/domain"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(ctx context.Context, page pagination.Page) ([]*devicedomain.Device, int, error) {
	var devices []*devicedomain.Device
	var total int64

	if err := r.db.WithContext(ctx).Model(&devicedomain.Device{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Select(`device.devices.*,
			v.vin        AS vehicle_vin,
			v.plate      AS vehicle_plate,
			v.color      AS vehicle_color,
			vm.name      AS vehicle_model_name,
			vmy.year     AS vehicle_year,
			vmy.model_url AS vehicle_model_url`).
		Joins("LEFT JOIN vehicle.vehicles v ON v.id = device.devices.vehicle_id AND v.deleted_at IS NULL").
		Joins("LEFT JOIN vehicle.vehicle_model_years vmy ON vmy.id = v.model_year_id").
		Joins("LEFT JOIN vehicle.vehicle_models vm ON vm.id = vmy.model_id").
		Offset(page.Offset()).
		Limit(page.PerPage).
		Find(&devices).Error; err != nil {
		return nil, 0, err
	}

	return devices, int(total), nil
}

func (r *GormRepository) Save(ctx context.Context, device *devicedomain.Device) error {
	return r.db.WithContext(ctx).Save(device).Error
}

func (r *GormRepository) GetByID(ctx context.Context, id devicedomain.DeviceID) (*devicedomain.Device, error) {
	var device devicedomain.Device
	err := r.db.WithContext(ctx).First(&device, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *GormRepository) GetByName(ctx context.Context, name string) (*devicedomain.Device, error) {
	var device devicedomain.Device
	err := r.db.WithContext(ctx).Unscoped().First(&device, "name = ?", name).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &device, nil
}


func (r *GormRepository) GetByVehicleID(ctx context.Context, vehicleID uuid.UUID) (*devicedomain.Device, error) {
	var device devicedomain.Device
	err := r.db.WithContext(ctx).First(&device, "vehicle_id = ?", vehicleID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *GormRepository) GetCommissionedByVIN(ctx context.Context, vin string) (*devicedomain.Device, error) {
	var device devicedomain.Device
	err := r.db.WithContext(ctx).
		Joins("JOIN vehicle.vehicles v ON v.id = device.devices.vehicle_id AND v.deleted_at IS NULL").
		Where("v.vin = ? AND device.devices.vehicle_id IS NOT NULL AND device.devices.deleted_at IS NULL", vin).
		First(&device).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &device, nil
}
