# serverscom-mcp

MCP (Model Context Protocol) server for managing [Servers.com](https://servers.com) dedicated server infrastructure. Enables AI assistants (Claude, etc.) to interact with the Servers.com API directly — query servers, manage SSH keys, configure networks, reinstall operating systems, and more.

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

### Hosts

| Tool | Description |
|---|---|
| `list_hosts` | List all hosts (dedicated servers, K8s baremetal nodes, SBM) with optional filtering by search pattern, location, type, and labels |

### Dedicated Servers

| Tool | Description |
|---|---|
| `get_dedicated_server` | Get detailed server info: configuration, status, `operational_status`, `power_status` |
| `update_dedicated_server` | Update title, labels, `user_data` (cloud-init), or `ipxe_config` |
| `list_dedicated_server_features` | List features and their status for a server |

#### Features

Each feature supports `activate_*` and `deactivate_*` operations. All feature changes are **asynchronous** — poll `list_dedicated_server_features` until the desired status is reached.

| Feature | Description |
|---|---|
| `disaggregated_public_ports` | Disaggregated public ports |
| `disaggregated_private_ports` | Disaggregated private ports |
| `no_public_ip_address` | No default public IP; restricts rescue mode, OOB access, additional public networks |
| `no_private_ip` | No default private IP |
| `no_public_network` | Public interface is not configured at all (more restrictive than `no_public_ip_address`) |
| `host_rescue_mode` | Rescue mode (requires `auth_methods` and optionally `ssh_key_fingerprints`) |
| `oob_public_access` | Out-of-band public access |
| `private_ipxe_boot` | Private iPXE boot (requires `ipxe_config`) |

### SSH Keys

| Tool | Description |
|---|---|
| `list_ssh_keys` | List all SSH keys in the account |
| `get_ssh_key` | Get key details by fingerprint |
| `create_ssh_key` | Create a new SSH key |
| `update_ssh_key` | Update key name or labels |
| `delete_ssh_key` | Delete an SSH key |
| `list_dedicated_server_ssh_keys` | List keys attached to a server |
| `attach_ssh_keys_to_dedicated_server` | Attach SSH keys to a server |
| `detach_ssh_key_from_dedicated_server` | Detach an SSH key from a server |

### Power Management

All power operations are **asynchronous** — poll `get_dedicated_server` and check `power_status`.

| Tool | Description |
|---|---|
| `power_on_dedicated_server` | Power on |
| `power_off_dedicated_server` | Power off |
| `power_cycle_dedicated_server` | Hard reboot |

### OS Reinstallation

| Tool | Description |
|---|---|
| `reinstall_dedicated_server` | Reinstall OS with custom partition layout, RAID config, and SSH key injection. **Async** — `operational_status` transitions to `installation`, then back to `normal` |

Supports custom partition layouts and RAID configuration. Max 1 SSH key per reinstall.

### Networks

| Tool | Description |
|---|---|
| `list_dedicated_server_networks` | List all networks (public/private, IPv4/IPv6) attached to a server |
| `get_dedicated_server_network` | Get network details: CIDR, family, interface type, distribution method, status |
| `get_dedicated_server_network_usage` | Check current and committed bandwidth utilization |
| `add_dedicated_server_public_ipv4_network` | Add additional public IPv4 network (`gateway` or `route`) |
| `add_dedicated_server_private_ipv4_network` | Add additional private IPv4 network (`gateway` or `route`) |
| `activate_dedicated_server_public_ipv6_network` | Activate IPv6 (one allocation per server) |
| `delete_dedicated_server_network` | Remove an additional network (default network cannot be deleted) |

#### Network Quotas

- **Route (alias) networks**: max 32 IPs per server
- **Additional gateway networks per family** (public/private IPv4): max 2 additional
- **Total IPs across all gateway networks per family**: max 72 (e.g. if default is /29 = 8 IPs, only 64 remain)
- **IPv6**: max 1 allocation per server (either a single /64 or a /125 + /64 depending on location)

### Locations & Infrastructure Options

| Tool | Description |
|---|---|
| `list_locations` / `get_location` | Data center locations |
| `list_server_model_options` / `get_server_model_option` | Available server models per location |
| `list_ram_options` | RAM upgrade options |
| `list_os_options` / `get_os_option` | Available operating systems |
| `list_drive_model_options` / `get_drive_model_option` | Available drive models |
| `list_uplink_options` / `get_uplink_option` | Network uplink options |
| `list_bandwidth_options` / `get_bandwidth_option` | Bandwidth plans |
| `list_sbm_flavor_options` / `get_sbm_flavor_option` | Scalable Bare Metal (SBM) flavors |
| `list_sbm_os_options` | OS options for SBM servers |

### Storage

| Tool | Description |
|---|---|
| `list_dedicated_server_drive_slots` | List all drive slots and installed drives for a server |

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
    └── networks.go               # Network management
```
