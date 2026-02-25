package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

func registerL2SegmentTools(server *mcp.Server, h *handler) {
	registerListL2Segments(server, h)
	registerGetL2Segment(server, h)
	registerCreateL2Segment(server, h)
	registerUpdateL2Segment(server, h)
	registerDeleteL2Segment(server, h)
	registerListL2SegmentMembers(server, h)
	registerListL2SegmentNetworks(server, h)
	registerChangeL2SegmentNetworks(server, h)
	registerListL2LocationGroups(server, h)
}

// --- list_l2_segments ---

func registerListL2Segments(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_l2_segments",
		Description: "List all L2 segments in the account. An L2 segment unites dedicated servers within a location group into a single broadcast domain so they communicate directly via MAC addresses without routing. Returns segment ID, name, type (public/private), status, location group, and labels",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		segments, err := h.client.L2Segments.Collection().Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(segments)
	})
}

// --- get_l2_segment ---

type l2SegmentIDArgs struct {
	SegmentID string `json:"segment_id" jsonschema:"L2 segment ID,required"`
}

func registerGetL2Segment(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_l2_segment",
		Description: "Get details of a specific L2 segment by ID. Returns name, type (public/private), status, location group code, and labels",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args l2SegmentIDArgs) (*mcp.CallToolResult, any, error) {
		segment, err := h.client.L2Segments.Get(ctx, args.SegmentID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(segment)
	})
}

// --- create_l2_segment ---

type l2SegmentMemberInput struct {
	ID   string `json:"id" jsonschema:"dedicated server host ID,required"`
	Mode string `json:"mode" jsonschema:"port mode: 'native' or 'trunk' — see tool description,required"`
}

type createL2SegmentArgs struct {
	Name            string                 `json:"name,omitempty" jsonschema:"human-readable segment name"`
	Type            string                 `json:"type" jsonschema:"'public' (bridges public interfaces, traffic counted against server package) or 'private' (bridges private interfaces, traffic not counted),required"`
	LocationGroupID int64                  `json:"location_group_id" jsonschema:"location group ID from list_l2_location_groups — must match segment type (public group for public segment, private group for private),required"`
	Members         []l2SegmentMemberInput `json:"members" jsonschema:"at least one dedicated server to add as initial member,required"`
	Labels          map[string]string      `json:"labels,omitempty" jsonschema:"optional key-value labels"`
}

func registerCreateL2Segment(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "create_l2_segment",
		Description: `Create a new L2 segment that unites dedicated servers into a single broadcast domain, enabling direct Layer 2 communication via MAC addresses without routing. Available only for Dedicated Servers in Enterprise locations.

Segment types:
- "public": bridges servers on their public network interfaces. Traffic between servers within a public segment counts against their traffic packages. Cannot be used if the server has public network failover activated
- "private": bridges servers on their private network interfaces. Traffic within a private segment is not counted against packages

Member modes:
- "native": no VLAN configuration needed on the server OS — servers communicate at Layer 2 without any interface changes. Each server supports at most 1 native public + 1 native private segment
- "trunk": the server OS must configure a VLAN sub-interface with the assigned VLAN number (visible in list_l2_segment_members). Each server supports up to 16 public trunk + 16 private trunk segments simultaneously. When added in trunk mode, the server's alias IPs become part of the L2 segment. Billed per trunk VLAN

Constraints:
- At least 1 host required
- Server must not be in rescue mode, undergoing OS reinstallation, or in maintenance mode
- ACL rules do not apply to traffic within an L2 segment
- Avoid IP ranges 10.0.0.0/8 and 192.168.0.0/16 for custom addressing — may overlap with Servers.com private ranges; prefer 172.16.0.0/12
- Maximum 34 L2 segments per server total (1 native public + 1 native private + 16 public trunk + 16 private trunk)

Typical use cases:
- Floating IP / Failover IP: an IP address migrates between servers automatically on failure, using Keepalived (Linux) or CARP (FreeBSD). Requires an additional network or alias IP added to the segment via change_l2_segment_networks
- Gateway redundancy with VRRP: entire subnets can float between servers, providing default gateway redundancy
- Network isolation / subnetting: isolate projects or departments into separate broadcast domains without routing between them
- Custom IP addressing: use any IP range for DHCP/PXE servers or legacy software — servers communicate without a gateway
- Multicast: enables Ethernet multicast (one-to-many) for streaming, gaming, or financial data feeds

Use list_l2_location_groups to find a valid location_group_id.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createL2SegmentArgs) (*mcp.CallToolResult, any, error) {
		members := make([]serverscom.L2SegmentMemberInput, len(args.Members))
		for i, m := range args.Members {
			members[i] = serverscom.L2SegmentMemberInput{ID: m.ID, Mode: m.Mode}
		}

		input := serverscom.L2SegmentCreateInput{
			Type:            args.Type,
			LocationGroupID: args.LocationGroupID,
			Members:         members,
		}
		if args.Name != "" {
			input.Name = &args.Name
		}
		if args.Labels != nil {
			input.Labels = args.Labels
		}

		segment, err := h.client.L2Segments.Create(ctx, input)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(segment)
	})
}

// --- update_l2_segment ---

type updateL2SegmentArgs struct {
	SegmentID string                 `json:"segment_id" jsonschema:"L2 segment ID,required"`
	Name      string                 `json:"name,omitempty" jsonschema:"new segment name"`
	Members   []l2SegmentMemberInput `json:"members,omitempty" jsonschema:"full updated member list — replaces the existing set entirely; omit to leave unchanged"`
	Labels    map[string]string      `json:"labels,omitempty" jsonschema:"labels to set (replaces existing labels)"`
}

func registerUpdateL2Segment(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "update_l2_segment",
		Description: `Update an L2 segment's name, members, or labels.

When members are provided, the list replaces the existing member set entirely — include all servers that should remain in the segment, not just the ones being added or removed. Use list_l2_segment_members first to retrieve current members.

Mode rules and constraints are the same as in create_l2_segment. Servers being added must not be in rescue mode, OS reinstallation, or maintenance. When a server is removed from a trunk segment, its alias IPs are also removed from the segment.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateL2SegmentArgs) (*mcp.CallToolResult, any, error) {
		input := serverscom.L2SegmentUpdateInput{}
		if args.Name != "" {
			input.Name = &args.Name
		}
		if args.Labels != nil {
			input.Labels = args.Labels
		}
		if args.Members != nil {
			members := make([]serverscom.L2SegmentMemberInput, len(args.Members))
			for i, m := range args.Members {
				members[i] = serverscom.L2SegmentMemberInput{ID: m.ID, Mode: m.Mode}
			}
			input.Members = members
		}

		segment, err := h.client.L2Segments.Update(ctx, args.SegmentID, input)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(segment)
	})
}

