# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run a service
go run ./cmd/api
go run ./cmd/mqtt-guard

# Build binaries
go build ./cmd/api
go build ./cmd/mqtt-guard

# Run all tests
go test ./...

# Run a specific test
go test -run TestName ./internal/modules/vehicle/...

# Vet
go vet ./...

# Start infrastructure (Postgres, Vault, EMQX, Traefik, Kratos, etc.)
docker compose -f infra/docker-compose.yml up -d
```

## Architecture

Two binaries share a single Go module:

- **`cmd/api`** — REST API server (go-chi). Authenticated via headers `x-user-id` / `x-user-role` injected by Ory Oathkeeper (reverse proxy in front).
- **`cmd/mqtt-guard`** — EMQX webhook service for MQTT client auth and ACL. Requests validated via `X-Guard-Secret` header.

### Clean Architecture layers (per module)

```
internal/modules/<module>/
  domain/         → entities, value objects, repository interfaces
  application/    → use cases, DTOs
  infrastructure/ → GORM repository implementations
```

Modules: `vehicle`, `device`, `telemetry` (stub).

### Core packages (`internal/core/`)

| Package | Purpose |
|---|---|
| `apperr` | App error type with `Kind` (BadRequest, NotFound, Conflict, etc.) mapped to HTTP status codes |
| `appctx` | Carries `AuthContext` (UserID + Role) through `context.Context` |
| `db` | GORM Postgres connection + auto-migration; `AuditableModel` base struct |
| `pki` | HashiCorp Vault client for certificate issuance / revocation |
| `pagination` | `Page` struct with safe defaults (1–100 items) |
| `logger` | Uber Zap (JSON or human-readable via `LOG_FORMAT` env) |

### Infrastructure stack

| Component | Role |
|---|---|
| Traefik | Reverse proxy, TLS termination |
| Ory Oathkeeper | Enforces auth policies; injects identity headers into upstream |
| Ory Kratos | Passwordless (email code) identity management |
| PostgreSQL 17 | Primary datastore |
| TimescaleDB | Time-series telemetry data |
| HashiCorp Vault | PKI: issues/revokes device TLS certificates |
| EMQX 5.8 | MQTT broker (mTLS on 8883); calls mqtt-guard webhooks for auth/ACL |
| RabbitMQ | Message broker (telemetry pipeline) |
| MinIO | S3-compatible object storage (3D model assets) |

### Error handling pattern

Use `apperr.New(kind, message)` in domain/use-case code. HTTP handlers in `cmd/*/httperr/` translate these to appropriate HTTP status responses via the `Kind` → status mapping.

### Device lifecycle

1. `CreateDevice` — provisions a device record and issues a Vault PKI certificate (private key stored only at creation).
2. `CommissionDevice` — links device to a vehicle; authorizes MQTT topic `device/<id>/telemetry`.
3. `DecommissionDevice` — unlinks device; Vault certificate revoked.
4. EMQX calls `mqtt-guard /mqtt/auth` (client cert CN = device ID) and `/mqtt/acl` per publish/subscribe.
