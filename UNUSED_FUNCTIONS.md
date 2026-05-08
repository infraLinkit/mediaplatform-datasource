# Unused / Redundant Function Audit

Tanggal scan: 2026-05-08
Tool: `staticcheck -checks U1000` + manual grep (untuk exported funcs).

## 1. REMOVED (private/unused)

Fungsi-fungsi privat berikut **sudah dihapus** karena dead code (tidak ada caller):

| File | Function | Alasan |
|------|----------|--------|
| `handler/incoming_campaign_summary_handler.go` | `getPreviousValueRevenue` | duplicate dari `getPreviousValueRedirection` di file lain, tidak dipakai |
| `handler/incoming_campaign_summary_handler.go` | `countPercentageRevenue` | duplicate dari `countPercentageRedirection`, tidak dipakai |
| `handler/incoming_campaign_summary_handler.go` | `countBudgetUsage` | tidak dipakai + ada **bug** typo: cek pakai `target_daily_budgets` (plural) tapi divide pakai `target_daily_budget` (singular) |
| `handler/incoming_redirection_time_handler.go` | `groupOperatorRedirection` | unused, chain caller (group → group → group → group) tidak ada entry point |
| `handler/incoming_redirection_time_handler.go` | `groupPartnerRedirection` | unused, dipanggil cuma dari `groupOperatorRedirection` (juga unused) |
| `handler/incoming_redirection_time_handler.go` | `groupServiceRedirection` | unused, dipanggil cuma dari `groupPartnerRedirection` (juga unused) |
| `handler/incoming_redirection_time_handler.go` | `groupAdnetRedirection` | unused, dipanggil cuma dari `groupServiceRedirection` (juga unused) |
| `handler/incoming_redirection_time_handler.go` | `countAverageRedirectionHourly` | unused, ada `countAverageRedirection` non-hourly yang dipakai |

**Verifikasi**: `staticcheck -checks U1000 ./...` setelah remove → 0 warnings. `go build ./...` exit 0.

---

## 2. UNUSED EXPORTED FUNCTIONS (KEPT, recommend review)

Fungsi exported (Capitalized) yang tidak dipanggil dimanapun di module `mediaplatform`. **Tidak dihapus** karena:
- Exported = bisa dipanggil dari module/repo lain (gak bisa staticcheck verify cross-module)
- Mungkin sengaja disediakan sebagai public API utility
- Bisa dipakai via reflection / struct embedding

Recommend manual review untuk decide: hapus, deprecate, atau dokumentasikan.

### 2.1. helper/utils.go

| Function | Signature | Catatan |
|----------|-----------|---------|
| `Encrypt` | `(plaintext string) (string, error)` | AES-GCM encrypt. **0 callers** di datasource, lp, worker, cms. Sudah di-refactor di #1 (env-based key). Aman dihapus kalau memang tidak dipakai. |
| `Decrypt` | `(ciphertext string) (string, error)` | Pair dari Encrypt. **0 callers**. Hapus barengan Encrypt. |
| `CompressGzip` | `(f string) string` | Gzip file utility. **0 callers**. Mungkin legacy? |
| `GetYesterday` | `(loc, day, layout) time.Time` | **0 callers**. `GetCurrentTime` (84 callers) udah handle relatif date. |
| `GetTomorrow` | `(loc, day, layout) time.Time` | **0 callers**. Sama spt GetYesterday. |
| `GetTrxID` | `(loc) string` | **0 callers**. `GetUniqId` (6 callers) lebih populer. |
| `GetDateFromInt` | `(loc, dateInt) time.Time` | **0 callers**. |

### 2.2. helper/http.go

| Function | Catatan |
|----------|---------|
| `Throw` | exception-pattern util. **0 callers**. Designed untuk dipakai inside `Block.Try()` tapi tidak ada code yang call Throw. |
| `HttpDial` | `(url, timeout) error` — TCP dial check. **0 callers**. |
| `HttpDial2` | `(url, timeout) bool` — HTTP GET check. **0 callers**. Selain itu, fn ini leak goroutine + ignore response body — hapus saja. |

### 2.3. helper/redis_safe.go (BARU — #4)

Fungsi baru yang ditambahkan saat fix #4 (degraded mode). **Belum dipakai**, sengaja disediakan untuk migrate hot path:

| Function | Catatan |
|----------|---------|
| `SafeRedisGet` | nil-safe Get. Migration target: replace `h.RCP.Get(key).Val()` |
| `SafeRedisSet` | nil-safe Set |
| `SafeRedisDel` | nil-safe Del |
| `SafeRueidisJSONGet` | nil-safe JSON.GET |
| `SafeRueidisExists` | nil-safe EXISTS |

**Status**: Intentionally unused. Akan dipakai progressive saat hot path dimigrate. Lihat `AUDIT_FIX.md` #4 untuk detail.

---

## 3. STATICCHECK WARNINGS LAIN (bukan unused)

Issues yang muncul di `staticcheck ./...` selain U1000. Bukan task ini, tapi flag untuk follow-up:

| File:Line | Code | Issue |
|-----------|------|-------|
| `handler/incoming_budget_target_handler.go:81` | SA5009 | `Printf` format `%s` dipakai dgn arg type `entity.TargetBudgetDetail` (struct, bukan string) |
| `handler/incoming_campaign_summary_handler.go:127` | SA4006 | value `BudgetDetailPerMonth` di-assign tapi tidak dipakai |
| `handler/incoming_campaign_summary_handler.go:128` | SA4006 | value `TargetBudget` di-assign tapi tidak dipakai |
| `model/dashboard.go:400` | S1002 | `dsp.IsDsp == true` → simplify ke `dsp.IsDsp` |
| `model/dashboard.go:417` | SA4006 | `where_non_dsp` di-assign tapi tidak dipakai |
| `model/dashboard.go:494` | SA4010 | `append` result tidak dipakai (kemungkinan bug — append harus di-reassign) |

⚠️ `model/dashboard.go:494` SA4010 = potential bug (append yang tidak di-reassign biasanya silent data loss).

---

## 4. SUMMARY

| Kategori | Count | Action |
|----------|-------|--------|
| Removed (private unused) | 8 | ✅ Done — `git diff` to review |
| Unused exported (legacy) | 10 | 📝 Review manual, putuskan hapus / deprecate |
| Unused exported (intentional, baru) | 5 | ⏳ Akan dipakai saat migration #4 |
| Other staticcheck issues | 6 | 🔍 Follow-up separate task |

---

## 5. CARA RE-VERIFY

```bash
# Re-scan unused private/unexported
cd cores/datasource
staticcheck -checks U1000 ./...

# Full staticcheck (all checks)
staticcheck ./...

# Cek caller fungsi exported tertentu
grep -rn "helper\.Encrypt" --include="*.go" /path/to/mediaplatform

# Build verify
go build ./...
```

## 6. NEXT STEPS REKOMENDASI

1. **Confirm tidak dipakai cross-repo** untuk 10 fungsi exported di section 2.1 + 2.2. Cek di repo lain (cms, fe, dll) yang import module ini.
2. **Bulk delete** kalau confirmed (Encrypt/Decrypt, HttpDial*, Throw, Get(Yesterday|Tomorrow|TrxID|DateFromInt), CompressGzip).
3. **Fix bug** di `model/dashboard.go:494` (SA4010 append result lost).
4. **Setup CI** untuk fail build pada U1000 baru — prevent regression. Lihat audit PDF section 10 (Suggested Tools): `golangci-lint` dgn preset `unused`.
