# Audit Fix Implementation Log

Tindak lanjut Code Audit Report `Data Source.pdf` (analysis date 2026-02-11).
Tanggal mulai implementasi: 2026-05-08.

## Ringkasan

| # | Task | Severity | File Utama | Status |
|---|------|----------|------------|--------|
| 1 | TLS verify + AES key externalization | High | `helper/http.go`, `helper/utils.go` | ✅ |
| 2 | Nil-deref panic di HTTP helper | Major | `helper/http.go` | ✅ |
| 3 | DB pool unbounded → bounded | Major | `config/c.go` | ✅ |
| 4 | Redis panic → retry + degraded mode | Major | `config/c.go`, `helper/redis_safe.go` | ✅ |
| 5 | Hardcoded APP_PATH + config redact | Minor | `config/c.go` | ✅ |
| 6 | JWT validation (aud/iss/nbf/jti) + auth default | High | `handler/incoming_auth_handler.go`, `app/url_mapping.go` | ✅ |
| 7 | AutoMigrate gating + `migrate` sub-cmd | Minor | `cmd/migrate.go`, `cmd/server.go` | ✅ |
| 8 | Router split + RabbitMQ reconnect | Smell | `app/routes/*`, `cmd/server.go` | ✅ |

---

## #1. TLS Verify + AES Key

### Issue
- `helper/http.go:88` — `InsecureSkipVerify: true` hardcoded → MITM risk (CWE-295)
- `helper/utils.go:183` — AES key `"N1PCdw3M2B1TfJho"` hardcoded → CWE-321

### Fix
- `PHttp` struct ditambah field `InsecureSkipVerify bool` (default false = secure).
- TLS skip verify hanya aktif kalau opt-in DAN `APP_ENV != production` (double guard).
- `MinVersion: tls.VersionTLS12` ditambahkan.
- AES key dipindah ke env `AES_SECRET_KEY` (16/24/32 byte → AES-128/192/256).
- `Encrypt`/`Decrypt` ubah signature jadi `(string, error)` — drop semua `panic`.
- Output ciphertext base64-encoded (safe transport).

### Env baru
```
AES_SECRET_KEY=<32-byte string>
APP_ENV=production
```

⚠️ **Action sebelum prod**: rotate `AES_SECRET_KEY`, simpan di secret manager (Vault, K8s Secret, AWS Secrets Manager).

---

## #2. Nil Response Dereference

### Issue
- `helper/http.go:163, 256` — saat `httpClient.Do(req)` error, code akses `response.Status`/`StatusCode` padahal `response` nil → service crash.
- Bonus: `req` juga bisa nil kalau `NewRequest` error tapi diakses sebelum check.

### Fix
- Cek `err` dari `NewRequest` SEBELUM akses `req.Header`.
- Drop `response.Status/StatusCode` di error branch.
- Tambah `response == nil` guard.
- Non-OK status return real status code + wrapped error (sebelumnya `nil` err).

### Behavior change
Caller yg sebelumnya assume `err==nil` saat status≠200 sekarang akan dapat error. Cek caller di `model/arpu_linkitdashboard.go` jika perlu adjust.

---

## #3. DB Connection Pool

### Issue
- `config/c.go:317-323` — `SetMaxIdleConns(0)`, `SetMaxOpenConns(0)`, `SetConnMaxLifetime(0)` = unbounded → file-descriptor exhaustion + DB overload.

### Fix
- Tambah field `DBMaxIdleConns`, `DBMaxOpenConns`, `DBConnMaxLifetime`, `DBConnMaxIdleTime` ke `Cfg`.
- Helper `envIntDefault(key, default)`.
- Apply via `sqlDB.SetMax*Conns(...)` + `SetConnMaxIdleTime`.

### Env baru
```
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME_MIN=30
DB_CONN_MAX_IDLE_TIME_MIN=10
```

⚠️ **Tuning**: `MaxOpenConns × replicas ≤ PG max_connections` (default 100). Multi-instance: bagi proporsional.

---

## #4. Redis Panic + Degraded Mode

### Issue
- `config/c.go:249, 274` — `panic(errRedis)` saat connection fail → no retry, no degraded path.

### Fix
- Helper `retryWithBackoff`: 5 attempts, exp backoff 1s→16s.
- `InitRedis` & `InitRedisJSON` return `(*Client, error)` instead of panic.
- `rueidis.New` panic recovered di goroutine wrapper (lib internal panic).
- `Initiate` return `(*Setup, error)`.
- Caller `cmd/server.go` handle error dgn `log.Fatalf`.

### Degraded mode (bonus)
- Env baru `REDIS_REQUIRED` (default `true`).
- Kalau `false` + Redis fail → log warn, app continue dgn `nil` client + `Setup.RedisAvailable=false`.
- File baru `helper/redis_safe.go` — nil-safe wrappers:
  - `SafeRedisGet/Set/Del`
  - `SafeRueidisJSONGet/Exists`
