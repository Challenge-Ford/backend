package devicerepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"torque/internal/core/pagination"
	devicedomain "torque/internal/modules/device/domain"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(ctx context.Context, page pagination.Page) ([]*devicedomain.Device, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM device.devices WHERE deleted_at IS NULL
	`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT d.id, d.name, d.vehicle_id, d.certificate_cn, d.certificate_sn,
		       d.created_at, d.updated_at, d.created_by, d.updated_by, d.deleted_at, d.deleted_by,
		       v.vin        AS vehicle_vin,
		       v.plate      AS vehicle_plate,
		       v.color      AS vehicle_color,
		       vm.name      AS vehicle_model_name,
		       vmy.year     AS vehicle_year,
		       vmy.model_url AS vehicle_model_url
		FROM device.devices d
		LEFT JOIN vehicle.vehicles v ON v.id = d.vehicle_id AND v.deleted_at IS NULL
		LEFT JOIN catalog.vehicle_model_years vmy ON vmy.id = v.model_year_id
		LEFT JOIN catalog.vehicle_models vm ON vm.id = vmy.model_id
		WHERE d.deleted_at IS NULL
		ORDER BY d.created_at DESC
		OFFSET $1 LIMIT $2
	`, page.Offset(), page.PerPage)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var devices []*devicedomain.Device
	for rows.Next() {
		d, err := scanDevice(rows)
		if err != nil {
			return nil, 0, err
		}
		devices = append(devices, d)
	}
	return devices, total, rows.Err()
}

func (r *Repository) Save(ctx context.Context, device *devicedomain.Device) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO device.devices (id, name, vehicle_id, certificate_cn, certificate_sn,
		                            created_at, updated_at, created_by, updated_by, deleted_at, deleted_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			name            = EXCLUDED.name,
			vehicle_id      = EXCLUDED.vehicle_id,
			certificate_cn  = EXCLUDED.certificate_cn,
			certificate_sn  = EXCLUDED.certificate_sn,
			updated_at      = EXCLUDED.updated_at,
			updated_by      = EXCLUDED.updated_by,
			deleted_at      = EXCLUDED.deleted_at,
			deleted_by      = EXCLUDED.deleted_by
	`,
		device.ID, device.Name, device.VehicleID,
		device.CertificateCN, device.CertificateSN,
		device.CreatedAt, device.UpdatedAt,
		device.CreatedBy, device.UpdatedBy,
		device.DeletedAt, device.DeletedBy,
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id devicedomain.DeviceID) (*devicedomain.Device, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, vehicle_id, certificate_cn, certificate_sn,
		       created_at, updated_at, created_by, updated_by, deleted_at, deleted_by
		FROM device.devices
		WHERE id = $1 AND deleted_at IS NULL
	`, id)

	return scanDeviceSimple(row)
}

// GetByName looks up a device by name including soft-deleted records.
// This is intentional: CreateDevice needs to find soft-deleted devices
// to re-issue certificates for re-created devices with the same name.
func (r *Repository) GetByName(ctx context.Context, name string) (*devicedomain.Device, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, vehicle_id, certificate_cn, certificate_sn,
		       created_at, updated_at, created_by, updated_by, deleted_at, deleted_by
		FROM device.devices
		WHERE name = $1
	`, name)

	return scanDeviceSimple(row)
}

func (r *Repository) GetByVehicleID(ctx context.Context, vehicleID uuid.UUID) (*devicedomain.Device, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, vehicle_id, certificate_cn, certificate_sn,
		       created_at, updated_at, created_by, updated_by, deleted_at, deleted_by
		FROM device.devices
		WHERE vehicle_id = $1 AND deleted_at IS NULL
	`, vehicleID)

	return scanDeviceSimple(row)
}

func (r *Repository) GetCommissionedByVIN(ctx context.Context, vin string) (*devicedomain.Device, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT d.id, d.name, d.vehicle_id, d.certificate_cn, d.certificate_sn,
		       d.created_at, d.updated_at, d.created_by, d.updated_by, d.deleted_at, d.deleted_by
		FROM device.devices d
		JOIN vehicle.vehicles v ON v.id = d.vehicle_id AND v.deleted_at IS NULL
		WHERE v.vin = $1 AND d.vehicle_id IS NOT NULL AND d.deleted_at IS NULL
	`, vin)

	return scanDeviceSimple(row)
}

type deviceScanner interface {
	Scan(dest ...any) error
}

func scanDevice(s deviceScanner) (*devicedomain.Device, error) {
	d := &devicedomain.Device{}
	err := s.Scan(
		&d.ID, &d.Name, &d.VehicleID, &d.CertificateCN, &d.CertificateSN,
		&d.CreatedAt, &d.UpdatedAt, &d.CreatedBy, &d.UpdatedBy, &d.DeletedAt, &d.DeletedBy,
		&d.VehicleVIN, &d.VehiclePlate, &d.VehicleColor,
		&d.VehicleModelName, &d.VehicleYear, &d.VehicleModelURL,
	)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func scanDeviceSimple(s deviceScanner) (*devicedomain.Device, error) {
	d := &devicedomain.Device{}
	err := s.Scan(
		&d.ID, &d.Name, &d.VehicleID, &d.CertificateCN, &d.CertificateSN,
		&d.CreatedAt, &d.UpdatedAt, &d.CreatedBy, &d.UpdatedBy, &d.DeletedAt, &d.DeletedBy,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}
