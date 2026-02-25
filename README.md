# serverscom-mcp

MCP (Model Context Protocol) server for managing [Servers.com](https://servers.com) dedicated server infrastructure. Enables AI assistants (Claude, etc.) to interact with the Servers.com API directly — query servers, manage SSH keys, configure networks, manage L2 segments, provision Remote Block Storage, reinstall operating systems, and more.

## Quick Start

```bash
npx @servers.com/mcp --token your-api-token

# or via env var
SC_TOKEN=your-api-token npx @servers.com/mcp
```

## Usage with Claude Desktop

```json
{
  "mcpServers": {
    "serverscom": {
      "command": "npx",
      "args": ["-y", "@servers.com/mcp"],
      "env": {
        "SC_TOKEN": "your-api-token"
      }
    }
  }
}
```

## Configuration

| Flag | Env var | Required | Default | Description |
|---|---|---|---|---|
| `--token`, `-t` | `SC_TOKEN` | yes | — | Servers.com API token |
| `--endpoint`, `-e` | `SC_ENDPOINT` | no | `https://api.servers.com/v1` | Custom API endpoint |

## Available Tools

78 tools across 7 categories — see **[TOOLS.md](TOOLS.md)** for the full reference.

## Async Operations

Many operations are asynchronous. After calling them, poll the relevant status field:

| Operation | Poll with | Field to watch |
|---|---|---|
| Feature changes | `list_dedicated_server_features` | feature `status` |
| Rescue mode | `get_dedicated_server` | `operational_status` |
| OS reinstallation | `get_dedicated_server` | `operational_status` |
| Power changes | `get_dedicated_server` | `power_status` |

`operational_status` values: `normal` → `provisioning` → `installation` → `entering_rescue_mode` → `rescue_mode` → `exiting_rescue_mode` → `maintenance`

## License

[MIT](LICENSE)

---

## Development

### Building from source

```bash
go build -o serverscom-mcp .
```

### Releasing

Releases are automated via [GoReleaser](https://goreleaser.com) and GitHub Actions. Push a version tag to trigger the pipeline:

```bash
git tag v1.2.3
git push origin v1.2.3
```

The workflow will:
1. Build binaries for Linux, macOS, Windows (amd64 + arm64)
2. Create a GitHub Release with archives and checksums
3. Publish `@servers.com/mcp` and platform packages to npm using [trusted publishing](https://docs.npmjs.com/trusted-publishers/) (OIDC, no long-lived tokens)

### Project Structure

```
serverscom-mcp/
├── main.go                       # Entry point
└── internal/tools/
    ├── tools.go                  # Tool registration hub, shared helpers
    ├── hosts.go                  # list_hosts
    ├── dedicated_servers.go      # Server CRUD and feature management
    ├── ssh_keys.go               # SSH key operations
    ├── locations.go              # Location and infrastructure options
    ├── power.go                  # Power management
    ├── drives.go                 # Drive slot listing
    ├── reinstall.go              # OS reinstallation
    ├── networks.go               # Network management
    ├── l2_segments.go            # L2 segment management
    └── rbs.go                    # Remote Block Storage
```