- Hot path bisa migrate pake wrapper supaya graceful kalau Redis down.

### Env baru
```
REDIS_REQUIRED=true
```

---

## #5. APP_PATH + Config Redact

### Issue
- `config/c.go:25` — `APP_PATH` default `/Users/wiliewahyuhidayat/Documents/GO/...` (developer path).
- `config/c.go:217` — `l.Info(fmt.Sprintf("Config Loaded : %#v\n", c))` dump full struct termasuk password/private key.

### Fix
- Default `APP_PATH = "./"`. Warning log kalau env `APPPATH` tidak diset.
- Method `Cfg.Redacted()` — return copy dgn fields sensitif jadi `"***REDACTED***"`:
  - `PSQLPassword`, `RedisPwd`, `RabbitMQPassword`
  - `ARPUUsername`, `ARPUPassword`
  - `GSPrivateKey`, `GSPrivateKeyID`, `GSClientID`
- `Initiate` log pakai `c.Redacted()`.

---

## #6. JWT Validation + Auth Middleware

### Issue
- `handler/incoming_auth_handler.go` — JWT cuma validate `exp`, `type`, `jti` keberadaan. Missing `aud`, `iss`, `nbf`, jti revocation.
- `app/url_mapping.go` — banyak admin endpoint (`/v1/management/*`, `/v1/int/*`, sebagian `/v1/report/*`) tanpa `AuthMiddleware`.

### Fix JWT
- `nbf` check (30s leeway).
- `aud` check — opt-in via `JWT_AUDIENCE` env, support string + array (RFC 7519).
- `iss` check — opt-in via `JWT_ISSUER` env.
- `jti` blacklist via Redis: key `jwt:blacklist:<jti>`. Skip kalau Redis nil (degraded mode).
- Helper `audienceMatches(claim, expected)`.
- Method `RevokeJWT(jti, ttl)` — pakai saat logout.

### Fix Route Exposure
- Group-level middleware ditambah ke `/dashboard`, `/v1/report`, `/v1/int`, `/v1/management`.
- Gate via `AUTH_ENFORCE_DEFAULT` env (default `false` = back-compat, `true` = enforce).
- Postback (`/v1/postback*`, `/v1/inquire/campid`) explicit public dengan comment alasan.

### Env baru
```
JWT_AUDIENCE=
JWT_ISSUER=
AUTH_ENFORCE_DEFAULT=false
```

⚠️ **Migration**: pastikan FE kirim `Authorization: Bearer <token>` ke semua admin/management endpoint sebelum flip `AUTH_ENFORCE_DEFAULT=true`.

---

## #7. AutoMigrate Gating + Migrate Sub-Command

### Issue
- `cmd/server.go:22` — AutoMigrate one-line monster, jalan setiap startup, no error handling, no timeout.

### Fix
- Slice `migrateEntities` extract dari one-liner (readable + reusable).
- **AutoMigrate dihapus total dari `cmd/server.go`** — server tidak lagi handle schema migration.
- Sub-command `cmd/migrate.go`: `./datasource migrate` — standalone migrasi, exit setelah selesai. Context timeout via `AUTO_MIGRATE_TIMEOUT_MIN` (default 5m).
- `migrateEntities` slice dipindah dari `server.go` ke `migrate.go` (cohesion: dipakai cuma oleh migrate cmd).

### Bonus: Docker / CI
- `Dockerfile.datasource.migrate` — clone server image dgn `CMD ["/datasource", "migrate"]`.
- `.github/workflows/main.yaml` — matrix include 2 entries (server + migrate). Tag push → both image dibuild & push paralel ke Docker Hub.
  - `infralinkit/mediaplatform-datasource-server:<tag>`
  - `infralinkit/mediaplatform-datasource-migrate:<tag>`
- `resources/dockerver/yaml/production/docker/docker-compose.migrate.yaml` — run-once compose, no ports, `restart_policy: none`.

### Env baru
```
AUTO_MIGRATE_TIMEOUT_MIN=5
```

### Recommended deploy flow
```bash
# 1. Run migrate (block sampai exit code 0)
docker compose -f docker-compose.migrate.yaml up --abort-on-container-exit
docker compose -f docker-compose.migrate.yaml down

# 2. Deploy/restart server
docker compose -f docker-compose.be.ds.yaml up -d
```

Atau pakai Kubernetes init-container dengan image `migrate` sebelum app container start.

---

## #8. Router Split + RabbitMQ Reconnect

### Issue Router
- `app/url_mapping.go` — single file 287 baris, semua route di satu fn → merge conflict, hard to test, hard to reason.

### Fix Router
- File `app/url_mapping.go` slim down jadi orchestrator (~110 baris).
- Sub-folder `app/routes/` per domain:
  - `dashboard.go` — `RegisterDashboard`
  - `postback.go` — `RegisterPostback` (public)
  - `report.go` — `RegisterReport`
  - `internal.go` — `RegisterInternal`
  - `management.go` — `RegisterManagement` + sub-fn (`registerCampaign`, `registerMenu`, `registerRole`, `registerUser`, `registerUserLog`, `registerBudgetIO`, `registerCountryService`, `registerIPRange`, `registerCampaignSetting`)

