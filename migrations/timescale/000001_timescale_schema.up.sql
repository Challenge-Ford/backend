CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS vehicle_state_message_ids (
    message_id  UUID        PRIMARY KEY,
    observed_at TIMESTAMPTZ NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS vehicle_state_observations (
    observed_at TIMESTAMPTZ NOT NULL,
    message_id  UUID        NOT NULL,
    device_id   UUID        NOT NULL,
    vehicle_id  UUID        NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    state       JSONB       NOT NULL,
    observation JSONB       NOT NULL DEFAULT '{}'::jsonb,
    raw_payload JSONB       NOT NULL,
    PRIMARY KEY (observed_at, message_id)
);

SELECT create_hypertable('vehicle_state_observations', 'observed_at', if_not_exists => true);

CREATE INDEX IF NOT EXISTS idx_vehicle_state_observations_vehicle_time
    ON vehicle_state_observations (vehicle_id, observed_at DESC);

CREATE INDEX IF NOT EXISTS idx_vehicle_state_observations_diagnostics
    ON vehicle_state_observations USING GIN ((state -> 'diagnostics'));
