
CREATE SCHEMA IF NOT EXISTS vehicle;
CREATE SCHEMA IF NOT EXISTS device;

CREATE TABLE IF NOT EXISTS vehicle.vehicle_models (
    id   UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS vehicle.vehicle_model_years (
    id         UUID PRIMARY KEY,
    model_id   UUID NOT NULL REFERENCES vehicle.vehicle_models(id),
    year       INT NOT NULL,
    model_url  TEXT
);

CREATE TABLE IF NOT EXISTS vehicle.vehicle_model_year_colors (
    id            UUID PRIMARY KEY,
    model_year_id UUID NOT NULL REFERENCES vehicle.vehicle_model_years(id),
    name          VARCHAR(100) NOT NULL,
    hex           VARCHAR(7) NOT NULL
);

CREATE TABLE IF NOT EXISTS vehicle.vehicles (
    id             UUID PRIMARY KEY,
    customer_id    UUID,
    model_year_id  UUID NOT NULL REFERENCES vehicle.vehicle_model_years(id),
    vin            VARCHAR(17) NOT NULL,
    plate          VARCHAR(7) NOT NULL,
    color          VARCHAR(7) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    updated_by     UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    deleted_at     TIMESTAMPTZ,
    deleted_by     UUID
);

CREATE TABLE IF NOT EXISTS device.devices (
    id              UUID PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    vehicle_id      UUID REFERENCES vehicle.vehicles(id),
    certificate_cn  VARCHAR(255) NOT NULL,
    certificate_sn  VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    updated_by      UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    deleted_at      TIMESTAMPTZ,
    deleted_by      UUID
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicle_model_years_unique
    ON vehicle.vehicle_model_years (model_id, year);
CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicles_vin_active
    ON vehicle.vehicles (vin) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicles_plate_active
    ON vehicle.vehicles (plate) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_name
    ON device.devices (name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_certificate_cn
    ON device.devices (certificate_cn);
CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_certificate_sn
    ON device.devices (certificate_sn);
