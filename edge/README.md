# Shlink Edge Guard

這個目錄包含一個 Cloudflare Worker，以及一個部署在 EKS 的 CronJob。用途是把無效的 Shlink short-code request 擋在 Cloudflare edge，避免 request 進到 Shlink 後被記成 orphan visit。

流程：

```text
EKS Postgres -> sync CronJob -> Cloudflare Workers KV -> Worker on parin.dev/* -> Shlink origin
```

Worker 只檢查可能是短網址的路徑。其他路徑會查詢 KV 裡的 `shlink:v2:<host>` manifest；manifest 內沒有對應 short code 時，Worker 會直接在 Cloudflare 回 `404`，Shlink 就不會記錄 orphan visit。

## Cloudflare 設定

1. 建立 Workers KV namespace，例如 `shlink_short_codes`。
2. 把 namespace ID 填到 `edge/worker/wrangler.toml`。
3. 部署 Worker：

```bash
cd edge/worker
npm install
npm run deploy
```

Worker route 已設定為 `parin.dev/*`，zone 是 `parin.dev`。

### Cloudflare Custom Rule

Worker + KV 是主要的 allowlist 防護；Cloudflare Custom Rule 可以保留在 Worker 前面，用來先擋掉明顯不需要進 Worker 的 path，降低 Worker invocation 與 Shlink origin 流量。

目前建議的 Block rule：

```text
http.host eq "parin.dev" and (
  http.request.uri.path eq "/" or
  http.request.uri.path wildcard r"/*/*" or
  http.request.uri.path wildcard r"*.*"
)
```

這條 rule 的行為：

```text
/       -> Cloudflare 直接 block，不進 Worker
/*/*    -> Cloudflare 直接 block，不進 Worker
*.*     -> Cloudflare 直接 block，不進 Worker
```

如果未來需要讓 `https://parin.dev/` 做首頁或 base redirect，請移除 `http.request.uri.path eq "/"`。如果需要支援 `/shortCode/extra/path`，請移除 `http.request.uri.path wildcard r"/*/*"`，讓 Worker 依第一段 short code 查 KV 後決定是否放行。

## EKS 設定

1. Build 並 push 同步用 image：

```bash
nerdctl.lima build --platform linux/amd64,linux/arm64 -t harbor.x300-local.parinl.com/shlink-security/shlink-edge-sync:0.6 edge/sync
nerdctl.lima push --all-platforms harbor.x300-local.parinl.com/shlink-security/shlink-edge-sync:0.6
```

2. 依照 `edge/k8s/secret.example.yaml` 建立同步用 Secret。必要值如下：

```text
DATABASE_URL
CLOUDFLARE_API_TOKEN
CF_ACCOUNT_ID
CF_KV_NAMESPACE_ID
```

Cloudflare token 需要有寫入該 account Workers KV 的權限。

3. 在每一台 k3s node 設定 Harbor registry 權限。k3s 使用 containerd 拉 image，不會讀 Docker CLI 的 login 狀態。請在每一台可能排程 workload 的 server/agent node 建立 `/etc/rancher/k3s/registries.yaml`：

```yaml
mirrors:
  harbor.x300-local.parinl.com:
    endpoint:
      - "https://harbor.x300-local.parinl.com"
configs:
  harbor.x300-local.parinl.com:
    auth:
      username: "<harbor-robot-username>"
      password: "<harbor-robot-token>"
```

建議在 Harbor 建立 project-scoped robot account，只給 `shlink-security` project 的 pull 權限。

修改後需要在每一台 node 重啟 k3s：

```bash
sudo systemctl restart k3s
# agent node 使用：
sudo systemctl restart k3s-agent
```

4. 套用 ConfigMap 與 CronJob：

```bash
kubectl apply -f edge/k8s/configmap.yaml
kubectl apply -f edge/k8s/secret.example.yaml
kubectl apply -f edge/k8s/cronjob.yaml
```

正式環境請用你既有的 secret 管理方式取代 `secret.example.yaml`，不要直接提交真實 token 或 database password。

## 驗證

手動跑一次同步 job：

```bash
kubectl -n shlink create job --from=cronjob/shlink-edge-sync shlink-edge-sync-manual
kubectl -n shlink logs job/shlink-edge-sync-manual
```

接著測試：

```bash
curl -i https://parin.dev/<known-short-code>
curl -i https://parin.dev/api
curl -i https://parin.dev/graphql
```

已存在的 short code 應該會到 Shlink。未知或掃描用路徑應該會直接回 `404`，並帶有：

```text
x-shlink-edge-guard: blocked
```

## KV 免費方案用量

同步器不再為每個 short code 建立和刷新獨立 KV key。每個 domain 只使用一個 manifest：

```text
shlink:v2:<host>
```

CronJob 每 5 分鐘讀取現有 manifest 並與 Postgres 結果比較。內容沒有變化時不執行 PUT；新增或刪除 short code 時，每個受影響 domain 只會產生一次 PUT。這可避免超過 Workers KV 免費方案每天 1,000 次 write 的限制。

為避免 Postgres query 意外回傳空資料時清空 allowlist，同步器在完全沒有 short code 時會跳過更新。

## 注意事項

如果你的 Shlink DB schema 跟預設 query 不同，請在 `edge/k8s/configmap.yaml` 設定 `SHLINK_SHORT_CODES_SQL`。查詢結果必須依序回傳 `domain` 與 `short_code` 兩個欄位。
