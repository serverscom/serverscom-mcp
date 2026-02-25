package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

func registerRBSTools(server *mcp.Server, h *handler) {
	registerListRBSVolumes(server, h)
	registerGetRBSVolume(server, h)
	registerCreateRBSVolume(server, h)
	registerUpdateRBSVolume(server, h)
	registerDeleteRBSVolume(server, h)
	registerGetRBSVolumeCredentials(server, h)
	registerResetRBSVolumeCredentials(server, h)
	registerListRBSFlavors(server, h)
}

// --- list_rbs_volumes ---

func registerListRBSVolumes(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_rbs_volumes",
		Description: "List all Remote Block Storage (RBS) volumes in the account. RBS provides network-attached block storage (iSCSI) mountable to Dedicated Servers and Kubernetes nodes. Returns volume ID, name, size (GB), status, location, flavor, IOPS, bandwidth, and labels",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		volumes, err := h.client.RemoteBlockStorageVolumes.Collection().Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(volumes)
	})
}

// --- get_rbs_volume ---

type rbsVolumeIDArgs struct {
	VolumeID string `json:"volume_id" jsonschema:"RBS volume ID,required"`
}

func registerGetRBSVolume(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_rbs_volume",
		Description: "Get details of a specific Remote Block Storage volume. Returns name, size (GB), status, location, flavor, IOPS, bandwidth, iSCSI target IQN, and Volume IP address (used as the iSCSI portal address when connecting from a server)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args rbsVolumeIDArgs) (*mcp.CallToolResult, any, error) {
		volume, err := h.client.RemoteBlockStorageVolumes.Get(ctx, args.VolumeID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(volume)
	})
}

// --- create_rbs_volume ---

type createRBSVolumeArgs struct {
	Name       string            `json:"name" jsonschema:"volume name,required"`
	Size       int64             `json:"size" jsonschema:"volume size in GB,required"`
	LocationID int               `json:"location_id" jsonschema:"location ID from list_locations — must support Remote Block Storage,required"`
	FlavorID   int               `json:"flavor_id" jsonschema:"flavor ID from list_rbs_flavors — defines IOPS and bandwidth characteristics,required"`
	Labels     map[string]string `json:"labels,omitempty" jsonschema:"optional key-value labels"`
}

func registerCreateRBSVolume(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "create_rbs_volume",
		Description: `Create a new Remote Block Storage (RBS) volume. RBS is a scalable block storage solution (built on Ceph) that attaches to Dedicated Servers and Kubernetes nodes via iSCSI — no physical disk changes required.

Key concepts:
- Each volume has a Volume IP address (iSCSI portal), a target IQN, and CHAP credentials for authentication
- IOPS and bandwidth limits are determined by the flavor and scale with volume size (optimized for 4KB block sizes)
- Billing: charged per GB/month based on flavor; size increases are prorated; cancellation charges the full month

Account limits: max 99 volumes, max 1 TB per volume, max 10 TB total across all volumes

Workflow after creation:
1. Use get_rbs_volume_credentials to get the iSCSI username, password, target IQN, and Volume IP
2. On Linux: iscsiadm discovery → configure CHAP → login → format → mount (port 3260)
3. On Windows: configure iSCSI Initiator with Volume IP:3260, enable CHAP, initialize disk
4. On Kubernetes: create a CHAP Secret, reference the volume in a pod spec as an iscsi volume type

Use list_rbs_flavors to see available flavors per location. Use list_locations to find a location that supports RBS.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createRBSVolumeArgs) (*mcp.CallToolResult, any, error) {
		input := serverscom.RemoteBlockStorageVolumeCreateInput{
			Name:       args.Name,
			Size:       args.Size,
			LocationID: args.LocationID,
			FlavorID:   args.FlavorID,
			Labels:     args.Labels,
		}
		volume, err := h.client.RemoteBlockStorageVolumes.Create(ctx, input)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(volume)
	})
}

// --- update_rbs_volume ---

type updateRBSVolumeArgs struct {
	VolumeID string            `json:"volume_id" jsonschema:"RBS volume ID,required"`
	Name     string            `json:"name,omitempty" jsonschema:"new volume name"`
	Size     int64             `json:"size,omitempty" jsonschema:"new size in GB — volumes can only be expanded, not shrunk"`
	Labels   map[string]string `json:"labels,omitempty" jsonschema:"labels to set (replaces existing labels)"`
}

func registerUpdateRBSVolume(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_rbs_volume",
		Description: "Update a Remote Block Storage volume's name, size, or labels. Size can only be increased — volumes cannot be shrunk. Size changes are prorated within the current billing month. All fields are optional — only provided fields are changed",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateRBSVolumeArgs) (*mcp.CallToolResult, any, error) {
		input := serverscom.RemoteBlockStorageVolumeUpdateInput{
			Name:   args.Name,
			Size:   args.Size,
			Labels: args.Labels,
		}
		volume, err := h.client.RemoteBlockStorageVolumes.Update(ctx, args.VolumeID, input)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(volume)
	})
}

// --- delete_rbs_volume ---

func registerDeleteRBSVolume(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_rbs_volume",
		Description: "Delete a Remote Block Storage volume. The volume must be unmounted and the iSCSI session logged out on all connected servers before deletion. Cancellation charges the full current month regardless of timing. This action is irreversible — all data will be lost",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args rbsVolumeIDArgs) (*mcp.CallToolResult, any, error) {
		err := h.client.RemoteBlockStorageVolumes.Delete(ctx, args.VolumeID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(map[string]string{"status": "deleted", "volume_id": args.VolumeID})
	})
}

// --- get_rbs_volume_credentials ---

func registerGetRBSVolumeCredentials(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_rbs_volume_credentials",
		Description: "Get iSCSI CHAP credentials for a Remote Block Storage volume. Returns username, password, target IQN, and Volume IP (portal address on port 3260). These are used to authenticate and connect: on Linux via iscsiadm, on Windows via iSCSI Initiator, on Kubernetes via a CHAP Secret manifest",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args rbsVolumeIDArgs) (*mcp.CallToolResult, any, error) {
		creds, err := h.client.RemoteBlockStorageVolumes.GetCredentials(ctx, args.VolumeID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(creds)
	})
}

// --- reset_rbs_volume_credentials ---

func registerResetRBSVolumeCredentials(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "reset_rbs_volume_credentials",
		Description: "Reset iSCSI CHAP credentials for a Remote Block Storage volume. Generates a new password, immediately invalidating the previous one. Any servers or Kubernetes nodes currently connected will lose access — disconnect gracefully (unmount → iscsiadm logout) before resetting, then reconnect with the new credentials",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args rbsVolumeIDArgs) (*mcp.CallToolResult, any, error) {
		volume, err := h.client.RemoteBlockStorageVolumes.ResetCredentials(ctx, args.VolumeID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(volume)
	})
}

// --- list_rbs_flavors ---

type rbsFlavorsArgs struct {
	LocationID int64 `json:"location_id" jsonschema:"location ID from list_locations,required"`
}

func registerListRBSFlavors(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_rbs_flavors",
		Description: "List available Remote Block Storage flavors for a location. A flavor defines the performance tier: IOPS per GB and bandwidth per GB — actual volume IOPS and bandwidth scale with the size you provision. Returns flavor ID, name, IOPS/GB, bandwidth/GB, and minimum volume size. Use flavor_id when creating an RBS volume",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args rbsFlavorsArgs) (*mcp.CallToolResult, any, error) {
		flavors, err := h.client.Locations.RemoteBlockStorageFlavors(args.LocationID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(flavors)
	})
}
