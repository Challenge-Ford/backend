#!/bin/sh
set -e

VAULT_ADDR="http://vault:8200"
VAULT_TOKEN="torque"

echo "Waiting for Vault..."
until vault status > /dev/null 2>&1; do
  sleep 2
done
echo "Vault is ready"

export VAULT_ADDR VAULT_TOKEN

echo "Enabling PKI engine..."
vault secrets enable pki || echo "PKI already enabled, skipping"

echo "Configuring PKI max TTL..."
vault secrets tune -max-lease-ttl=87600h pki

echo "Generating root CA..."
vault write -format=json pki/root/generate/internal \
  common_name="Torque Device CA" \
  ttl=87600h > /dev/null && echo "Root CA generated" || echo "Root CA may already exist, skipping"

echo "Configuring PKI URLs..."
vault write pki/config/urls \
  issuing_certificates="$VAULT_ADDR/v1/pki/ca" \
  crl_distribution_points="$VAULT_ADDR/v1/pki/crl"

echo "Creating device role..."
vault write pki/roles/device \
  allow_any_name=true \
  enforce_hostnames=false \
  key_type=ec \
  key_bits=256 \
  ttl=8760h \
  max_ttl=87600h \
  no_store=false

echo "Vault setup complete"
