#!/bin/sh
# torque backend local development setup
#
# Usage:
#   ./scripts/setup.sh
#
# Bootstraps the full local environment in order:
#   1. Start infrastructure
#   2. Run database migrations
#   3. Wait for step-ca to become healthy
#   4. Configure step-ca certificate duration limits
#   5. Copy root CA certificate
#   6. Create demo Kratos identity
#   7. Issue test device certificate and seed in Postgres
#   8. Seed reference data
#
# Idempotent — safe to run multiple times.
# Requires: docker, docker compose
set -e

COMPOSE="docker compose -f $(cd "$(dirname "$0")/../infra" && pwd)/docker-compose.yml"
REPO_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CERTS_DIR="$REPO_DIR/certs"
MIGRATIONS_DIR="$REPO_DIR/migrations"

mkdir -p "$CERTS_DIR/device"

step_exec() { $COMPOSE exec -T step-ca "$@"; }
psql_exec() { $COMPOSE exec -T postgres psql -U torque -d torque -c "$1"; }
psql_exec_file() { $COMPOSE exec -T postgres psql -U torque -d torque -f - < "$1"; }
kratos_exec() { $COMPOSE exec -T kratos "$@"; }

issue_cert() {
  CN="$1"; OUT_CRT="$2"; OUT_KEY="$3"
  step_exec sh -c "
    step ca certificate '$CN' /tmp/out.crt /tmp/out.key \
      --provisioner=torque \
      --provisioner-password-file=/run/secrets/ca-password \
      --not-after=8760h -f > /dev/null 2>&1
  "
  $COMPOSE cp step-ca:/tmp/out.crt "$OUT_CRT"
  $COMPOSE cp step-ca:/tmp/out.key "$OUT_KEY"
}

create_demo_identity() {
  RESP=$(kratos_exec wget -qO- \
    --post-data '{
      "schema_id": "default",
      "traits": {
        "email": "demo@torque.dev",
        "name": { "first": "Demo", "last": "User" },
        "role": "admin"
      }
    }' \
    --header "Content-Type: application/json" \
    http://localhost:4434/admin/identities 2>&1) || true
  if echo "$RESP" | grep -q '"id"'; then
    echo "  ✓ demo@torque.dev (admin)"
  elif echo "$RESP" | grep -q "409"; then
    echo "  ⏭ demo@torque.dev (admin) — already exists"
  else
    echo "  ✗ demo@torque.dev FAILED: $RESP"
  fi
}

print_step() {
  echo ""
  echo "──────────────────────────────────────────────────"
  printf " %s\n" "$1"
  echo "──────────────────────────────────────────────────"
}

# ──────────────────────────────────────────────────────────────
print_step "1/8  Starting infrastructure"
# ──────────────────────────────────────────────────────────────
$COMPOSE up -d --remove-orphans
echo "  ✓ services started"

# ──────────────────────────────────────────────────────────────
print_step "2/8  Running database migrations"
# ──────────────────────────────────────────────────────────────
until $COMPOSE exec -T postgres pg_isready -U torque > /dev/null 2>&1; do
  printf "  waiting for postgres...\r"
  sleep 2
done

MIGRATE="migrate -path $MIGRATIONS_DIR/main -database postgres://torque:torque@localhost:5432/torque?sslmode=disable"

$MIGRATE up 2>&1 | while IFS= read -r line; do echo "  $line"; done
echo "  ✓ main database migrations complete"

# Timescale migrations (separate database)
until $COMPOSE exec -T timescaledb pg_isready -U torque > /dev/null 2>&1; do
  printf "  waiting for timescaledb...\r"
  sleep 2
done

MIGRATE_TS="migrate -path $MIGRATIONS_DIR/timescale -database postgres://torque:torque@localhost:5433/torque?sslmode=disable"

$MIGRATE_TS up 2>&1 | while IFS= read -r line; do echo "  $line"; done
echo "  ✓ timescaledb migrations complete"

# ──────────────────────────────────────────────────────────────
print_step "3/8  Waiting for step-ca to become healthy"
# ──────────────────────────────────────────────────────────────
until step_exec step ca health \
  --ca-url=https://step-ca:9000 \
  --root=/home/step/certs/root_ca.crt > /dev/null 2>&1; do
  printf "  waiting...\r"
  sleep 2
done
echo "  ✓ step-ca is healthy"

