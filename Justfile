# List available recipes
default:
    @just --list

# Build binary for current platform (version = "dev")
build:
    go build -o serverscom-mcp .

# Build binary with explicit version
build-version version:
    go build -ldflags "-s -w -X main.Version={{version}}" -o serverscom-mcp .

# Run tests
test:
    go test ./...

# Check goreleaser config
check:
    goreleaser check

# Snapshot build for all platforms (no publish)
snapshot:
    goreleaser build --snapshot --clean

# Build DXT locally (run `just snapshot` first)
dxt:
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always --dirty)
    rm -rf dxt && mkdir -p dxt/server
    cp server/run.sh dxt/server/
    cp dist/mcp_linux_amd64_v1/serverscom-mcp   dxt/server/serverscom-mcp-linux-amd64
    cp dist/mcp_linux_arm64_v8.0/serverscom-mcp  dxt/server/serverscom-mcp-linux-arm64
    cp dist/mcp_darwin_amd64_v1/serverscom-mcp   dxt/server/serverscom-mcp-darwin-amd64
    cp dist/mcp_darwin_arm64_v8.0/serverscom-mcp dxt/server/serverscom-mcp-darwin-arm64
    cp dist/mcp_windows_amd64_v1/serverscom-mcp.exe dxt/server/serverscom-mcp-windows-amd64.exe
    chmod +x dxt/server/serverscom-mcp-*
    cp icon.svg dxt/
    jq --arg v "$VERSION" '.version = $v' manifest.json > dxt/manifest.json
    cd dxt && zip -r ../serverscom-mcp.dxt .
    echo "Built serverscom-mcp.dxt ($VERSION)"
