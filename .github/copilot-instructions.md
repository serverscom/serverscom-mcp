# Copilot Instructions

This is an MCP (Model Context Protocol) server written in Go that exposes Servers.com infrastructure management as AI-callable tools. It wraps the [serverscom-go-client](https://github.com/serverscom/serverscom-go-client) (v1.0.29) and uses the [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk).

## Project Structure

```
serverscom-mcp/
├── main.go             # CLI entry point (urfave/cli v3), MCP server init
├── version.go          # var Version = "dev" — overridden at build time via ldflags
├── Justfile            # developer tasks: build, test, snapshot, dxt
├── .goreleaser.yaml    # cross-platform release builds + ldflags version injection
├── manifest.json       # DXT manifest (Claude Desktop extension)
└── internal/tools/
    ├── tools.go            # handler struct, Register(…, version), jsonResult(), errorResult()
    ├── hosts.go            # list_hosts
    ├── dedicated_servers.go # get/update server, features (activate_*/deactivate_*)
    ├── ssh_keys.go         # account SSH keys + server attach/detach
    ├── locations.go        # order options (server models, OS, drives, uplinks, SBM, RBS flavors)
    ├── power.go            # power on/off/cycle
    ├── reinstall.go        # OS reinstall with partition/RAID layout
    ├── networks.go         # server networks (IPv4/IPv6, gateway/route)
    ├── l2_segments.go      # L2 segment management
    ├── drives.go           # drive slot listing
    └── rbs.go              # Remote Block Storage volumes
```

## How to Add a New Tool

Every tool follows the same 4-step pattern. Use existing files as reference.

### 1. Define an args struct

```go
type myResourceArgs struct {
    ResourceID string `json:"resource_id" jsonschema:"resource ID,required"`
    Name       string `json:"name,omitempty" jsonschema:"optional name"`
}
```

- Use `json:"field,required"` in the jsonschema tag for required fields
- Use `json:"field,omitempty"` for optional fields
- For tools with no arguments use `_ struct{}` as the parameter type (see `registerListL2Segments`)

### 2. Register with mcp.AddTool

```go
func registerGetMyResource(server *mcp.Server, h *handler) {
    mcp.AddTool(server, &mcp.Tool{
        Name:        "get_my_resource",
        Description: "One-line summary. Include async notes, parameter constraints, and cross-tool references here.",
    }, func(ctx context.Context, req *mcp.CallToolRequest, args myResourceArgs) (*mcp.CallToolResult, any, error) {
        result, err := h.client.MyService.Get(ctx, args.ResourceID)
        if err != nil {
            return errorResult(err), nil, nil
        }
        return jsonResult(result)
    })
}
```

- Tool names are `snake_case`
- Return signature is always `(*mcp.CallToolResult, any, error)` — the middle `any` is always `nil`
- Never return a Go error from the handler; wrap it with `errorResult(err)` instead

### 3. Add to the category register function

```go
func registerMyServiceTools(server *mcp.Server, h *handler) {
    registerGetMyResource(server, h)
    registerListMyResources(server, h)
    // ...
}
```

### 4. Call from Register() in tools.go

```go
func Register(server *mcp.Server, token, endpoint, version string) {
    // ...existing registrations...
    registerMyServiceTools(server, h)
}
```

## Client Library

The client is `*serverscom.Client` from `github.com/serverscom/serverscom-go-client/pkg`.

Top-level services:

| Service | Description |
|---|---|
| `h.client.Hosts` | `GetDedicatedServer`, `ListHosts`, etc. |
| `h.client.Hosts.DedicatedServers` | Feature management, networks, SSH keys, reinstall |
| `h.client.SSHKeys` | Account-level SSH key CRUD |
| `h.client.Locations` | Location list, order options, RBS flavors |
| `h.client.L2Segments` | L2 segment CRUD, members, networks, location groups |
| `h.client.RemoteBlockStorageVolumes` | RBS volume CRUD, credentials |

Collection endpoints return a collection object; call `.Collect(ctx)` to fetch all pages:

```go
volumes, err := h.client.RemoteBlockStorageVolumes.Collection().Collect(ctx)
```

RBS flavors live under Locations, not RemoteBlockStorageVolumes:

```go
flavors, err := h.client.Locations.RemoteBlockStorageFlavors(locationID).Collect(ctx)
```

## API Coverage

The underlying REST API is documented at:
**https://developers.servers.com/api-documentation/v1/index.json**

Key facts:
- Base URL: `https://api.servers.com` (versioned at `/v1`)
- Auth: Bearer JWT token (read-only or read-write, valid 20 years)
- Rate limit: 2,000 requests/hour; headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`
- Pagination: `X-Total` header on list responses

API tag groups map to tool files as follows:

| API tag group | Tool file |
|---|---|
| Host, Dedicated Server | `hosts.go`, `dedicated_servers.go`, `power.go`, `reinstall.go`, `networks.go`, `drives.go` |
| SSH Key | `ssh_keys.go` |
| Locations, Order options | `locations.go` |
| L2 Segment | `l2_segments.go` |
| Remote Block Storage (RBS) | `rbs.go` |

Not yet implemented: Cloud Computing, Cloud Block Storage, Load Balancers, Network Pool, Kubernetes Cluster, Racks, Billing/Account, Metrics.

## Async Operations

Many Servers.com API operations are asynchronous. When implementing tools for async endpoints:

1. Document it in the tool description: "**Async** — poll `X` until `field` reaches `value`"
2. Note what field to poll and expected transition states

| Operation type | Poll tool | Field |
|---|---|---|
| Feature activate/deactivate | `list_dedicated_server_features` | `status` |
| OS reinstall | `get_dedicated_server` | `operational_status` |
| Rescue mode | `get_dedicated_server` | `operational_status` |
| Power changes | `get_dedicated_server` | `power_status` |

`operational_status` values: `normal` → `provisioning` → `installation` → `entering_rescue_mode` → `rescue_mode` → `exiting_rescue_mode` → `maintenance`

## Conventions

- One `register<Tool>()` function per tool, one `register<Category>Tools()` per file
- All tool descriptions are written for an AI model reading them — be specific about required parameters, constraints, cross-tool dependencies, and async behaviour
- Descriptions for multi-concept tools use backtick-fenced sections or bullet lists (see `create_l2_segment` and `reinstall_dedicated_server` for examples)
- Update `TOOLS.md` when adding new tools

## Development

### Version injection

`version.go` declares `var Version = "dev"` at package `main`. GoReleaser overrides it at build time:

```
-ldflags "-s -w -X main.Version={{.Version}}"
```

`Version` is passed into `tools.Register()` and used as the HTTP `User-Agent` header (`serverscom-mcp/<version>`) and as the MCP server version field. Local builds always report `dev`.

### Common tasks (requires `just`)

| Command | Description |
|---|---|
| `just build` | Build binary for current platform (`version = "dev"`) |
| `just build-version v1.2.3` | Build with explicit version string |
| `just test` | Run all tests |
| `just check` | Validate `.goreleaser.yaml` |
| `just snapshot` | Build all platform binaries locally via GoReleaser (no publish) |
| `just dxt` | Package `serverscom-mcp.dxt` from `dist/` artifacts (run `just snapshot` first) |
