CREATE TABLE IF NOT EXISTS catalog.dtc_catalog (
    code            TEXT PRIMARY KEY,
    description     TEXT NOT NULL,
    system          TEXT,
    severity        TEXT NOT NULL,
    requires_stop   BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS catalog.dtc_vehicle_estimates (
    dtc_code         TEXT NOT NULL REFERENCES catalog.dtc_catalog(code) ON DELETE CASCADE,
    model_year_id    UUID NOT NULL REFERENCES catalog.vehicle_model_years(id) ON DELETE CASCADE,
    cost_min_cents   INT,
    cost_max_cents   INT,
    time_min         INT,
    time_max         INT,
    PRIMARY KEY (dtc_code, model_year_id)
);
