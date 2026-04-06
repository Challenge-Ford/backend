-- Test vehicles seed data
-- Creates one example of each vehicle model for local testing
-- Safe to run multiple times (idempotent via ON CONFLICT DO NOTHING)

-- Ford Ranger Raptor 2022, Azul Belize
INSERT INTO vehicle.vehicles (
    id, model_year_id, vin, plate, color,
    created_by, updated_by
) VALUES (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a22',
    '9B000000000000001',
    'ABC1D23',
    '#003A8F',
    '00000000-0000-0000-0000-000000000000',
    '00000000-0000-0000-0000-000000000000'
) ON CONFLICT (id) DO NOTHING;

-- Ford Territory 2023, Azul Profundo
INSERT INTO vehicle.vehicles (
    id, model_year_id, vin, plate, color,
    created_by, updated_by
) VALUES (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a23',
    '9B000000000000002',
    'DEF2G34',
    '#1A2F6C',
    '00000000-0000-0000-0000-000000000000',
    '00000000-0000-0000-0000-000000000000'
) ON CONFLICT (id) DO NOTHING;
