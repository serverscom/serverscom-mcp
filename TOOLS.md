# Tool Reference

## Order Options

Tools for browsing available configurations before ordering servers or storage.

| Tool | Description |
|---|---|
| `list_locations` | List all data center locations with status and supported capabilities |
| `get_location` | Get details of a specific location |
| `list_server_model_options` | List available server models in a location: CPU, RAM, drive slots |
| `get_server_model_option` | Get full details of a server model including drive slot configuration |
| `list_ram_options` | List RAM upgrade options for a server model in a location |
| `list_os_options` | List available operating systems for a server model in a location |
| `get_os_option` | Get OS details including supported filesystems |
| `list_drive_model_options` | List available drive models: capacity, interface, media type |
| `get_drive_model_option` | Get details of a specific drive model |
| `list_uplink_options` | List network uplink models: speed and redundancy |
| `get_uplink_option` | Get details of a specific uplink model |
| `list_bandwidth_options` | List bandwidth plans for an uplink model |
| `get_bandwidth_option` | Get details of a specific bandwidth plan |
| `list_sbm_flavor_options` | List Scalable Bare Metal (SBM) flavors in a location |
| `get_sbm_flavor_option` | Get details of a specific SBM flavor |
| `list_sbm_os_options` | List available operating systems for an SBM flavor |
| `list_rbs_flavors` | List Remote Block Storage flavors in a location: IOPS/GB and bandwidth/GB per tier |

## Hosts

`list_hosts` returns all host types in the account — dedicated servers, Kubernetes baremetal nodes, and SBM servers. Dedicated servers are a subtype of host.

| Tool | Description |
|---|---|
| `list_hosts` | List all hosts with optional filtering by search pattern, location, type, and labels |

## SSH Keys

