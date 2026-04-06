-- +migrate Down
DROP TABLE IF EXISTS device.devices;
DROP TABLE IF EXISTS vehicle.vehicles;
DROP TABLE IF EXISTS vehicle.vehicle_model_year_colors;
DROP TABLE IF EXISTS vehicle.vehicle_model_years;
DROP TABLE IF EXISTS vehicle.vehicle_models;
DROP SCHEMA IF EXISTS device;
DROP SCHEMA IF EXISTS vehicle;
