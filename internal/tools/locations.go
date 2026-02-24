package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerLocationTools(server *mcp.Server, h *handler) {
	registerListLocations(server, h)
	registerGetLocation(server, h)
	registerListServerModelOptions(server, h)
	registerGetServerModelOption(server, h)
	registerListRAMOptions(server, h)
	registerListOSOptions(server, h)
	registerGetOSOption(server, h)
	registerListDriveModelOptions(server, h)
	registerGetDriveModelOption(server, h)
	registerListUplinkOptions(server, h)
	registerGetUplinkOption(server, h)
	registerListBandwidthOptions(server, h)
	registerGetBandwidthOption(server, h)
	registerListSBMFlavorOptions(server, h)
	registerGetSBMFlavorOption(server, h)
	registerListSBMOSOptions(server, h)
}

// --- list_locations ---

func registerListLocations(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_locations",
		Description: "List all available locations. Returns location ID, name, status, code, supported features, and enabled capabilities (L2 segments, private racks, load balancers)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		locations, err := h.client.Locations.Collection().Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(locations)
	})
}

// --- get_location ---

type getLocationArgs struct {
	LocationID int64 `json:"location_id" jsonschema:"location ID,required"`
}

func registerGetLocation(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_location",
		Description: "Get details of a specific location by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getLocationArgs) (*mcp.CallToolResult, any, error) {
		location, err := h.client.Locations.GetLocation(ctx, args.LocationID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(location)
	})
}

// --- list_server_model_options ---

type locationArgs struct {
	LocationID int64 `json:"location_id" jsonschema:"location ID,required"`
}

func registerListServerModelOptions(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_server_model_options",
		Description: "List available server models for ordering in a location. Returns model ID, name, CPU specs, RAM, drive slot count, and RAID controller info. Use location ID from list_locations",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args locationArgs) (*mcp.CallToolResult, any, error) {
		models, err := h.client.Locations.ServerModelOptions(args.LocationID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(models)
	})
}

// --- get_server_model_option ---

type serverModelArgs struct {
	LocationID    int64 `json:"location_id" jsonschema:"location ID,required"`
	ServerModelID int64 `json:"server_model_id" jsonschema:"server model ID,required"`
}

func registerGetServerModelOption(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_server_model_option",
		Description: "Get detailed info on a server model in a location, including full drive slot configuration",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args serverModelArgs) (*mcp.CallToolResult, any, error) {
		model, err := h.client.Locations.GetServerModelOption(ctx, args.LocationID, args.ServerModelID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(model)
	})
}

// --- list_ram_options ---

func registerListRAMOptions(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_ram_options",
		Description: "List available RAM upgrade options for a server model in a location. Returns RAM size (GB) and memory type",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args serverModelArgs) (*mcp.CallToolResult, any, error) {
		options, err := h.client.Locations.RAMOptions(args.LocationID, args.ServerModelID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(options)
	})
}

// --- list_os_options ---

func registerListOSOptions(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_os_options",
		Description: "List available operating systems for a server model in a location. Returns OS ID, full name, version, architecture, and supported filesystems",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args serverModelArgs) (*mcp.CallToolResult, any, error) {
		options, err := h.client.Locations.OperatingSystemOptions(args.LocationID, args.ServerModelID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(options)
	})
}

// --- get_os_option ---

type osOptionArgs struct {
	LocationID      int64 `json:"location_id" jsonschema:"location ID,required"`
	ServerModelID   int64 `json:"server_model_id" jsonschema:"server model ID,required"`
	OperatingSystemID int64 `json:"operating_system_id" jsonschema:"operating system ID,required"`
}

func registerGetOSOption(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_os_option",
		Description: "Get details of a specific operating system option for a server model in a location",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args osOptionArgs) (*mcp.CallToolResult, any, error) {
		option, err := h.client.Locations.GetOperatingSystemOption(ctx, args.LocationID, args.ServerModelID, args.OperatingSystemID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(option)
	})
}

// --- list_drive_model_options ---

func registerListDriveModelOptions(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_drive_model_options",
		Description: "List available drive models for a server model in a location. Returns drive ID, name, capacity (GB), interface, form factor, and media type (HDD/SSD/NVMe)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args serverModelArgs) (*mcp.CallToolResult, any, error) {
		drives, err := h.client.Locations.DriveModelOptions(args.LocationID, args.ServerModelID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(drives)
	})
}

// --- get_drive_model_option ---

type driveModelArgs struct {
	LocationID    int64 `json:"location_id" jsonschema:"location ID,required"`
	ServerModelID int64 `json:"server_model_id" jsonschema:"server model ID,required"`
	DriveModelID  int64 `json:"drive_model_id" jsonschema:"drive model ID,required"`
}