Account-level SSH key management. To attach or detach keys from a specific server, see [Dedicated Servers → SSH Keys](#ssh-keys-1).

| Tool | Description |
|---|---|
| `list_ssh_keys` | List all SSH keys in the account |
| `get_ssh_key` | Get key details by fingerprint |
| `create_ssh_key` | Add a new SSH public key to the account |
| `update_ssh_key` | Update key name or labels |
| `delete_ssh_key` | Delete an SSH key from the account |

## Dedicated Servers

### Servers

| Tool | Description |
|---|---|
| `get_dedicated_server` | Get server details: configuration, `operational_status`, `power_status` |
| `update_dedicated_server` | Update title, labels, `user_data` (cloud-init), or `ipxe_config` |
| `list_dedicated_server_features` | List features and their current status for a server |

### Features

Each feature supports `activate_*` and `deactivate_*` operations. All feature changes are **asynchronous** — poll `list_dedicated_server_features` until the status reaches `activated` or `deactivated`.

| Feature | Description |
|---|---|
| `disaggregated_public_ports` | Disaggregated public ports |
| `disaggregated_private_ports` | Disaggregated private ports |
| `no_public_ip_address` | No default public IP; restricts rescue mode, OOB access, and additional public networks |
| `no_private_ip` | No default private IP |
| `no_public_network` | Public interface not configured at all — more restrictive than `no_public_ip_address` |
| `host_rescue_mode` | Rescue mode; requires `auth_methods` and optionally `ssh_key_fingerprints` |
| `oob_public_access` | Out-of-band public access |
| `private_ipxe_boot` | Private iPXE boot; requires `ipxe_config` |

### SSH Keys

| Tool | Description |
|---|---|
| `list_dedicated_server_ssh_keys` | List SSH keys currently attached to a server |
| `attach_ssh_keys_to_dedicated_server` | Attach one or more SSH keys to a server |
| `detach_ssh_key_from_dedicated_server` | Detach an SSH key from a server |

### Power Management

All power operations are **asynchronous** — poll `get_dedicated_server` and check `power_status`.

| Tool | Description |
|---|---|
| `power_on_dedicated_server` | Power on |
| `power_off_dedicated_server` | Power off |
| `power_cycle_dedicated_server` | Hard reboot |

### OS Reinstallation

| Tool | Description |
|---|---|
| `reinstall_dedicated_server` | Reinstall OS with custom partition layout, RAID config, and SSH key injection |

**Async** — `operational_status` transitions to `installation`, then back to `normal`. Supports custom RAID levels (0, 1, 5, 6, 10, 50, 60) and partition layouts. Max 1 SSH key per reinstall.

### Networks

| Tool | Description |
|---|---|
| `list_dedicated_server_networks` | List all networks (public/private, IPv4/IPv6) attached to a server |
| `get_dedicated_server_network` | Get network details: CIDR, family, interface type, distribution method, status |
| `get_dedicated_server_network_usage` | Check current and committed bandwidth utilization |
| `add_dedicated_server_public_ipv4_network` | Add additional public IPv4 network (`gateway` or `route`) |
| `add_dedicated_server_private_ipv4_network` | Add additional private IPv4 network (`gateway` or `route`) |
| `activate_dedicated_server_public_ipv6_network` | Activate IPv6 (one allocation per server) |
| `delete_dedicated_server_network` | Remove an additional network (default network cannot be deleted) |

**Quotas:**
- Route (alias) networks: max 32 IPs per server
- Additional gateway networks per family (public/private IPv4): max 2
- Total IPs across all gateway networks per family: max 72
- IPv6: max 1 allocation per server (single /64 or /125 + /64 depending on location)

### Storage

| Tool | Description |
|---|---|
| `list_dedicated_server_drive_slots` | List all drive slots and installed drives for a server |

## Remote Block Storage

RBS provides iSCSI network-attached block storage mountable to Dedicated Servers and Kubernetes nodes without physical hardware changes. Built on Ceph.

| Tool | Description |
|---|---|
| `list_rbs_volumes` | List all RBS volumes in the account |
| `get_rbs_volume` | Get volume details: size, IOPS, bandwidth, target IQN, Volume IP |
| `create_rbs_volume` | Create a new volume (requires location and flavor) |
| `update_rbs_volume` | Update name, size (increase only), or labels |
| `delete_rbs_volume` | Delete a volume — must be disconnected first; irreversible |
| `get_rbs_volume_credentials` | Get iSCSI CHAP credentials: username, password, target IQN, Volume IP |
| `reset_rbs_volume_credentials` | Rotate iSCSI password — invalidates active connections |

**Account limits:** max 99 volumes · max 1 TB per volume · max 10 TB total

**Connection:** iSCSI on TCP port 3260 with CHAP authentication. After `create_rbs_volume`, call `get_rbs_volume_credentials` to retrieve connection details. Use `list_rbs_flavors` (see [Order Options](#order-options)) to pick a performance tier.

## L2 Segments

L2 segments unite dedicated servers within a location group into a single broadcast domain, enabling direct Layer 2 communication via MAC addresses without routing. Available only for Dedicated Servers in Enterprise locations.

| Tool | Description |
|---|---|
| `list_l2_segments` | List all L2 segments in the account |
| `get_l2_segment` | Get segment details: type, status, location group |
| `create_l2_segment` | Create a new L2 segment (public or private, native or trunk mode) |
| `update_l2_segment` | Update name, members, or labels |
| `delete_l2_segment` | Delete a segment and release all associated networks |
| `list_l2_segment_members` | List member servers with mode and VLAN number |
| `list_l2_segment_networks` | List networks (subnets) assigned to a segment |
| `change_l2_segment_networks` | Add or remove networks on a segment |
| `list_l2_location_groups` | List available location groups (required for segment creation) |

**Member modes:** `native` (no OS config needed, 1 per interface per type) · `trunk` (requires VLAN sub-interface on server OS, up to 16 per type)

**Limits per server:** 1 native public + 1 native private + 16 public trunk + 16 private trunk = 34 max

## Kubernetes Clusters

Managed Kubernetes clusters on Servers.com infrastructure. Clusters are provisioned through the Servers.com portal; these tools provide read and label-management access.

| Tool | Description |
|---|---|
| `list_kubernetes_clusters` | List all Kubernetes clusters in the account |
| `get_kubernetes_cluster` | Get cluster details: status, location, labels |
| `update_kubernetes_cluster` | Update cluster labels (replaces all existing labels) |
| `list_kubernetes_cluster_nodes` | List all nodes in a cluster: role, type, status, IP addresses |
| `get_kubernetes_cluster_node` | Get full details of a specific cluster node |

**Node roles:** `master` · `node`

**Node types:** `cloud` · `baremetal`