### Issue RabbitMQ
- `cmd/server.go:31-36` — `c.Rmqp.SetUpChannel(...)` return error tapi diabaikan. Tidak ada reconnect strategy.

### Fix RabbitMQ
- `config.InitMessageBroker` switch dari `SetUpConnectionAmqp` ke `SetupConnectionAmqpAndReconnect` (lib built-in 5x retry+backoff). Return error.
- `cmd/server.go` — channel setup pakai loop + error check + `log.Fatalf` kalau initial gagal.
- Goroutine `rmqpReconnectWatcher`: blocks di `Connection.NotifyClose()`, retry connection + re-setup channel, infinite loop.

---

## File Tree Setelah Refactor (DDD-style)

```
cores/datasource/
├── main.go                                 # entry point
├── go.mod / go.sum
├── AUDIT_FIX.md                            # this doc
├── UNUSED_FUNCTIONS.md
├── Dockerfile.datasource.server
├── Dockerfile.datasource.migrate           # NEW (#7)
├── .github/workflows/main.yaml             # UPDATED matrix (#7)
├── .env
└── src/                                    # all source code dipindah ke src/ (DDD)
    ├── app/
    │   ├── url_mapping.go                  # SLIM ~110 lines (#8)
    │   └── routes/                         # NEW (#8)
    │       ├── dashboard.go
    │       ├── postback.go
    │       ├── report.go
    │       ├── internal.go
    │       └── management.go
    ├── cmd/
    │   ├── root.go                         # +migrateCmd register (#7)
    │   ├── server.go                       # +rmqp reconnect (#8); migrate dihapus (#7)
    │   └── migrate.go                      # NEW (#7) — owns migrateEntities slice
    ├── config/
    │   └── c.go                            # major refactor (#3, #4, #5, #8)
    ├── handler/
    │   ├── incoming_handler.go
    │   └── incoming_auth_handler.go        # +nbf/aud/iss/jti blacklist (#6)
    ├── helper/
    │   ├── http.go                         # TLS + nil-guard (#1, #2)
    │   ├── utils.go                        # AES env (#1)
    │   └── redis_safe.go                   # NEW (#4)
    └── domain/                             # DDD layer
        ├── entity/                         # ex /entity (data structs / GORM models)
        └── repository/                     # ex /model (DB access; package renamed)
```

### DDD restructure notes

- Module path tetap `github.com/infraLinkit/mediaplatform-datasource`. Sub-paths sekarang prefix `/src/...`:
  - `/app` → `/src/app`
  - `/cmd` → `/src/cmd`
  - `/config` → `/src/config`
  - `/handler` → `/src/handler`
  - `/helper` → `/src/helper`
  - `/entity` → `/src/domain/entity`
  - `/model` → `/src/domain/repository` (folder + package rename: `package model` → `package repository`)
- Caller `model.X` → `repository.X` (cross-file otomatis via sed).
- `main.go` tetap di root. Dockerfile build context tidak berubah (`COPY . .` + `go build .`).

---

## Env Vars Summary (yang baru)

```bash
# Security
APP_ENV=production
AES_SECRET_KEY=<32-byte>             # rotate before prod, secret manager

# DB pool
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME_MIN=30
DB_CONN_MAX_IDLE_TIME_MIN=10

# Redis
REDIS_REQUIRED=true                  # false = degraded mode

# JWT
JWT_AUDIENCE=                        # opt-in
JWT_ISSUER=                          # opt-in
AUTH_ENFORCE_DEFAULT=false           # flip true after FE migrate

# Migrate (hanya dipakai oleh `./datasource migrate`)
AUTO_MIGRATE_TIMEOUT_MIN=5
```

---

## Pre-Production Checklist

- [ ] `AES_SECRET_KEY` di-rotate, dipindah ke secret manager
- [ ] `APP_ENV=production` set di prod
- [ ] `DB_MAX_OPEN_CONNS × replicas ≤ PG max_connections`
- [ ] FE coverage test: semua admin endpoint kirim JWT
- [ ] `AUTH_ENFORCE_DEFAULT=true` setelah FE confirmed (staging dulu, monitor 401)
- [ ] `JWT_AUDIENCE` + `JWT_ISSUER` di-set kalau issuer service include claim ini
- [ ] Migration flow: wajib pakai `./datasource migrate` sub-cmd / external tool sebelum deploy `server`. Server cmd tidak menjalankan AutoMigrate.
- [ ] Tag release `v1.x.y` → CI build 2 image otomatis
- [ ] Update version di `docker-compose.migrate.yaml` + `docker-compose.be.ds.yaml` saat deploy

---

## Build Verification

Setiap step di-verify dengan `go build ./...` exit 0. No lint/vet/test introduced (audit recommends adding golangci-lint + gosec — separate task).
