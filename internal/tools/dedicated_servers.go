package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

func registerDedicatedServerTools(server *mcp.Server, h *handler) {
	registerGetDedicatedServer(server, h)
	registerUpdateDedicatedServer(server, h)
	registerListDedicatedServerFeatures(server, h)
	registerFeatureTools(server, h)
}

// --- get_dedicated_server ---

type getDedicatedServerArgs struct {
	ServerID string `json:"server_id" jsonschema:"dedicated server ID,required"`
}

func registerGetDedicatedServer(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_dedicated_server",
		Description: "Get detailed information about a dedicated server including configuration, status, operational_status, and power_status. Operational status indicates async operations: normal, provisioning, installation, entering_rescue_mode, exiting_rescue_mode, rescue_mode, maintenance",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getDedicatedServerArgs) (*mcp.CallToolResult, any, error) {
		ds, err := h.client.Hosts.GetDedicatedServer(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(ds)
	})
}

// --- update_dedicated_server ---

type updateDedicatedServerArgs struct {
	ServerID   string            `json:"server_id" jsonschema:"dedicated server ID,required"`
	Title      string            `json:"title,omitempty" jsonschema:"server title in the Servers.com portal (does not rename the OS hostname)"`
	Labels     map[string]string `json:"labels,omitempty" jsonschema:"labels to attach to the server (replaces existing labels)"`
	UserData   *string           `json:"user_data,omitempty" jsonschema:"cloud-init user data processed during next initialization (Linux: cloud-init, Windows: cloudbase-init). Omit to leave unchanged; send empty string to clear"`
	IPXEConfig *string           `json:"ipxe_config,omitempty" jsonschema:"iPXE configuration (under development, currently unavailable). Omit to leave unchanged; send empty string to clear"`
}

func registerUpdateDedicatedServer(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_dedicated_server",
		Description: "Update a dedicated server's title, labels, user_data (cloud-init), and/or ipxe_config. All fields are optional — only provided fields are changed",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateDedicatedServerArgs) (*mcp.CallToolResult, any, error) {
		ds, err := h.client.Hosts.UpdateDedicatedServer(ctx, args.ServerID, serverscom.DedicatedServerUpdateInput{
			Title:      args.Title,
			Labels:     args.Labels,
			UserData:   args.UserData,
			IPXEConfig: args.IPXEConfig,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(ds)
	})
}

// --- list_dedicated_server_features ---

type listDedicatedServerFeaturesArgs struct {
	ServerID string `json:"server_id" jsonschema:"dedicated server ID,required"`
}

func registerListDedicatedServerFeatures(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_dedicated_server_features",
		Description: "List features for a dedicated server. Each feature has a name and status (activation, activated, deactivation, deactivated, incompatible, unavailable). Features: disaggregated_public_ports, disaggregated_private_ports, no_public_ip_address, no_private_ip, host_rescue_mode, oob_public_access, no_public_network, private_ipxe_boot",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listDedicatedServerFeaturesArgs) (*mcp.CallToolResult, any, error) {
		features, err := h.client.Hosts.DedicatedServerFeatures(args.ServerID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(features)
	})
}

// --- feature activate/deactivate ---

type featureArgs struct {
	ServerID string `json:"server_id" jsonschema:"dedicated server ID,required"`
}

type rescueModeActivateArgs struct {
	ServerID           string   `json:"server_id" jsonschema:"dedicated server ID,required"`
	AuthMethods        []string `json:"auth_methods" jsonschema:"authentication methods: password and/or ssh_key,required"`
	SSHKeyFingerprints []string `json:"ssh_key_fingerprints,omitempty" jsonschema:"SSH key fingerprints, required when ssh_key auth method is selected"`
}

type ipxeBootActivateArgs struct {
	ServerID   string `json:"server_id" jsonschema:"dedicated server ID,required"`
	IPXEConfig string `json:"ipxe_config" jsonschema:"iPXE configuration string,required"`
}

