package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

func registerKubernetesClusterTools(server *mcp.Server, h *handler) {
	registerListKubernetesClusters(server, h)
	registerGetKubernetesCluster(server, h)
	registerUpdateKubernetesCluster(server, h)
	registerListKubernetesClusterNodes(server, h)
	registerGetKubernetesClusterNode(server, h)
}

type kubernetesClusterIDArgs struct {
	ClusterID string `json:"cluster_id" jsonschema:"cluster ID,required"`
}

type updateKubernetesClusterArgs struct {
	ClusterID string            `json:"cluster_id" jsonschema:"cluster ID,required"`
	Labels    map[string]string `json:"labels,omitempty" jsonschema:"key-value labels (replaces all existing labels)"`
}

type listKubernetesClusterNodesArgs struct {
	ClusterID string `json:"cluster_id" jsonschema:"cluster ID,required"`
}

type kubernetesClusterNodeArgs struct {
	ClusterID string `json:"cluster_id" jsonschema:"cluster ID,required"`
	NodeID    string `json:"node_id" jsonschema:"node ID,required"`
}

func registerListKubernetesClusters(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_kubernetes_clusters",
		Description: "List all Kubernetes clusters in the account. Returns cluster ID, name, status, location_id, labels, and timestamps for each cluster.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		clusters, err := h.client.KubernetesClusters.Collection().Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(clusters)
	})
}

func registerGetKubernetesCluster(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_kubernetes_cluster",
		Description: "Get details of a Kubernetes cluster by ID. Returns name, status, location_id, labels, and timestamps.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args kubernetesClusterIDArgs) (*mcp.CallToolResult, any, error) {
		cluster, err := h.client.KubernetesClusters.Get(ctx, args.ClusterID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(cluster)
	})
}

func registerUpdateKubernetesCluster(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_kubernetes_cluster",
		Description: "Update labels on a Kubernetes cluster. The provided labels replace all existing labels on the cluster.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateKubernetesClusterArgs) (*mcp.CallToolResult, any, error) {
		input := serverscom.KubernetesClusterUpdateInput{
			Labels: args.Labels,
		}
		cluster, err := h.client.KubernetesClusters.Update(ctx, args.ClusterID, input)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(cluster)
	})
}

func registerListKubernetesClusterNodes(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_kubernetes_cluster_nodes",
		Description: "List all nodes in a Kubernetes cluster. Each node includes: role (node/master), type (cloud/baremetal), status, private and public IPv4 addresses, hostname, configuration, and ref_id.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listKubernetesClusterNodesArgs) (*mcp.CallToolResult, any, error) {
		nodes, err := h.client.KubernetesClusters.Nodes(args.ClusterID).Collect(ctx)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(nodes)
	})
}

func registerGetKubernetesClusterNode(server *mcp.Server, h *handler) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_kubernetes_cluster_node",
		Description: "Get details of a specific node in a Kubernetes cluster. Returns role, type, status, IP addresses, hostname, configuration, ref_id, labels, and timestamps. Use list_kubernetes_cluster_nodes to find node IDs.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args kubernetesClusterNodeArgs) (*mcp.CallToolResult, any, error) {
		node, err := h.client.KubernetesClusters.GetNode(ctx, args.ClusterID, args.NodeID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return jsonResult(node)
	})
}
