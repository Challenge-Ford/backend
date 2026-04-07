-- DTC catalog seed data
-- Common OBD-II diagnostic trouble codes with descriptions
-- Safe to run multiple times (idempotent via ON CONFLICT DO NOTHING)

INSERT INTO catalog.dtc_catalog (code, description, system, severity, requires_stop) VALUES
  -- Engine misfire / combustion
  ('P0300', 'Random/Multiple Cylinder Misfire Detected', 'Engine', 'high', true),
  ('P0301', 'Cylinder 1 Misfire Detected', 'Engine', 'high', true),
  ('P0302', 'Cylinder 2 Misfire Detected', 'Engine', 'high', true),
  ('P0303', 'Cylinder 3 Misfire Detected', 'Engine', 'high', true),
  ('P0304', 'Cylinder 4 Misfire Detected', 'Engine', 'high', true),

  -- Fuel / air metering
  ('P0171', 'System Too Lean (Bank 1)', 'Engine', 'medium', false),
  ('P0172', 'System Too Rich (Bank 1)', 'Engine', 'medium', false),
  ('P0174', 'System Too Lean (Bank 2)', 'Engine', 'medium', false),
  ('P0175', 'System Too Rich (Bank 2)', 'Engine', 'medium', false),

  -- Catalyst / emissions
  ('P0420', 'Catalyst System Efficiency Below Threshold (Bank 1)', 'Emissions', 'medium', false),
  ('P0430', 'Catalyst System Efficiency Below Threshold (Bank 2)', 'Emissions', 'medium', false),

  -- EGR / EVAP
  ('P0401', 'Exhaust Gas Recirculation Flow Insufficient', 'Emissions', 'low', false),
  ('P0440', 'Evaporative Emission Control System Malfunction', 'Emissions', 'low', false),
  ('P0442', 'Evaporative Emission Control System Leak Detected (Small)', 'Emissions', 'low', false),
  ('P0455', 'Evaporative Emission Control System Leak Detected (Large)', 'Emissions', 'medium', false),

  -- Oxygen sensor
  ('P0130', 'O2 Sensor Circuit Malfunction (Bank 1, Sensor 1)', 'Engine', 'medium', false),
  ('P0135', 'O2 Sensor Heater Circuit Malfunction (Bank 1, Sensor 1)', 'Engine', 'medium', false),
  ('P0141', 'O2 Sensor Heater Circuit Malfunction (Bank 1, Sensor 2)', 'Engine', 'low', false),

  -- Ignition / electrical
  ('P0340', 'Camshaft Position Sensor Circuit Malfunction', 'Engine', 'high', true),
  ('P0351', 'Ignition Coil A Primary/Secondary Circuit Malfunction', 'Engine', 'high', true),

  -- Temperature / cooling
  ('P0217', 'Engine Overheating Condition', 'Cooling', 'critical', true),
  ('P0128', 'Coolant Thermostat (Coolant Temperature Below Thermostat Regulating Temperature)', 'Cooling', 'medium', false),

  -- Oil pressure
  ('P0520', 'Engine Oil Pressure Sensor/Switch Circuit Malfunction', 'Engine', 'high', true),
  ('P0522', 'Engine Oil Pressure Low', 'Engine', 'critical', true),

  -- Transmission
  ('P0700', 'Transmission Control System Malfunction', 'Transmission', 'high', false),
  ('P0730', 'Incorrect Gear Ratio', 'Transmission', 'high', true),

  -- Battery / charging
  ('P0562', 'System Voltage Low', 'Electrical', 'medium', false),
  ('P0563', 'System Voltage High', 'Electrical', 'medium', false),

  -- Throttle / pedal
  ('P0120', 'Throttle/Pedal Position Sensor/Switch A Circuit Malfunction', 'Engine', 'high', true),
  ('P0220', 'Throttle/Pedal Position Sensor/Switch B Circuit Malfunction', 'Engine', 'high', true)
ON CONFLICT (code) DO NOTHING;

-- Cost/time estimates for all model years
-- Ranger Raptor (2019-2022, 2024)
INSERT INTO catalog.dtc_vehicle_estimates (dtc_code, model_year_id, cost_min_cents, cost_max_cents, time_min, time_max)
SELECT dtc_code, vmy.id, cost_min_cents, cost_max_cents, time_min, time_max
FROM (VALUES
  ('P0300', 30000, 150000, 60, 240),
  ('P0301', 20000, 80000, 30, 120),
  ('P0420', 100000, 300000, 120, 360),
  ('P0171', 15000, 60000, 30, 120),
  ('P0217', 50000, 200000, 120, 480),
  ('P0522', 20000, 80000, 30, 90),
  ('P0700', 50000, 500000, 120, 600),
  ('P0128', 10000, 40000, 30, 60)
) AS estimates(dtc_code, cost_min_cents, cost_max_cents, time_min, time_max)
CROSS JOIN catalog.vehicle_model_years vmy
JOIN catalog.vehicle_models vm ON vm.id = vmy.model_id
WHERE vm.name = 'Ranger Raptor'
ON CONFLICT (dtc_code, model_year_id) DO NOTHING;

-- Territory (2020-2024)
INSERT INTO catalog.dtc_vehicle_estimates (dtc_code, model_year_id, cost_min_cents, cost_max_cents, time_min, time_max)
SELECT dtc_code, vmy.id, cost_min_cents, cost_max_cents, time_min, time_max
FROM (VALUES
  ('P0300', 35000, 160000, 60, 240),
  ('P0301', 25000, 85000, 30, 120),
  ('P0420', 120000, 350000, 120, 360),
  ('P0171', 18000, 65000, 30, 120),
  ('P0217', 60000, 220000, 120, 480),
  ('P0522', 22000, 85000, 30, 90),
  ('P0700', 60000, 520000, 120, 600),
  ('P0128', 12000, 45000, 30, 60)
) AS estimates(dtc_code, cost_min_cents, cost_max_cents, time_min, time_max)
CROSS JOIN catalog.vehicle_model_years vmy
JOIN catalog.vehicle_models vm ON vm.id = vmy.model_id
WHERE vm.name = 'Territory'
ON CONFLICT (dtc_code, model_year_id) DO NOTHING;