func registerFeatureTools(server *mcp.Server, h *handler) {
	// disaggregated_public_ports
	mcp.AddTool(server, &mcp.Tool{
		Name:        "activate_disaggregated_public_ports",
		Description: "Activate disaggregated public ports feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'activated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.ActivateDisaggregatedPublicPortsFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deactivate_disaggregated_public_ports",
		Description: "Deactivate disaggregated public ports feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'deactivated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.DeactivateDisaggregatedPublicPortsFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	// disaggregated_private_ports
	mcp.AddTool(server, &mcp.Tool{
		Name:        "activate_disaggregated_private_ports",
		Description: "Activate disaggregated private ports feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'activated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.ActivateDisaggregatedPrivatePortsFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deactivate_disaggregated_private_ports",
		Description: "Deactivate disaggregated private ports feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'deactivated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.DeactivateDisaggregatedPrivatePortsFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	// no_public_ip_address
	mcp.AddTool(server, &mcp.Tool{
		Name: "activate_no_public_ip_address",
		Description: `Activate no_public_ip_address feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'activated'.

When active, no default public IP network is assigned to the server. Implications:
- Additional public networks (gateway or route) cannot be ordered
- Features requiring a public IP (e.g. host_rescue_mode, oob_public_access) become unavailable
- The public network interface itself still exists and is configured, so the server can be added to a public L2 segment or Firewall`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.ActivateNoPublicIpAddressFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deactivate_no_public_ip_address",
		Description: "Deactivate no_public_ip_address feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'deactivated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.DeactivateNoPublicIpAddressFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	// no_private_ip
	mcp.AddTool(server, &mcp.Tool{
		Name:        "activate_no_private_ip",
		Description: "Activate no_private_ip feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'activated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.ActivateNoPrivateIpFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deactivate_no_private_ip",
		Description: "Deactivate no_private_ip feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'deactivated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.DeactivateNoPrivateIpFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	// oob_public_access
	mcp.AddTool(server, &mcp.Tool{
		Name:        "activate_oob_public_access",
		Description: "Activate OOB public access feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'activated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.ActivateOobPublicAccessFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deactivate_oob_public_access",
		Description: "Deactivate OOB public access feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'deactivated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.DeactivateOobPublicAccessFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	// no_public_network
	mcp.AddTool(server, &mcp.Tool{
		Name: "activate_no_public_network",
		Description: `Activate no_public_network feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'activated'.

More radical than no_public_ip_address: not only is no public IP assigned, but the public network interface is not configured at all. Implications:
- All restrictions of no_public_ip_address apply (no additional public networks, no rescue mode, no OOB public access)
- The server cannot be added to a public L2 segment or Firewall
- Any functionality that assumes the existence of a public interface is unavailable`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.ActivateNoPublicNetworkFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deactivate_no_public_network",
		Description: "Deactivate no_public_network feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'deactivated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.DeactivateNoPublicNetworkFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	// host_rescue_mode — activate requires auth_methods
	mcp.AddTool(server, &mcp.Tool{
		Name:        "activate_host_rescue_mode",
		Description: "Activate host rescue mode for a dedicated server. Async — operational_status changes to 'entering_rescue_mode', then 'rescue_mode'. Poll with get_dedicated_server or list_dedicated_server_features",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args rescueModeActivateArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.ActivateHostRescueModeFeature(ctx, args.ServerID, serverscom.HostRescueModeFeatureInput{
			AuthMethods:        args.AuthMethods,
			SSHKeyFingerprints: args.SSHKeyFingerprints,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deactivate_host_rescue_mode",
		Description: "Deactivate host rescue mode for a dedicated server. Async — operational_status changes to 'exiting_rescue_mode', then back to 'normal'. Poll with get_dedicated_server or list_dedicated_server_features",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.DeactivateHostRescueModeFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	// private_ipxe_boot — activate requires ipxe_config
	mcp.AddTool(server, &mcp.Tool{
		Name:        "activate_private_ipxe_boot",
		Description: "Activate private iPXE boot feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'activated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args ipxeBootActivateArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.ActivatePrivateIpxeBootFeature(ctx, args.ServerID, serverscom.PrivateIpxeBootFeatureInput{
			IPXEConfig: args.IPXEConfig,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deactivate_private_ipxe_boot",
		Description: "Deactivate private iPXE boot feature for a dedicated server. Async — poll list_dedicated_server_features until status is 'deactivated'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args featureArgs) (*mcp.CallToolResult, any, error) {
		result, err := h.client.Hosts.DeactivatePrivateIpxeBootFeature(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(result)
	})
}
