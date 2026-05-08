# SPEC Reference (untuk Claude Code & developer)

> Portable & committable. Bisa di-share ke tim.

## Purpose

Saat ada permintaan **tanya / tambah / ubah fitur** di repo `mediaplatform-datasource`, **baca dokumen SDD dulu** sebelum implementasi:

| Doc | Kapan dipakai |
|-----|---------------|
| [`SPEC.md`](../SPEC.md) | Source-of-truth — 8 section (business context, actors, use case, API, data model, constraints, edge cases, AC) + workflow |
| [`SPEC.html`](../SPEC.html) | Versi HTML, buka di browser untuk review yg lebih readable |
| [`AUDIT_FIX.md`](../AUDIT_FIX.md) | Log 8 audit fix items (TLS, AES, DB pool, Redis, config, JWT, migrate, router/rmq) — context kenapa code yg ada bentuknya begini |
| [`UNUSED_FUNCTIONS.md`](../UNUSED_FUNCTIONS.md) | Dead code report — reference sblm hapus/refactor fungsi |
| [`TESTING.md`](../TESTING.md) | Test strategy + roadmap 5 phase |

## Workflow saat tambah fitur

1. **Identify domain**: campaign / budget / master-data / report / postback / dll
2. **Update SPEC.md (+ SPEC.html)** dulu sebelum koding:
   - Section 3: tambah Use Case (UC-XX)
   - Section 4.4: tambah endpoint definition
   - Section 5: tambah entity (kalau perlu DB schema baru)
   - Section 7: tambah edge cases
   - Section 8.1: tambah Acceptance Criteria
3. **Implement** sesuai DDD layer:
   - `src/domain/entity/<file>.go` (data struct + register di `migrateEntities`)
   - `src/domain/repository/<file>.go` (DB access)
   - `src/handler/incoming_<file>_handler.go` (HTTP handler)
   - `src/app/routes/<group>.go` (route registration)
4. **Test**: `make test`, `make staticcheck`, `make build`
5. **Migrate**: `./datasource migrate` (kalau ada entity baru)
6. **Deploy**: tag git → CI auto build → run migrate compose → deploy server compose

## Definition of Done (lihat SPEC.md §8.2)

Sebelum merge, pastikan:
- Endpoint terdaftar di `src/app/routes/<domain>.go`
- Handler di `src/handler/incoming_<domain>_handler.go`
- Repository (kalau ada DB) di `src/domain/repository/<domain>.go`
- Entity (kalau ada schema baru) registered di `migrateEntities` (`src/cmd/migrate.go`)
- JWT auth applied (kalau bukan public — public hanya postback callback)
- Scope filter (company/agency/adnet) sesuai actor
- Error response konsisten: `{code, desc, error}` envelope
- Response time: <1s read, <3s export
- Logs structured (logrus) — no `fmt.Println` bocor
- Secrets tidak di-log (pakai `Cfg.Redacted()` pattern)
- `staticcheck ./...` zero new warning
- `go build ./...` exit 0
- SPEC.md + SPEC.html updated

## Architecture (DDD)

```
cores/datasource/
├── main.go                  # entry point (root)
├── go.mod
├── Makefile                 # build/test/docker targets
├── Dockerfile.datasource.{server,migrate}
├── .github/workflows/       # CI matrix build (server + migrate image)
├── SPEC.md / SPEC.html      # source-of-truth
├── AUDIT_FIX.md
├── UNUSED_FUNCTIONS.md
├── TESTING.md
├── .claude/
│   └── SPEC_REFERENCE.md    # this file
└── src/
    ├── app/
    │   ├── url_mapping.go   # orchestrator (auth gate, group setup)
    │   └── routes/          # per-domain route registration
    ├── cmd/                 # cobra sub-commands (server, migrate)
    ├── config/              # Cfg struct, env parsing, init Redis/DB/Rmq/GS
    ├── handler/             # HTTP handlers (188+ methods)
    ├── helper/              # http, utils (crypto), redis_safe, logger
    └── domain/
        ├── entity/          # data structs / GORM models
        └── repository/      # DB access layer (was /model)
```

## Stack

- **Lang**: Go 1.24
- **Framework**: Fiber v2
- **DB**: PostgreSQL via GORM
- **Cache**: Redis (go-redis + rueidis untuk JSON)
- **Queue**: RabbitMQ via `wiliehidayat87/rmqp`
- **Sheets**: Google Sheets API
- **CLI**: Cobra
- **JWT**: `golang-jwt/jwt v4` HS256
- **Test**: testify + testing stdlib

## Constraints critical (must-not-violate)

| Constraint | Detail | Source |
|------------|--------|--------|
| TLS verify default secure | Skip verify hanya non-prod opt-in | `helper/http.go` (#1 fix) |
| AES key dari env | `AES_SECRET_KEY` (16/24/32 byte), bukan hardcoded | `helper/utils.go` (#1 fix) |
| DB pool bounded | Default `max_open=100`, `max_idle=10` | `config/c.go` (#3 fix) |
| Redis fail-fast atau degraded | Gate `REDIS_REQUIRED` env | `config/c.go` (#4 fix) |
| Server tidak run AutoMigrate | Pakai `./datasource migrate` sub-cmd | `cmd/migrate.go` (#7 fix) |
| Logs redact secrets | `Cfg.Redacted()` mask password/key | `config/c.go` (#5 fix) |
| Group-level auth default | `AUTH_ENFORCE_DEFAULT` env gate | `app/url_mapping.go` (#6 fix) |
| JWT validation full | aud/iss/nbf/jti blacklist | `handler/incoming_auth_handler.go` (#6 fix) |

## Quick reference: Adding endpoint

```go
// 1. src/domain/entity/feature.go (kalau perlu schema)
type Feature struct {
    ID   uint   `gorm:"primaryKey"`
    Name string
}

// 2. src/cmd/migrate.go — tambah ke migrateEntities slice
&entity.Feature{},

// 3. src/domain/repository/feature.go
func (b *BaseModel) GetFeatures() ([]entity.Feature, error) {
    var f []entity.Feature
    return f, b.DB.Find(&f).Error
}

// 4. src/handler/incoming_feature_handler.go
func (h *IncomingHandler) DisplayFeatures(c *fiber.Ctx) error {
    items, err := h.DS.GetFeatures()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"code": 500, "desc": "internal", "error": err.Error()})
    }
    return c.JSON(fiber.Map{"code": 200, "desc": "OK", "data": items})
}

// 5. src/app/routes/<group>.go — register
grp.Get("/features", h.DisplayFeatures).Name("List features")

// 6. SPEC.md update section 3, 4.4, 5, 7, 8.1
// 7. make test && make build && ./datasource migrate
```
