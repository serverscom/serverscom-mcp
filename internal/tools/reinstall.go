package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

func registerReinstallTools(server *mcp.Server, h *handler) {
	registerReinstallDedicatedServer(server, h)
}

// reinstallPartitionInput mirrors OperatingSystemReinstallPartitionInput.
//
// target — mount point: "/", "/boot", "swap", etc.
// size   — size in MB (ignored when fill=true, but still required by API).
// fs     — filesystem type; use the values from OperatingSystemOption.Filesystems
//
//	(returned by list_os_options / get_os_option) — that list is the authoritative
//	source for the chosen OS. Omit for swap partitions.
//
// fill   — expand the partition to fill all remaining space (only one partition per layout
//
//	group may have fill=true).
type reinstallPartitionInput struct {
	Target string  `json:"target"`
	Size   int     `json:"size"`
	Fs     *string `json:"fs,omitempty"`
	Fill   bool    `json:"fill,omitempty"`
}

// reinstallLayoutInput defines how a set of drive slots should be configured.
//
// slot_positions — list of drive slot positions from list_dedicated_server_drive_slots
//
//	(zero-based). All positions in one entry are grouped together and
//	will form a single RAID volume (or a JBOD if raid is omitted/null).
//
// raid           — RAID level: 0, 1, 5, 6, 10, 50 or 60.
//
//	Omit (null) to leave the drives unformatted / as individual disks.
//
// ignore         — when true the slot positions are skipped entirely during
//
//	partitioning (no RAID, no partitions). Use this for drives that
//	should remain untouched.
//
// partitions     — ordered list of partitions to create on this volume.
//
//	Must be omitted (empty) when ignore=true.
type reinstallLayoutInput struct {
	SlotPositions []int                     `json:"slot_positions"`
	Raid          *int                      `json:"raid,omitempty"`
	Ignore        *bool                     `json:"ignore,omitempty"`
	Partitions    []reinstallPartitionInput `json:"partitions,omitempty"`
}

// reinstallDedicatedServerArgs is the full input for reinstall_dedicated_server.
//
// Workflow to build a correct request:
//  1. Call list_dedicated_server_drive_slots to find which slot positions exist
//     and which drive models are installed in them.
//  2. Call list_os_options (with the server's location_id and server_model_id from
//     get_dedicated_server → configuration_details) to pick an operating_system_id.
//  3. Fill in drives.layout — at minimum one entry covering the boot drive(s).
//
// The reinstall is asynchronous: after the call the server's operational_status
// changes to "installation". Poll get_dedicated_server until it returns "normal".
type reinstallDedicatedServerArgs struct {
	// ServerID is the target dedicated server.
	ServerID string `json:"server_id" jsonschema:"dedicated server ID,required"`

	// Hostname to set inside the OS (max 253 chars; max 15 chars for Windows).
	Hostname string `json:"hostname" jsonschema:"server hostname to configure inside the OS,required"`

	// OperatingSystemID from list_os_options. If omitted the server is
	// reinstalled with the same OS it currently runs.
	OperatingSystemID *int64 `json:"operating_system_id,omitempty" jsonschema:"OS ID from list_os_options; omit to keep the current OS"`

	// Drives describes the disk layout. Required: at least one layout entry.
	Drives reinstallDrivesInput `json:"drives" jsonschema:"drive layout configuration,required"`

	// SSHKeyFingerprints — up to one SSH key fingerprint to inject into the new OS.
	SSHKeyFingerprints []string `json:"ssh_key_fingerprints,omitempty" jsonschema:"SSH key fingerprints to inject (max 1 item)"`
}

// reinstallDrivesInput wraps the layout slice for JSON compatibility.
type reinstallDrivesInput struct {
	// Layout — one entry per RAID group (or individual drive).
	// Every slot position that should participate in the OS installation
	// must appear in exactly one layout entry.
	Layout []reinstallLayoutInput `json:"layout" jsonschema:"ordered list of drive layout groups,required"`
}

func registerReinstallDedicatedServer(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "reinstall_dedicated_server",
		Description: `Reinstall the operating system on a dedicated server. This is an ASYNC operation.

Workflow:
1. Call list_dedicated_server_drive_slots to learn slot positions and installed drives.
2. Call list_os_options (location_id + server_model_id from get_dedicated_server → configuration_details) to get a valid operating_system_id.
3. Build the drives.layout array — each entry groups slot_positions into one RAID/disk unit and lists partitions to create on it.
4. Submit this tool. The server's operational_status will become "installation" immediately.
5. Poll get_dedicated_server until operational_status returns to "normal".

drives.layout rules:
- slot_positions: zero-based positions matching list_dedicated_server_drive_slots output.
- raid: RAID level (0,1,5,6,10,50,60) or omit for a single-disk volume.
- ignore: set true to skip a slot entirely (no partitioning, no RAID).
- partitions: list of {target, size (MB), fs, fill}. Exactly one partition per layout group may have fill=true to consume remaining space.
- fs: use the values from the "filesystems" field returned by list_os_options / get_os_option for the chosen OS — that list is the authoritative source of supported filesystems. Omit fs for swap partitions.

Example layout for two-disk RAID-1 with /boot + swap + /:
  [{"slot_positions":[0,1],"raid":1,"partitions":[{"target":"/boot","size":1024,"fs":"ext4"},{"target":"swap","size":4096},{"target":"/","size":100000,"fs":"ext4","fill":true}]}]`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args reinstallDedicatedServerArgs) (*mcp.CallToolResult, any, error) {
		layout := make([]serverscom.OperatingSystemReinstallLayoutInput, 0, len(args.Drives.Layout))
		for _, l := range args.Drives.Layout {
			partitions := make([]serverscom.OperatingSystemReinstallPartitionInput, 0, len(l.Partitions))
			for _, p := range l.Partitions {
				partitions = append(partitions, serverscom.OperatingSystemReinstallPartitionInput{
					Target: p.Target,
					Size:   p.Size,
					Fs:     p.Fs,
					Fill:   p.Fill,
				})
			}
			layout = append(layout, serverscom.OperatingSystemReinstallLayoutInput{
				SlotPositions: l.SlotPositions,
				Raid:          l.Raid,
				Ignore:        l.Ignore,
				Partitions:    partitions,
			})
		}

		input := serverscom.OperatingSystemReinstallInput{
			Hostname: args.Hostname,
			Drives: serverscom.OperatingSystemReinstallDrivesInput{
				Layout: layout,
			},
			OperatingSystemID:  args.OperatingSystemID,
			SSHKeyFingerprints: args.SSHKeyFingerprints,
		}

		ds, err := h.client.Hosts.ReinstallOperatingSystemForDedicatedServer(ctx, args.ServerID, input)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(ds)
	})
}