# ──────────────────────────────────────────────────────────────
print_step "4/8  Configuring step-ca certificate duration limits"
# ──────────────────────────────────────────────────────────────
step_exec sh -c "
  jq '.authority.claims = {
    \"minTLSCertDuration\": \"5m\",
    \"maxTLSCertDuration\": \"8760h\",
    \"defaultTLSCertDuration\": \"8760h\"
  }' /home/step/config/ca.json > /tmp/ca.json.tmp &&
  mv /tmp/ca.json.tmp /home/step/config/ca.json
"
$COMPOSE restart step-ca
until step_exec step ca health \
  --ca-url=https://step-ca:9000 \
  --root=/home/step/certs/root_ca.crt > /dev/null 2>&1; do
  printf "  waiting...\r"
  sleep 2
done
echo "  ✓ certificate duration configured (max: 8760h)"

# ──────────────────────────────────────────────────────────────
print_step "5/8  Copying root CA certificate"
# ──────────────────────────────────────────────────────────────
$COMPOSE cp step-ca:/home/step/certs/root_ca.crt "$CERTS_DIR/ca.crt"
echo "  ✓ root cert written to $CERTS_DIR/ca.crt"

# ──────────────────────────────────────────────────────────────
print_step "6/8  Creating demo Kratos identity"
# ──────────────────────────────────────────────────────────────
until kratos_exec wget -qO- http://localhost:4434/health/ready > /dev/null 2>&1; do
  printf "  waiting for kratos...\r"
  sleep 2
done

create_demo_identity

# ──────────────────────────────────────────────────────────────
print_step "7/8  Issuing test device certificate and seeding Postgres"
# ──────────────────────────────────────────────────────────────
META="$CERTS_DIR/device/meta.json"

# Reset state: if meta exists but DB is fresh, re-seed
if [ -f "$META" ]; then
  DEVICE_ID=$(grep -o '"device_id": *"[^"]*"' "$META" | cut -d'"' -f4)
  EXISTS=$(psql_exec "SELECT count(*) FROM device.devices WHERE id = '$DEVICE_ID';" 2>/dev/null || echo "0")
  if echo "$EXISTS" | grep -q "^[[:space:]]*0"; then
    echo "  ⚠ meta.json exists but device not found in DB, re-seeding"
    rm -f "$META"
  fi
fi

if [ -f "$META" ]; then
  DEVICE_ID=$(grep -o '"device_id": *"[^"]*"' "$META" | cut -d'"' -f4)
  echo "  ⏭ device TRQ-1 already exists (ID: $DEVICE_ID), skipping"
else
  DEVICE_ID=$(cat /proc/sys/kernel/random/uuid 2>/dev/null || uuidgen | tr '[:upper:]' '[:lower:]')

  issue_cert "$DEVICE_ID" "$CERTS_DIR/device/device.crt" "$CERTS_DIR/device/device.key"
  cp "$CERTS_DIR/ca.crt" "$CERTS_DIR/device/ca.crt"

  CERT_SN=$(openssl x509 -in "$CERTS_DIR/device/device.crt" -noout -serial 2>/dev/null \
    | cut -d= -f2 | tr '[:upper:]' '[:lower:]' | sed 's/../&:/g;s/:$//')

  psql_exec "
    INSERT INTO device.devices (id, name, certificate_cn, certificate_sn)
    VALUES ('$DEVICE_ID', 'TRQ-1', '$DEVICE_ID', '$CERT_SN')
    ON CONFLICT (certificate_cn) DO NOTHING;
  "

  printf '{\n  "device_id": "%s",\n  "certificate_sn": "%s"\n}\n' \
    "$DEVICE_ID" "$CERT_SN" > "$META"

  echo "  ✓ test device seeded (name: TRQ-1, CN: $DEVICE_ID)"
fi

# ──────────────────────────────────────────────────────────────
print_step "8/8  Seeding reference data"
# ──────────────────────────────────────────────────────────────
SEEDS_DIR="$REPO_DIR/seeds"

for f in "$SEEDS_DIR"/*.sql; do
  [ -f "$f" ] || continue
  echo "  seeding $(basename "$f")"
  psql_exec_file "$f"
done
echo "  ✓ reference data seeded"

echo ""
echo "══════════════════════════════════════════════════"
echo " Setup complete!"
echo "══════════════════════════════════════════════════"
echo ""
echo "  Identity:"
echo "    demo@torque.dev (role: admin)"
echo ""
echo "  Device: TRQ-1 (ID: $DEVICE_ID)"
echo ""
echo "  Certs:"
echo "    CA     → $CERTS_DIR/ca.crt"
echo "    Device → $CERTS_DIR/device/"
echo ""
echo "  Services:"
echo "    minio    → http://localhost:9000 (API) / http://localhost:9001 (Console)"
echo "    mailhog  → http://localhost:8025"
echo ""
