package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerHostTools(server *mcp.Server, h *handler) {
	registerListHosts(server, h)
}

// --- list_hosts ---

type listHostsArgs struct {
	SearchPattern string `json:"search_pattern,omitempty" jsonschema:"free-text search pattern"`
	LocationID    int64  `json:"location_id,omitempty" jsonschema:"filter by location ID"`
	Type          string `json:"type,omitempty" jsonschema:"filter by host type: dedicated_server, kubernetes_baremetal_node, sbm_server"`
	LabelSelector string `json:"label_selector,omitempty" jsonschema:"filter by labels, e.g. 'env=prod,team=infra'"`
}

func registerListHosts(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_hosts",
		Description: "List all hosts (dedicated servers, kubernetes baremetal nodes, sbm servers) with optional filtering by search pattern, location, type, and label selector",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listHostsArgs) (*mcp.CallToolResult, any, error) {
		collection := h.client.Hosts.Collection()

		if args.SearchPattern != "" {
			collection.SetParam("search_pattern", args.SearchPattern)
		}
		if args.LocationID != 0 {
			collection.SetParam("location_id", fmt.Sprintf("%d", args.LocationID))
		}
		if args.Type != "" {
			collection.SetParam("type", args.Type)
		}
		if args.LabelSelector != "" {
			collection.SetParam("label_selector", args.LabelSelector)
		}

		hosts, err := collection.Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}

		return jsonResult(hosts)
	})
}
