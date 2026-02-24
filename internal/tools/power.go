package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerPowerTools(server *mcp.Server, h *handler) {
	registerPowerOnDedicatedServer(server, h)
	registerPowerOffDedicatedServer(server, h)
	registerPowerCycleDedicatedServer(server, h)
}

type powerArgs struct {
	ServerID string `json:"server_id" jsonschema:"dedicated server ID,required"`
}

// --- power_on_dedicated_server ---

func registerPowerOnDedicatedServer(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "power_on_dedicated_server",
		Description: "Send a power-on command to a dedicated server. The server's power_status will transition to 'powering_on', then 'powered_on'. Monitor progress with get_dedicated_server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args powerArgs) (*mcp.CallToolResult, any, error) {
		ds, err := h.client.Hosts.PowerOnDedicatedServer(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(ds)
	})
}

// --- power_off_dedicated_server ---

func registerPowerOffDedicatedServer(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "power_off_dedicated_server",
		Description: "Send a power-off command to a dedicated server. The server's power_status will transition to 'powering_off', then 'powered_off'. Monitor progress with get_dedicated_server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args powerArgs) (*mcp.CallToolResult, any, error) {
		ds, err := h.client.Hosts.PowerOffDedicatedServer(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(ds)
	})
}

// --- power_cycle_dedicated_server ---

func registerPowerCycleDedicatedServer(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "power_cycle_dedicated_server",
		Description: "Send a power-cycle (hard reboot) command to a dedicated server. The server will be powered off and back on. Monitor power_status with get_dedicated_server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args powerArgs) (*mcp.CallToolResult, any, error) {
		ds, err := h.client.Hosts.PowerCycleDedicatedServer(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(ds)
	})
}
