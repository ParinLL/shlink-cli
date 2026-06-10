# Shlink Short-Code KV Sync

這個 container 會從 EKS 裡的 Postgres 匯出有效的 Shlink short code，並寫入 Cloudflare Workers KV。

Worker 會使用以下 key 格式：

```text
shlink:v1:<host>:<short_code>
```

必要環境變數：

```text
DATABASE_URL=postgres://user:password@host:5432/shlink?sslmode=require
CLOUDFLARE_API_TOKEN=...
CF_ACCOUNT_ID=...
CF_KV_NAMESPACE_ID=...
```

可選環境變數：

```text
SHLINK_DEFAULT_DOMAIN=parin.dev
KV_KEY_PREFIX=shlink
KV_EXPIRATION_TTL=86400
CF_KV_BATCH_SIZE=9000
```

如果你的 Shlink schema 不同，可以用 `SHLINK_SHORT_CODES_SQL` 覆蓋預設查詢。查詢結果必須依照以下順序回傳兩個不含 tab 的欄位：

```sql
SELECT domain, short_code FROM ...;
```

Build 並發布：

```bash
nerdctl.lima build --platform linux/amd64,linux/arm64 -t harbor.x300-local.parinl.com/shlink-security/shlink-edge-sync:0.3 edge/sync
nerdctl.lima push --all-platforms harbor.x300-local.parinl.com/shlink-security/shlink-edge-sync:0.3
```
