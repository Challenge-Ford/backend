-- Move tables back to vehicle schema
ALTER TABLE catalog.vehicle_model_year_colors SET SCHEMA vehicle;
ALTER TABLE catalog.vehicle_model_years SET SCHEMA vehicle;
ALTER TABLE catalog.vehicle_models SET SCHEMA vehicle;

-- Recreate unique index with old schema
DROP INDEX IF EXISTS catalog.idx_vehicle_model_years_unique;
CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicle_model_years_unique
    ON vehicle.vehicle_model_years (model_id, year);
