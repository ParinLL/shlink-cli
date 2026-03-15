# shlink-cli

A command-line tool for the Shlink URL shortener, built on the [Shlink REST API v3](https://api-spec.shlink.io/).

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
go install shlink-cli@latest
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
