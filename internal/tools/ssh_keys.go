package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

func registerSSHKeyTools(server *mcp.Server, h *handler) {
	registerListSSHKeys(server, h)
	registerGetSSHKey(server, h)
	registerCreateSSHKey(server, h)
	registerUpdateSSHKey(server, h)
	registerDeleteSSHKey(server, h)
	registerListDedicatedServerSSHKeys(server, h)
	registerAttachSSHKeysToDedicatedServer(server, h)
	registerDetachSSHKeyFromDedicatedServer(server, h)
}

// --- list_ssh_keys ---

func registerListSSHKeys(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_ssh_keys",
		Description: "List all SSH keys in the account. Returns name, fingerprint, labels, and timestamps for each key",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		keys, err := h.client.SSHKeys.Collection().Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(keys)
	})
}

// --- get_ssh_key ---

type sshKeyFingerprintArgs struct {
	Fingerprint string `json:"fingerprint" jsonschema:"SSH key fingerprint,required"`
}

func registerGetSSHKey(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_ssh_key",
		Description: "Get details of a specific SSH key by its fingerprint",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args sshKeyFingerprintArgs) (*mcp.CallToolResult, any, error) {
		key, err := h.client.SSHKeys.Get(ctx, args.Fingerprint)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(key)
	})
}

// --- create_ssh_key ---

type createSSHKeyArgs struct {
	Name      string            `json:"name" jsonschema:"human-readable key name,required"`
	PublicKey string            `json:"public_key" jsonschema:"public key content in OpenSSH format,required"`
	Labels    map[string]string `json:"labels,omitempty" jsonschema:"optional labels to attach to the key"`
}

func registerCreateSSHKey(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_ssh_key",
		Description: "Add a new SSH public key to the account. Returns the created key with its fingerprint, which is used to reference the key in all other operations",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createSSHKeyArgs) (*mcp.CallToolResult, any, error) {
		key, err := h.client.SSHKeys.Create(ctx, serverscom.SSHKeyCreateInput{
			Name:      args.Name,
			PublicKey: args.PublicKey,
			Labels:    args.Labels,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(key)
	})
}

// --- update_ssh_key ---

type updateSSHKeyArgs struct {
	Fingerprint string            `json:"fingerprint" jsonschema:"SSH key fingerprint,required"`
	Name        string            `json:"name,omitempty" jsonschema:"new key name"`
	Labels      map[string]string `json:"labels,omitempty" jsonschema:"labels to set on the key (replaces existing labels)"`
}

func registerUpdateSSHKey(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_ssh_key",
		Description: "Update an SSH key's name and/or labels",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateSSHKeyArgs) (*mcp.CallToolResult, any, error) {
		key, err := h.client.SSHKeys.Update(ctx, args.Fingerprint, serverscom.SSHKeyUpdateInput{
			Name:   args.Name,
			Labels: args.Labels,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(key)
	})
}

// --- delete_ssh_key ---

func registerDeleteSSHKey(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_ssh_key",
		Description: "Delete an SSH key from the account by its fingerprint. The key must be detached from all dedicated servers before deletion",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args sshKeyFingerprintArgs) (*mcp.CallToolResult, any, error) {
		if err := h.client.SSHKeys.Delete(ctx, args.Fingerprint); err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(map[string]string{"status": "deleted", "fingerprint": args.Fingerprint})
	})
}

// --- list_dedicated_server_ssh_keys ---

type dedicatedServerSSHKeyListArgs struct {
	ServerID string `json:"server_id" jsonschema:"dedicated server ID,required"`
}

func registerListDedicatedServerSSHKeys(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_dedicated_server_ssh_keys",
		Description: "List SSH keys currently attached to a dedicated server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args dedicatedServerSSHKeyListArgs) (*mcp.CallToolResult, any, error) {
		keys, err := h.client.Hosts.ListDedicatedServerSSHKeys(ctx, args.ServerID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(keys)
	})
}

// --- attach_ssh_keys_to_dedicated_server ---

type attachSSHKeysArgs struct {
	ServerID     string   `json:"server_id" jsonschema:"dedicated server ID,required"`
	Fingerprints []string `json:"ssh_key_fingerprints" jsonschema:"list of SSH key fingerprints to attach,required"`
}

func registerAttachSSHKeysToDedicatedServer(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "attach_ssh_keys_to_dedicated_server",
		Description: "Attach one or more SSH keys to a dedicated server. Use list_ssh_keys to find fingerprints. Keys will be available on the server after the next OS reinstall",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args attachSSHKeysArgs) (*mcp.CallToolResult, any, error) {
		keys, err := h.client.Hosts.AttachSSHKeysToDedicatedServer(ctx, args.ServerID, serverscom.SSHKeyAttachInput{
			SSHKeyFingerprints: args.Fingerprints,
		})
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(keys)
	})
}

// --- detach_ssh_key_from_dedicated_server ---

type detachSSHKeyArgs struct {
	ServerID    string `json:"server_id" jsonschema:"dedicated server ID,required"`
	Fingerprint string `json:"fingerprint" jsonschema:"fingerprint of the SSH key to detach,required"`
}

func registerDetachSSHKeyFromDedicatedServer(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "detach_ssh_key_from_dedicated_server",
		Description: "Detach a single SSH key from a dedicated server by its fingerprint",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args detachSSHKeyArgs) (*mcp.CallToolResult, any, error) {
		if err := h.client.Hosts.DetachSSHKeyFromDedicatedServer(ctx, args.ServerID, args.Fingerprint); err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(map[string]string{"status": "detached", "fingerprint": args.Fingerprint})
	})
}