func registerGetDriveModelOption(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_drive_model_option",
		Description: "Get details of a specific drive model option for a server model in a location",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args driveModelArgs) (*mcp.CallToolResult, any, error) {
		drive, err := h.client.Locations.GetDriveModelOption(ctx, args.LocationID, args.ServerModelID, args.DriveModelID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(drive)
	})
}

// --- list_uplink_options ---

func registerListUplinkOptions(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_uplink_options",
		Description: "List available uplink models for a server model in a location. Returns uplink ID, name, type, speed (Mbps), and redundancy flag",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args serverModelArgs) (*mcp.CallToolResult, any, error) {
		uplinks, err := h.client.Locations.UplinkOptions(args.LocationID, args.ServerModelID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(uplinks)
	})
}

// --- get_uplink_option ---

type uplinkOptionArgs struct {
	LocationID    int64 `json:"location_id" jsonschema:"location ID,required"`
	ServerModelID int64 `json:"server_model_id" jsonschema:"server model ID,required"`
	UplinkModelID int64 `json:"uplink_model_id" jsonschema:"uplink model ID,required"`
}

func registerGetUplinkOption(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_uplink_option",
		Description: "Get details of a specific uplink model option for a server model in a location",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args uplinkOptionArgs) (*mcp.CallToolResult, any, error) {
		uplink, err := h.client.Locations.GetUplinkOption(ctx, args.LocationID, args.ServerModelID, args.UplinkModelID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(uplink)
	})
}

// --- list_bandwidth_options ---

type bandwidthListArgs struct {
	LocationID    int64 `json:"location_id" jsonschema:"location ID,required"`
	ServerModelID int64 `json:"server_model_id" jsonschema:"server model ID,required"`
	UplinkModelID int64 `json:"uplink_model_id" jsonschema:"uplink model ID,required"`
}

func registerListBandwidthOptions(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_bandwidth_options",
		Description: "List available bandwidth plans for a specific uplink model in a location. Returns bandwidth ID, name, type, and optional commit (Mbps). Use list_uplink_options to get uplink model IDs",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args bandwidthListArgs) (*mcp.CallToolResult, any, error) {
		options, err := h.client.Locations.BandwidthOptions(args.LocationID, args.ServerModelID, args.UplinkModelID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(options)
	})
}

// --- get_bandwidth_option ---

type bandwidthGetArgs struct {
	LocationID    int64 `json:"location_id" jsonschema:"location ID,required"`
	ServerModelID int64 `json:"server_model_id" jsonschema:"server model ID,required"`
	UplinkModelID int64 `json:"uplink_model_id" jsonschema:"uplink model ID,required"`
	BandwidthID   int64 `json:"bandwidth_id" jsonschema:"bandwidth option ID,required"`
}

func registerGetBandwidthOption(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_bandwidth_option",
		Description: "Get details of a specific bandwidth option for an uplink model in a location",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args bandwidthGetArgs) (*mcp.CallToolResult, any, error) {
		option, err := h.client.Locations.GetBandwidthOption(ctx, args.LocationID, args.ServerModelID, args.UplinkModelID, args.BandwidthID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(option)
	})
}

// --- list_sbm_flavor_options ---

func registerListSBMFlavorOptions(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_sbm_flavor_options",
		Description: "List available Scalable Bare Metal (SBM) flavors in a location. Returns flavor ID, name, CPU specs, RAM, drives configuration, and uplink/bandwidth details",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args locationArgs) (*mcp.CallToolResult, any, error) {
		flavors, err := h.client.Locations.SBMFlavorOptions(args.LocationID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(flavors)
	})
}

// --- get_sbm_flavor_option ---

type sbmFlavorArgs struct {
	LocationID      int64 `json:"location_id" jsonschema:"location ID,required"`
	SBMFlavorModelID int64 `json:"sbm_flavor_model_id" jsonschema:"SBM flavor model ID,required"`
}

func registerGetSBMFlavorOption(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_sbm_flavor_option",
		Description: "Get details of a specific SBM flavor in a location",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args sbmFlavorArgs) (*mcp.CallToolResult, any, error) {
		flavor, err := h.client.Locations.GetSBMFlavorOption(ctx, args.LocationID, args.SBMFlavorModelID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(flavor)
	})
}

// --- list_sbm_os_options ---

func registerListSBMOSOptions(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_sbm_os_options",
		Description: fmt.Sprintf("List available operating systems for a Scalable Bare Metal (SBM) flavor in a location. Use %q and %q to get the required IDs", "list_locations", "list_sbm_flavor_options"),
	}, func(ctx context.Context, req *mcp.CallToolRequest, args sbmFlavorArgs) (*mcp.CallToolResult, any, error) {
		options, err := h.client.Locations.SBMOperatingSystemOptions(args.LocationID, args.SBMFlavorModelID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(options)
	})
}
