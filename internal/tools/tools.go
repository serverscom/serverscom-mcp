package tools

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

const defaultEndpoint = "https://api.servers.com/v1"

type handler struct {
	client *serverscom.Client
}

func Register(server *mcp.Server, token, endpoint string) {
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	client := serverscom.NewClientWithEndpoint(token, endpoint)
	client.SetupUserAgent("serverscom-mcp/0.1.0")

	h := &handler{client: client}

	registerHostTools(server, h)
	registerDedicatedServerTools(server, h)
	registerSSHKeyTools(server, h)
	registerLocationTools(server, h)
	registerDriveSlotTools(server, h)
	registerPowerTools(server, h)
	registerReinstallTools(server, h)
	registerNetworkTools(server, h)
	registerL2SegmentTools(server, h)
	registerRBSTools(server, h)
}

// jsonResult serialises v as indented JSON and wraps it in a TextContent result.
func jsonResult(v interface{}) (*mcp.CallToolResult, any, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// errorResult wraps an error as a failed tool result.
func errorResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
		},
		IsError: true,
	}
}

