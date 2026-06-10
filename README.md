# shlink-cli

A command-line tool for the Shlink URL shortener, built on the [Shlink REST API v3](https://api-spec.shlink.io/).

## Use With OpenClaw

This repository includes a publish-ready skill file at [`clawhub-publish/SKILL.md`](./clawhub-publish/SKILL.md).

ClawHub page:

- https://clawhub.ai/ParinLL/shlink-cli

Install from ClawHub:

```bash
clawhub install shlink-cli
```

This is an instruction-only skill package (contains only `SKILL.md`) and guides users to install `shlink-cli` from GitHub.

## Environment Variables

```bash
export SHLINK_BASE_URL=https://your-shlink-instance.example.com
export SHLINK_API_KEY=your-api-key-here
```

## Installation (Go Build)

Requires Go 1.24+.

```bash
cd shlink-cli
go mod tidy
go build -o shlink-cli .
```

Install to PATH:

```bash
go build -o shlink-cli .
sudo mv shlink-cli /usr/local/bin/
```

Or use `go install`:

```bash
go install github.com/ParinLL/shlink-cli@latest
```

## Docker Usage

### Build Image (Multi-arch)

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t shlink-cli .
```

### Using Docker Compose

```bash
cp .env.example .env
# Edit .env and fill in your SHLINK_BASE_URL and SHLINK_API_KEY

# Run commands
docker compose run --rm shlink-cli short-url list
docker compose run --rm shlink-cli short-url create https://example.com
docker compose run --rm shlink-cli health
```

### Using Docker Directly

```bash
docker run --rm \
  -e SHLINK_BASE_URL=https://your-shlink.example.com \
  -e SHLINK_API_KEY=your-api-key \
  shlink-cli short-url list
```

## Shlink Edge Guard

本 repo 也包含一組 Cloudflare Worker + EKS/k3s CronJob，用來把無效的 Shlink short-code request 擋在 Cloudflare edge，避免 Shlink 記錄大量 orphan visits。

詳細部署文件在 [`edge/README.md`](./edge/README.md)。

架構：

```text
EKS Postgres -> sync CronJob -> Cloudflare Workers KV -> Worker on parin.dev/* -> Shlink origin
```

Worker 會用 Cloudflare KV 作為 short-code allowlist：

```text
有效 short code -> 放行到 Shlink
無效 short code -> Cloudflare Worker 直接回 404
```

同步器每 5 分鐘從 Shlink Postgres 匯出有效 short code，寫入 KV key：

```text
shlink:v1:<host>:<short_code>
```

### Edge Guard 部署摘要

部署 Worker：

```bash
cd edge/worker
npm install
npm run deploy
```

Build 並 push 同步器 image：

```bash
nerdctl.lima build --platform linux/amd64,linux/arm64 -t harbor.x300-local.parinl.com/shlink-security/shlink-edge-sync:0.3 edge/sync
nerdctl.lima push --all-platforms harbor.x300-local.parinl.com/shlink-security/shlink-edge-sync:0.3
```

套用 k3s/EKS 資源：

```bash
kubectl apply -f edge/k8s/configmap.yaml
kubectl apply -f edge/k8s/secret.example.yaml
kubectl apply -f edge/k8s/cronjob.yaml
```

正式環境請使用你自己的 secret 管理方式取代 `secret.example.yaml`，不要提交真實 database password 或 Cloudflare token。

### Cloudflare Custom Rule

Worker + KV 是主要 allowlist 防護；Cloudflare Custom Rule 可以保留在 Worker 前面，先擋掉明顯不需要進 Worker 的 path，降低 Worker invocation 與 Shlink origin 流量。

建議的 Block rule：

```text
http.host eq "parin.dev" and (
  http.request.uri.path eq "/" or
  http.request.uri.path wildcard r"/*/*" or
  http.request.uri.path wildcard r"*.*"
)
```

如果需要讓 `https://parin.dev/` 做首頁或 base redirect，請移除 `http.request.uri.path eq "/"`。如果需要支援 `/shortCode/extra/path`，請移除 `http.request.uri.path wildcard r"/*/*"`，讓 Worker 依第一段 short code 查 KV 後決定是否放行。

### Edge Guard 驗證

手動跑一次同步 job：

```bash
kubectl -n shlink create job --from=cronjob/shlink-edge-sync shlink-edge-sync-manual
kubectl -n shlink logs job/shlink-edge-sync-manual
```

測試 Worker：

```bash
curl -i https://parin.dev/<known-short-code>
curl -i https://parin.dev/api
curl -i https://parin.dev/graphql
```

未知或掃描用路徑應該回：

```text
x-shlink-edge-guard: blocked
```

## Testing

```bash
go test ./...
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--debug` | Enable debug output (API Key will be masked) |
| `--help` | Show help |
| `--base-url` | Shlink base URL (or set `SHLINK_BASE_URL`) |
| `--api-key` | Shlink API Key (or set `SHLINK_API_KEY`) |

## Commands

### Short URLs

```bash
shlink-cli short-url list [--page N] [--per-page N] [--search TERM] [--tags tag1,tag2] [--json]
shlink-cli short-url create <longUrl> [--slug SLUG] [--tags tag1,tag2] [--title TITLE] [--json]
shlink-cli short-url get <shortCode> [--domain DOMAIN] [--json]
shlink-cli short-url edit <shortCode> [--long-url URL] [--title TITLE] [--tags tag1,tag2]
shlink-cli short-url delete <shortCode> [--domain DOMAIN]
```

### Tags

```bash
shlink-cli tag list [--stats] [--json]
shlink-cli tag rename <oldName> <newName>
shlink-cli tag delete <tag1> [tag2...]
```

### Visits

```bash
shlink-cli visit overview [--json]
shlink-cli visit short-url <shortCode> [--domain DOMAIN] [--page N] [--per-page N] [--json]
shlink-cli visit tag <tag> [--page N] [--per-page N] [--json]
shlink-cli visit orphan [--page N] [--per-page N] [--json]
```

### Domains

```bash
shlink-cli domain list [--json]
shlink-cli domain set-redirects <domain> [--base-url-redirect URL] [--404-redirect URL] [--invalid-redirect URL]
```

### Health

```bash
shlink-cli health [--json]
```

## License

This project is licensed under the [Apache License 2.0](LICENSE).
