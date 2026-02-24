package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerDriveSlotTools(server *mcp.Server, h *handler) {
	registerListDedicatedServerDriveSlots(server, h)
}

// --- list_dedicated_server_drive_slots ---

type listDriveSlotArgs struct {
	ServerID string `json:"server_id" jsonschema:"dedicated server ID,required"`
}

func registerListDedicatedServerDriveSlots(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_dedicated_server_drive_slots",
		Description: "List all drive slots for a dedicated server. Returns each slot's position, interface type, form factor, and the installed drive model (if any). Useful for planning disk upgrades or understanding current storage configuration",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listDriveSlotArgs) (*mcp.CallToolResult, any, error) {
		slots, err := h.client.Hosts.DedicatedServerDriveSlots(args.ServerID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(slots)
	})
}
