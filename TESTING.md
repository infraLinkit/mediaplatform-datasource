# Testing Guide

## Status saat ini

| Package | Test file | Coverage |
|---------|-----------|----------|
| `src/helper` | `utils_test.go`, `redis_safe_test.go`, `http_test.go` | **36.4%** |
| `src/config` | `c_test.go` | **35.0%** |
| `src/cmd` | `migrate_test.go` | 4.4% (smoke only) |
| `src/app` | `url_mapping_test.go` | **90.9%** |
| `src/handler` | `incoming_auth_handler_test.go` | 0.2% (sample only) |
| `src/app/routes` | – | 0% |
| `src/domain/entity` | – | 0% (data structs only) |
| `src/domain/repository` | – | 0% (need DB mock) |

Run: `make test`. Coverage: `make cover` / `make cover-html`.

---

## Coverage rationale

### Sudah di-test (critical paths)

1. **`helper/utils.go`** — Encrypt/Decrypt round-trip, AES key validation, nonce randomness, tampering detection. Fix #1 (CWE-321) verified.
2. **`helper/redis_safe.go`** — Nil-safety semua wrapper. Fix #4 degraded mode verified.
3. **`helper/http.go`** — TLS config (default secure, opt-in non-prod), Get/Post sukses & error path, **nil response no panic** (fix #2). Block try/catch.
4. **`config/c.go`** — `envIntDefault`, `Cfg.Redacted()` (semua secret field), `InitCfg` env defaults & overrides. Fix #3 + #5 verified.
5. **`handler/incoming_auth_handler.go`** — `audienceMatches` (string + array per RFC 7519), `RevokeJWT` degraded mode. Fix #6 helper verified.
6. **`cmd/migrate.go`** — `migrateEntities` non-empty, all pointers, no duplicate, no nil. Sub-cmd registration.
7. **`app/url_mapping.go`** — `MapUrls` no panic, `authEnforceDefault` env parsing, route registration smoke (404 vs 500 distinction).

### Belum di-test (TODO terstruktur)

| Package | Yg masih kosong | Kenapa skip dulu |
|---------|----------------|------------------|
| `domain/repository/*.go` | 28 file repository | Butuh `sqlmock` per repository untuk test query. Setup ~30-50 lines per file. |
| `handler/incoming_*.go` | 188 handler method | Butuh fiber.Ctx mock + DB mock + Redis mock + Rmqp mock. Setup test infra dulu sebelum bulk write. |
| `app/routes/*.go` | Cuma map handler ke route — covered via `url_mapping_test.go` smoke |
| `domain/entity/*.go` | Plain struct, no logic — coverage tidak meaningful |

---

## Cara extend test

### Tambah test untuk repository (DB layer)

Setup `sqlmock`:
```go
import (
    "database/sql"
    "regexp"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/require"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err)
    db, err := gorm.Open(postgres.New(postgres.Config{
        Conn:                 sqlDB,
        PreferSimpleProtocol: true,
    }), &gorm.Config{})
    require.NoError(t, err)
    return db, mock
}

func TestSomeRepoQuery(t *testing.T) {
    db, mock := newMockDB(t)
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "campaigns"`)).
        WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
            AddRow(1, "Test"))
    // ... call repo method, assert
    require.NoError(t, mock.ExpectationsWereMet())
}
```

Add dep: `go get github.com/DATA-DOG/go-sqlmock`

### Tambah test untuk handler

Setup Fiber test:
```go
import (
    "net/http/httptest"
    "github.com/gofiber/fiber/v2"
)

func TestSomeHandler(t *testing.T) {
    app := fiber.New()
    h := &handler.IncomingHandler{
        DB:   mockDB,
        RCP:  nil, // degraded mode test
        Logs: logrus.New(),
    }
    app.Get("/test", h.SomeHandler)

    req := httptest.NewRequest("GET", "/test", nil)
    resp, err := app.Test(req, -1)
    require.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

### Mock untuk Redis (miniredis)

```go
import "github.com/alicebob/miniredis/v2"

func newMockRedis(t *testing.T) *redis.Client {
    s := miniredis.RunT(t)
    return redis.NewClient(&redis.Options{Addr: s.Addr()})
}
```

Add dep: `go get github.com/alicebob/miniredis/v2`

### Mock untuk RabbitMQ

`rmqp.AMQP` lib tidak punya interface — recommendation: introduce wrapper interface di `helper/`:
```go
type Publisher interface {
    Publish(exchange, queue string, payload []byte) error
}
```
Lalu refactor `c.Rmqp.IntegratePublish` jadi method via interface. Mock-able di test.

---

## Test conventions

### Struktur file

- File test: `<source>_test.go` di package yg sama
- Naming: `Test<FunctionName>_<Scenario>` (e.g., `TestEncrypt_KeyMissing`)
- Sub-test untuk variasi: `t.Run("name", func(t *testing.T){ ... })`
- Table-driven test untuk multi-case

### Tools

- `testify/assert` — soft assertion
- `testify/require` — hard assertion (stop on fail). Pakai untuk setup yg invalid → tidak guna lanjut test.
- `httptest` — HTTP server mock (stdlib)
- `t.Setenv` — set env per test (auto cleanup)

### Dilarang di test

- Hit DB beneran (gunakan sqlmock)
- Hit Redis beneran (gunakan miniredis)
- Hit external API (gunakan httptest server)
- Network call ke production
- Hardcoded sleep > 100ms (race + flaky)

### CI integration

Makefile target:
- `make ci` — `deps + check + test + build` (no fmt/tidy modifications)
- `make test` — semua test dengan `-race -count=1`
- `make cover` — coverage report

GitHub workflow `main.yaml` saat ini build Docker image. Future improvement: tambah test job sebelum build.

---

## Roadmap

### Phase 1 (done) — foundation
- [x] Critical helper coverage
- [x] Config redact + env parsing
- [x] Auth helper functions
- [x] Smoke test url_mapping
- [x] Migrate entities sanity

### Phase 2 — repository layer
- [ ] Setup sqlmock helper di `src/domain/repository/testutil/`
- [ ] Test 5 hot-path repos: `campaign_management`, `dashboard`, `postback`, `target_budget`, `summary_landing`

### Phase 3 — handler layer
- [ ] Mock infra untuk fiber.Ctx + DB + Redis (`src/handler/testutil/`)
- [ ] Test 10 critical handler: postback, auth flow, edit campaign, dashboard data, budget IO

### Phase 4 — integration
- [ ] Docker compose untuk test PG + Redis + RMQ
- [ ] E2E test untuk postback flow (HTTP → DB → RMQ message)
- [ ] Load test postback endpoint (≥ 500 req/sec)

### Phase 5 — CI
- [ ] Add test job ke `.github/workflows/main.yaml` (run sebelum docker build)
- [ ] Coverage threshold gate (e.g., min 30% target awal, naikkan progressive)
- [ ] Codecov / coveralls upload
