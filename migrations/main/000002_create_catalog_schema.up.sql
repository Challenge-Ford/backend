CREATE SCHEMA IF NOT EXISTS catalog;

-- Move catalog tables from vehicle schema to catalog schema
ALTER TABLE vehicle.vehicle_model_year_colors SET SCHEMA catalog;
ALTER TABLE vehicle.vehicle_model_years SET SCHEMA catalog;
ALTER TABLE vehicle.vehicle_models SET SCHEMA catalog;

-- Recreate unique index with new schema
DROP INDEX IF EXISTS vehicle.idx_vehicle_model_years_unique;
CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicle_model_years_unique
    ON catalog.vehicle_model_years (model_id, year);
