#!/bin/sh
# torque backend local development setup
#
# Usage:
#   ./scripts/setup.sh
#
# Bootstraps the full local environment in order:
#   1. Start infrastructure
#   2. Wait for step-ca to become healthy
#   3. Configure step-ca certificate duration limits
#   4. Copy root CA certificate
#   5. Create test Kratos identities (admin, support, mechanical)
#   6. Issue test device certificate and seed in Postgres
#
# Idempotent — safe to run multiple times.
# Requires: docker, docker compose
set -e

COMPOSE="docker compose -f $(cd "$(dirname "$0")/../infra" && pwd)/docker-compose.yml"
REPO_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CERTS_DIR="$REPO_DIR/certs"

mkdir -p "$CERTS_DIR/device"

step_exec() { $COMPOSE exec -T step-ca "$@"; }
psql_exec() { $COMPOSE exec -T postgres psql -U torque -d torque -c "$1"; }
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

create_identity() {
  EMAIL="$1"; FIRST="$2"; LAST="$3"; ROLE="$4"
  kratos_exec wget -qO- \
    --post-data "{
      \"schema_id\": \"default\",
      \"traits\": {
        \"email\": \"$EMAIL\",
        \"name\": { \"first\": \"$FIRST\", \"last\": \"$LAST\" },
        \"role\": \"$ROLE\"
      }
    }" \
    --header "Content-Type: application/json" \
    http://localhost:4434/admin/identities > /dev/null 2>&1 || true
}

print_step() {
  echo ""
  echo "──────────────────────────────────────────────────"
  printf " %s\n" "$1"
  echo "──────────────────────────────────────────────────"
}

# ──────────────────────────────────────────────────────────────
print_step "1/6  Starting infrastructure"
# ──────────────────────────────────────────────────────────────
$COMPOSE up -d
echo "  ✓ services started"

# ──────────────────────────────────────────────────────────────
print_step "2/6  Waiting for step-ca to become healthy"
# ──────────────────────────────────────────────────────────────
until step_exec step ca health \
  --ca-url=https://localhost:9000 \
  --root=/home/step/certs/root_ca.crt > /dev/null 2>&1; do
  printf "  waiting...\r"
  sleep 2
done
echo "  ✓ step-ca is healthy"

# ──────────────────────────────────────────────────────────────
print_step "3/6  Configuring step-ca certificate duration limits"
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
  --ca-url=https://localhost:9000 \
  --root=/home/step/certs/root_ca.crt > /dev/null 2>&1; do
  printf "  waiting...\r"
  sleep 2
done
echo "  ✓ certificate duration configured (max: 8760h)"

# ──────────────────────────────────────────────────────────────
print_step "4/6  Copying root CA certificate"
# ──────────────────────────────────────────────────────────────
$COMPOSE cp step-ca:/home/step/certs/root_ca.crt "$CERTS_DIR/ca.crt"
echo "  ✓ root cert written to $CERTS_DIR/ca.crt"

# ──────────────────────────────────────────────────────────────
print_step "5/6  Creating test Kratos identities"
# ──────────────────────────────────────────────────────────────
until kratos_exec wget -qO- http://localhost:4434/health/ready > /dev/null 2>&1; do
  printf "  waiting for kratos...\r"
  sleep 2
done

create_identity "admin@torque.dev"      "Admin"      "User" "admin"
create_identity "support@torque.dev"    "Support"    "User" "support"
create_identity "mechanical@torque.dev" "Mechanical" "User" "mechanical"

echo "  ✓ identities created:"
echo "      admin@torque.dev      (role: admin)"
echo "      support@torque.dev    (role: support)"
echo "      mechanical@torque.dev (role: mechanical)"

# ──────────────────────────────────────────────────────────────
print_step "6/6  Issuing test device certificate and seeding Postgres"
# ──────────────────────────────────────────────────────────────
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
  "$DEVICE_ID" "$CERT_SN" > "$CERTS_DIR/device/meta.json"

echo "  ✓ test device seeded (name: TRQ-1, CN: $DEVICE_ID)"

echo ""
echo "══════════════════════════════════════════════════"
echo " Setup complete!"
echo "══════════════════════════════════════════════════"
echo ""
echo "  Identities:"
echo "    admin@torque.dev      (role: admin)"
echo "    support@torque.dev    (role: support)"
echo "    mechanical@torque.dev (role: mechanical)"
echo ""
echo "  Device: TRQ-1 (ID: $DEVICE_ID)"
echo ""
echo "  Certs:"
echo "    CA     → $CERTS_DIR/ca.crt"
echo "    Device → $CERTS_DIR/device/"
echo ""
echo "  Services:"
echo "    step-ca  → https://localhost:9000"
echo "    postgres → localhost:5432"
echo "    api      → http://localhost:80"
echo ""
