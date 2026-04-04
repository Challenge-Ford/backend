CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS telemetry (
    time             TIMESTAMPTZ      NOT NULL,
    device_id        UUID             NOT NULL,
    vin              TEXT             NOT NULL,
    -- GPS (stored but not exposed via API)
    lat              DOUBLE PRECISION,
    lng              DOUBLE PRECISION,
    alt              DOUBLE PRECISION,
    gps_speed        DOUBLE PRECISION,
    heading          DOUBLE PRECISION,
    hdop             DOUBLE PRECISION,
    -- OBD
    rpm              INTEGER,
    speed            INTEGER,
    coolant_temp     DOUBLE PRECISION,
    intake_temp      DOUBLE PRECISION,
    engine_load      DOUBLE PRECISION,
    throttle_pos     DOUBLE PRECISION,
    fuel_level       DOUBLE PRECISION,
    fuel_trim_short  DOUBLE PRECISION,
    fuel_trim_long   DOUBLE PRECISION,
    maf              DOUBLE PRECISION,
    battery_voltage  DOUBLE PRECISION
);

SELECT create_hypertable('telemetry', by_range('time'), if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_telemetry_vin_time      ON telemetry (vin, time DESC);
CREATE INDEX IF NOT EXISTS idx_telemetry_device_time   ON telemetry (device_id, time DESC);

CREATE TABLE IF NOT EXISTS dtc_events (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id   UUID         NOT NULL,
    vin         TEXT         NOT NULL,
    code        TEXT         NOT NULL,
    opened_at   TIMESTAMPTZ  NOT NULL,
    closed_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_dtc_events_vin          ON dtc_events (vin, opened_at DESC);
CREATE INDEX IF NOT EXISTS idx_dtc_events_device_active ON dtc_events (device_id, code) WHERE closed_at IS NULL;
