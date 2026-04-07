-- Vehicle models and colors seed data
-- Safe to run multiple times (idempotent via ON CONFLICT DO NOTHING)

-- Ford Ranger Raptor
INSERT INTO catalog.vehicle_models (id, name, type)
VALUES ('1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a1a', 'Ranger Raptor', 'pickup')
ON CONFLICT (name) DO NOTHING;

-- 2019–2022 Ranger Raptor (same 3D model)
INSERT INTO catalog.vehicle_model_years (id, model_id, year, model_url) VALUES
  ('1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a19', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a1a', 2019, 'http://localhost:9000/vehicle-models/ranger_raptor.glb'),
  ('1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a20', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a1a', 2020, 'http://localhost:9000/vehicle-models/ranger_raptor.glb'),
  ('1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a21', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a1a', 2021, 'http://localhost:9000/vehicle-models/ranger_raptor.glb'),
  ('1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a22', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a1a', 2022, 'http://localhost:9000/vehicle-models/ranger_raptor.glb')
ON CONFLICT (model_id, year) DO NOTHING;

-- 2024 Ranger Raptor
INSERT INTO catalog.vehicle_model_years (id, model_id, year, model_url)
VALUES ('1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a24', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a1a', 2024, 'http://localhost:9000/vehicle-models/ranger_raptor.glb')
ON CONFLICT (model_id, year) DO NOTHING;

-- 2019–2022 colors
INSERT INTO catalog.vehicle_model_year_colors (id, model_year_id, name, hex) VALUES
  ('00000000-0019-0000-0001-000000000001', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a19', 'Azul Belize', '#003A8F'),
  ('00000000-0019-0000-0001-000000000002', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a19', 'Laranja Saara', '#D94F00'),
  ('00000000-0019-0000-0001-000000000003', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a19', 'Preto Astúrias', '#0B0B0B'),
  ('00000000-0019-0000-0001-000000000004', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a19', 'Branco Nevasca', '#F4F4F4'),
  ('00000000-0019-0000-0001-000000000005', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a19', 'Cinza Diamantina', '#6E7072'),
  ('00000000-0020-0000-0001-000000000001', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a20', 'Azul Belize', '#003A8F'),
  ('00000000-0020-0000-0001-000000000002', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a20', 'Laranja Saara', '#D94F00'),
  ('00000000-0020-0000-0001-000000000003', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a20', 'Preto Astúrias', '#0B0B0B'),
  ('00000000-0020-0000-0001-000000000004', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a20', 'Branco Nevasca', '#F4F4F4'),
  ('00000000-0020-0000-0001-000000000005', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a20', 'Cinza Diamantina', '#6E7072'),
  ('00000000-0021-0000-0001-000000000001', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a21', 'Azul Belize', '#003A8F'),
  ('00000000-0021-0000-0001-000000000002', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a21', 'Laranja Saara', '#D94F00'),
  ('00000000-0021-0000-0001-000000000003', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a21', 'Preto Astúrias', '#0B0B0B'),
  ('00000000-0021-0000-0001-000000000004', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a21', 'Branco Nevasca', '#F4F4F4'),
  ('00000000-0021-0000-0001-000000000005', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a21', 'Cinza Diamantina', '#6E7072'),
  ('00000000-0022-0000-0001-000000000001', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a22', 'Azul Belize', '#003A8F'),
  ('00000000-0022-0000-0001-000000000002', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a22', 'Laranja Saara', '#D94F00'),
  ('00000000-0022-0000-0001-000000000003', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a22', 'Preto Astúrias', '#0B0B0B'),
  ('00000000-0022-0000-0001-000000000004', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a22', 'Branco Nevasca', '#F4F4F4'),
  ('00000000-0022-0000-0001-000000000005', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a22', 'Cinza Diamantina', '#6E7072')
ON CONFLICT (id) DO NOTHING;

-- 2024 Ranger Raptor (different colors)
INSERT INTO catalog.vehicle_model_years (id, model_id, year, model_url)
VALUES ('1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a24', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a1a', 2024, 'http://localhost:9000/vehicle-models/ranger_raptor.glb')
ON CONFLICT (model_id, year) DO NOTHING;

INSERT INTO catalog.vehicle_model_year_colors (id, model_year_id, name, hex) VALUES
  ('00000000-0024-0000-0001-000000000001', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a24', 'Azul Belize', '#003A8F'),
  ('00000000-0024-0000-0001-000000000002', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a24', 'Laranja Saara (Code Orange)', '#FF5A1F'),
  ('00000000-0024-0000-0001-000000000003', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a24', 'Preto Astúria', '#0B0B0B'),
  ('00000000-0024-0000-0001-000000000004', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a24', 'Branco Nevasca', '#F4F4F4'),
  ('00000000-0024-0000-0001-000000000005', '1a1a1a1a-1a1a-1a1a-1a1a-1a1a1a1a1a24', 'Cinza Interlagos', '#5F6366')
ON CONFLICT (id) DO NOTHING;

-- Ford Territory
INSERT INTO catalog.vehicle_models (id, name, type)
VALUES ('2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a2a', 'Territory', 'suv')
ON CONFLICT (name) DO NOTHING;

-- 2020–2022 Territory
INSERT INTO catalog.vehicle_model_years (id, model_id, year, model_url) VALUES
  ('2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a20', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a2a', 2020, 'http://localhost:9000/vehicle-models/generic_suv.glb'),
  ('2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a21', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a2a', 2021, 'http://localhost:9000/vehicle-models/generic_suv.glb'),
  ('2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a22', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a2a', 2022, 'http://localhost:9000/vehicle-models/generic_suv.glb')
ON CONFLICT (model_id, year) DO NOTHING;

INSERT INTO catalog.vehicle_model_year_colors (id, model_year_id, name, hex) VALUES
  ('00000000-0020-0000-0002-000000000001', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a20', 'Azul Santorini', '#1F3F8C'),
  ('00000000-0020-0000-0002-000000000002', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a20', 'Branco Bariloche', '#F5F5F5'),
  ('00000000-0020-0000-0002-000000000003', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a20', 'Prata Maiorca', '#A6A9AD'),
  ('00000000-0020-0000-0002-000000000004', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a20', 'Marrom Roma', '#5A3E2B'),
  ('00000000-0020-0000-0002-000000000005', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a20', 'Preto Toronto', '#0A0A0A'),
  ('00000000-0021-0000-0002-000000000001', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a21', 'Azul Santorini', '#1F3F8C'),
  ('00000000-0021-0000-0002-000000000002', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a21', 'Branco Bariloche', '#F5F5F5'),
  ('00000000-0021-0000-0002-000000000003', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a21', 'Prata Maiorca', '#A6A9AD'),
  ('00000000-0021-0000-0002-000000000004', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a21', 'Marrom Roma', '#5A3E2B'),
  ('00000000-0021-0000-0002-000000000005', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a21', 'Preto Toronto', '#0A0A0A'),
  ('00000000-0022-0000-0002-000000000001', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a22', 'Azul Santorini', '#1F3F8C'),
  ('00000000-0022-0000-0002-000000000002', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a22', 'Branco Bariloche', '#F5F5F5'),
  ('00000000-0022-0000-0002-000000000003', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a22', 'Prata Maiorca', '#A6A9AD'),
  ('00000000-0022-0000-0002-000000000004', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a22', 'Marrom Roma', '#5A3E2B'),
  ('00000000-0022-0000-0002-000000000005', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a22', 'Preto Toronto', '#0A0A0A')
ON CONFLICT (id) DO NOTHING;

-- 2023–2024 Territory
INSERT INTO catalog.vehicle_model_years (id, model_id, year, model_url) VALUES
  ('2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a23', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a2a', 2023, 'http://localhost:9000/vehicle-models/generic_suv.glb'),
  ('2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a24', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a2a', 2024, 'http://localhost:9000/vehicle-models/generic_suv.glb')
ON CONFLICT (model_id, year) DO NOTHING;

INSERT INTO catalog.vehicle_model_year_colors (id, model_year_id, name, hex) VALUES
  ('00000000-0023-0000-0002-000000000001', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a23', 'Azul Profundo', '#1A2F6C'),
  ('00000000-0023-0000-0002-000000000002', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a23', 'Cinza Catar', '#7A7F85'),
  ('00000000-0023-0000-0002-000000000003', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a23', 'Cinza Dover', '#5C6064'),
  ('00000000-0023-0000-0002-000000000004', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a23', 'Branco Bariloche', '#F5F5F5'),
  ('00000000-0023-0000-0002-000000000005', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a23', 'Verde Oásis', '#3E5F4A'),
  ('00000000-0023-0000-0002-000000000006', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a23', 'Preto Toronto', '#0A0A0A'),
  ('00000000-0024-0000-0002-000000000001', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a24', 'Azul Profundo', '#1A2F6C'),
  ('00000000-0024-0000-0002-000000000002', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a24', 'Cinza Catar', '#7A7F85'),
  ('00000000-0024-0000-0002-000000000003', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a24', 'Cinza Dover', '#5C6064'),
  ('00000000-0024-0000-0002-000000000004', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a24', 'Branco Bariloche', '#F5F5F5'),
  ('00000000-0024-0000-0002-000000000005', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a24', 'Verde Oásis', '#3E5F4A'),
  ('00000000-0024-0000-0002-000000000006', '2a2a2a2a-2a2a-2a2a-2a2a-2a2a2a2a2a24', 'Preto Toronto', '#0A0A0A')
ON CONFLICT (id) DO NOTHING;
