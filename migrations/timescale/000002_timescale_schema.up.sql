-- +migrate Up
CREATE TABLE IF NOT EXISTS telemetry_entries (
    time            TIMESTAMPTZ      NOT NULL,
    device_id       UUID             NOT NULL,
    vin             TEXT             NOT NULL,
    lat             DOUBLE PRECISION,
    lng             DOUBLE PRECISION,
    alt             DOUBLE PRECISION,
    gps_speed       DOUBLE PRECISION,
    heading         DOUBLE PRECISION,
    hdop            DOUBLE PRECISION,
    rpm             INTEGER,
    speed           INTEGER,
    coolant_temp    DOUBLE PRECISION,
    intake_temp     DOUBLE PRECISION,
    engine_load     DOUBLE PRECISION,
    throttle_pos    DOUBLE PRECISION,
    fuel_level      DOUBLE PRECISION,
    fuel_trim_short DOUBLE PRECISION,
    fuel_trim_long  DOUBLE PRECISION,
    maf             DOUBLE PRECISION,
    battery_voltage DOUBLE PRECISION,
    PRIMARY KEY (time, device_id)
);

SELECT create_hypertable('telemetry_entries', 'time', if_not_exists => true);

CREATE INDEX IF NOT EXISTS idx_telemetry_entries_vin_time ON telemetry_entries (vin, time DESC);

CREATE TABLE IF NOT EXISTS dtc_entries (
    time      TIMESTAMPTZ NOT NULL,
    device_id UUID        NOT NULL,
    vin       TEXT        NOT NULL,
    code      TEXT        NOT NULL,
    status    TEXT        NOT NULL
);

SELECT create_hypertable('dtc_entries', 'time', if_not_exists => true);

CREATE INDEX IF NOT EXISTS idx_dtc_entries_vin_code_time ON dtc_entries (vin, code, time DESC);