// --- delete_l2_segment ---

func registerDeleteL2Segment(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_l2_segment",
		Description: "Delete an L2 segment. All member servers are removed from the broadcast domain, all networks and alias IPs assigned to the segment are released, and billing stops",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args l2SegmentIDArgs) (*mcp.CallToolResult, any, error) {
		err := h.client.L2Segments.Delete(ctx, args.SegmentID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(map[string]string{"status": "deleted", "segment_id": args.SegmentID})
	})
}

// --- list_l2_segment_members ---

func registerListL2SegmentMembers(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_l2_segment_members",
		Description: "List all member servers in an L2 segment. Returns host ID, title, mode (native/trunk), VLAN number configured for this server in the segment, and provisioning status",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args l2SegmentIDArgs) (*mcp.CallToolResult, any, error) {
		members, err := h.client.L2Segments.Members(args.SegmentID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(members)
	})
}

// --- list_l2_segment_networks ---

func registerListL2SegmentNetworks(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_l2_segment_networks",
		Description: `List all networks (subnets) assigned to an L2 segment. Returns CIDR, address family, distribution method, and status.

Networks within an L2 segment are billed per additional network. They enable routable addressing across the broadcast domain and are required for Failover IP / Floating IP setups (e.g. with Keepalived/VRRP): an IP can migrate between member servers without routing changes. Use change_l2_segment_networks to add or remove networks.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args l2SegmentIDArgs) (*mcp.CallToolResult, any, error) {
		networks, err := h.client.L2Segments.Networks(args.SegmentID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(networks)
	})
}

// --- change_l2_segment_networks ---

type l2SegmentCreateNetworkInput struct {
	Mask               int    `json:"mask" jsonschema:"subnet prefix length (e.g. 29 for /29, 28 for /28),required"`
	DistributionMethod string `json:"distribution_method" jsonschema:"'gateway' for a routable subnet with a gateway, or 'route' for /32 alias IPs,required"`
}

type changeL2SegmentNetworksArgs struct {
	SegmentID string                        `json:"segment_id" jsonschema:"L2 segment ID,required"`
	Create    []l2SegmentCreateNetworkInput `json:"create,omitempty" jsonschema:"networks to add"`
	Delete    []string                      `json:"delete,omitempty" jsonschema:"IDs of existing networks to remove — use list_l2_segment_networks to get IDs"`
}

func registerChangeL2SegmentNetworks(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "change_l2_segment_networks",
		Description: `Add or remove networks on an L2 segment in a single atomic operation.

Networks enable IP addressing within the broadcast domain and are required for Failover IP (Floating IP) scenarios — where a shared IP migrates between servers using Keepalived or CARP without routing changes.

Billing: each additional network is billed separately. Alias IPs added via 'route' distribution are billed daily (postpaid).

Use list_l2_segment_networks to retrieve existing network IDs before deleting.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args changeL2SegmentNetworksArgs) (*mcp.CallToolResult, any, error) {
		input := serverscom.L2SegmentChangeNetworksInput{
			Delete: args.Delete,
		}
		for _, n := range args.Create {
			input.Create = append(input.Create, serverscom.L2SegmentCreateNetworksInput{
				Mask:               n.Mask,
				DistributionMethod: n.DistributionMethod,
			})
		}

		segment, err := h.client.L2Segments.ChangeNetworks(ctx, args.SegmentID, input)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(segment)
	})
}

// --- list_l2_location_groups ---

func registerListL2LocationGroups(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_l2_location_groups",
		Description: `List available location groups for L2 segments. A location group defines which data centers can share an L2 segment.

Location group types:
- "public": for public L2 segments — spans a single location
- "private": for private L2 segments — may span multiple locations, allowing servers in different data centers to share a broadcast domain

Returns group ID (use as location_group_id in create_l2_segment), name, code, type, and list of included location IDs. Only Enterprise locations are eligible for L2 segments.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		groups, err := h.client.L2Segments.LocationGroups().Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(groups)
	})
}
