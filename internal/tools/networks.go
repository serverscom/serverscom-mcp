package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

func registerNetworkTools(server *mcp.Server, h *handler) {
	registerListDedicatedServerNetworks(server, h)
	registerGetDedicatedServerNetwork(server, h)
	registerGetDedicatedServerNetworkUsage(server, h)
	registerAddDedicatedServerPublicIPv4Network(server, h)
	registerAddDedicatedServerPrivateIPv4Network(server, h)
	registerActivateDedicatedServerPublicIPv6Network(server, h)
	registerDeleteDedicatedServerNetwork(server, h)
}

// --- list_dedicated_server_networks ---

type listNetworksArgs struct {
	ServerID string `json:"server_id" jsonschema:"dedicated server ID,required"`
}

func registerListDedicatedServerNetworks(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_dedicated_server_networks",
		Description: `List all networks attached to a dedicated server.

Each network has a distribution_method:
- "gateway" — a fully routable subnet with a gateway; all IPs are independently addressable
- "route" — /32 addresses configured as interface aliases; not available in all locations

IPv4 quotas (per host):
- Route (alias) networks: max 32 IPs total
- Additional gateway networks per family (public IPv4 / private IPv4): max 2 additional (1 default is usually present, but may be absent if features like no_public_ip, no_private_ip, or no_public_network are active)
- Total IPs across all gateway networks per family: max 72 (e.g. if the default network is /29 = 8 IPs, only 64 remain, so the largest additional prefix you can order is /26)

IPv6:
- Only one IPv6 allocation per server is possible
- Depending on the location, IPv6 is provisioned either as a single /64, or as a /125 (point-to-point) with a routed /64 — both are delivered as part of the same activation and count as one allocation`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listNetworksArgs) (*mcp.CallToolResult, any, error) {
		networks, err := h.client.Hosts.DedicatedServerNetworks(args.ServerID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(networks)
	})
}

// --- get_dedicated_server_network ---

type getNetworkArgs struct {
	ServerID  string `json:"server_id" jsonschema:"dedicated server ID,required"`
	NetworkID string `json:"network_id" jsonschema:"network ID,required"`
}

func registerGetDedicatedServerNetwork(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_dedicated_server_network",
		Description: "Get details of a specific network attached to a dedicated server. Returns CIDR, family, interface type, distribution method, and status",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getNetworkArgs) (*mcp.CallToolResult, any, error) {
		network, err := h.client.Hosts.GetDedicatedServerNetwork(ctx, args.ServerID, args.NetworkID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(network)
	})
}

// --- get_dedicated_server_network_usage ---

type serverIDArgs struct {
	ServerID string `json:"server_id" jsonschema:"dedicated server ID,required"`
}

func registerGetDedicatedServerNetworkUsage(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_dedicated_server_network_usage",
		Description: "Get network traffic utilization for a dedicated server. Returns committed and current utilization per network type. Useful for checking how much bandwidth is being consumed",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args serverIDArgs) (*mcp.CallToolResult, any, error) {
		usage, err := h.client.Hosts.GetDedicatedServerNetworkUsage(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(usage)
	})
}

// --- add_dedicated_server_public_ipv4_network ---

type addNetworkArgs struct {
	ServerID           string `json:"server_id" jsonschema:"dedicated server ID,required"`
	DistributionMethod string `json:"distribution_method" jsonschema:"network distribution method: 'gateway' for a fully routable subnet with a gateway or 'route' for /32 interface aliases (route availability depends on location),required"`
	Mask               int    `json:"mask" jsonschema:"subnet prefix length (e.g. 29 for /29, 28 for /28). For route (alias) networks always use 32. For gateway networks: available sizes depend on location and remaining quota (max 72 total IPs per family across all gateway networks),required"`
}

func registerAddDedicatedServerPublicIPv4Network(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "add_dedicated_server_public_ipv4_network",
		Description: `Add an additional public IPv4 network to a dedicated server.

Distribution methods:
- "gateway": a fully routable subnet with a gateway (e.g. /29, /28, /27, /26). Max 2 additional gateway networks per server. Total IPs across all public IPv4 gateway networks (including the default one) must not exceed 72
- "route": /32 addresses as interface aliases. Max 32 route IPs per server. Not available in all locations

Use list_dedicated_server_networks to check existing networks and remaining quota before ordering`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addNetworkArgs) (*mcp.CallToolResult, any, error) {
		network, err := h.client.Hosts.AddDedicatedServerPublicIPv4Network(ctx, args.ServerID, serverscom.NetworkInput{
			DistributionMethod: args.DistributionMethod,
			Mask:               args.Mask,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(network)
	})
}

// --- add_dedicated_server_private_ipv4_network ---

func registerAddDedicatedServerPrivateIPv4Network(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "add_dedicated_server_private_ipv4_network",
		Description: `Add an additional private IPv4 network to a dedicated server.

Distribution methods:
- "gateway": a fully routable private subnet with a gateway. Max 2 additional gateway networks per server. Total IPs across all private IPv4 gateway networks (including the default one) must not exceed 72
- "route": /32 addresses as interface aliases. Max 32 route IPs per server. Not available in all locations

Use list_dedicated_server_networks to check existing networks and remaining quota before ordering`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addNetworkArgs) (*mcp.CallToolResult, any, error) {
		network, err := h.client.Hosts.AddDedicatedServerPrivateIPv4Network(ctx, args.ServerID, serverscom.NetworkInput{
			DistributionMethod: args.DistributionMethod,
			Mask:               args.Mask,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(network)
	})
}

// --- activate_dedicated_server_public_ipv6_network ---

func registerActivateDedicatedServerPublicIPv6Network(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "activate_dedicated_server_public_ipv6_network",
		Description: `Activate a public IPv6 network for a dedicated server.

Only one IPv6 allocation per server is possible. The provisioning model depends on the location:
- Some locations assign a /64 directly routed to the server
- Other locations assign a /125 (point-to-point link) with a /64 routed over it

Both variants are activated through this single call and appear as separate network entries (the /125 and the /64) in list_dedicated_server_networks. IPv6 quotas are independent from IPv4 quotas`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args serverIDArgs) (*mcp.CallToolResult, any, error) {
		network, err := h.client.Hosts.ActivateDedicatedServerPubliIPv6Network(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(network)
	})
}

// --- delete_dedicated_server_network ---

func registerDeleteDedicatedServerNetwork(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_dedicated_server_network",
		Description: "Remove an additional network from a dedicated server. Only additional networks (additional=true) can be deleted; the default network cannot be removed. The operation frees up quota for new networks",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getNetworkArgs) (*mcp.CallToolResult, any, error) {
		network, err := h.client.Hosts.DeleteDedicatedServerNetwork(ctx, args.ServerID, args.NetworkID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(network)
	})
}
