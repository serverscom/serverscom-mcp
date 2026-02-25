package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/serverscom/serverscom-mcp/internal/tools"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "serverscom-mcp",
		Usage: "MCP server for Servers.com infrastructure management",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "Servers.com API token",
				Sources:  cli.EnvVars("SC_TOKEN"),
				Required: true,
			},
			&cli.StringFlag{
				Name:    "endpoint",
				Aliases: []string{"e"},
				Usage:   "Servers.com API endpoint",
				Sources: cli.EnvVars("SC_ENDPOINT"),
				Value:   "https://api.servers.com/v1",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			server := mcp.NewServer(&mcp.Implementation{
				Name:    "serverscom-mcp",
				Version: Version,
			}, nil)

			tools.Register(server, cmd.String("token"), cmd.String("endpoint"), Version)

			if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
				return fmt.Errorf("server error: %w", err)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
