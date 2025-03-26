# kubectl-meshsync-snapshot

A kubectl plugin for capturing MeshSync snapshots for Meshery.

## Overview

`kubectl-meshsync-snapshot` is a kubectl plugin that helps manage MeshSync snapshots for Meshery. It allows deploying MeshSync temporarily, capturing cluster state, importing snapshots to Meshery, and cleaning up resources.

## Installation

### Using Krew

```bash
kubectl krew install meshsync-snapshot
```

### Manual Installation

1. Download the appropriate binary for your platform from [Releases](https://github.com/Prajwal-kp-18/kubectl-meshsync-snapshot/releases)
2. Make it executable: `chmod +x kubectl-meshsync_snapshot`
3. Move it to a directory in your PATH: `mv kubectl-meshsync_snapshot /usr/local/bin/`

## Getting Started

### Prerequisites

- A running Kubernetes cluster
- kubectl installed and configured
- (Optional) Meshery server for importing snapshots

### Quick Start

Deploy MeshSync, capture a snapshot of all namespaces, and clean up:

```bash
# 1. Deploy MeshSync temporarily
kubectl meshsync-snapshot deploy

# 2. Capture a snapshot of all namespaces
kubectl meshsync-snapshot capture -A -o my-cluster-snapshot.yaml

# 3. Clean up MeshSync resources
kubectl meshsync-snapshot cleanup
```

### Import to Meshery

If you have a Meshery server running, you can import your snapshot:

```bash
kubectl meshsync-snapshot import -i my-cluster-snapshot.yaml -u https://meshery.example.com -t your-token
```

## Usage

The plugin provides the following commands:

### Deploy MeshSync

Deploy MeshSync temporarily to capture cluster state:

```bash
kubectl meshsync-snapshot deploy [flags]
```

Flags:
- `--namespace`, `-n`: Namespace to deploy MeshSync (default: "meshery")
- `--version`, `-v`: MeshSync version to deploy (default: "latest")
- `--timeout`, `-t`: Timeout for deployment (default: 2m0s)

### Capture Snapshot

Capture cluster state using MeshSync:

```bash
kubectl meshsync-snapshot capture [flags]
```

Flags:
- `--namespace`, `-n`: Namespace where MeshSync is deployed (default: "meshery")
- `--output`, `-o`: Output file for snapshot (default: "meshsync-snapshot.yaml")
- `--format`, `-f`: Output format (yaml or json) (default: "yaml")
- `--timeout`, `-t`: Timeout for capture operation (default: 1m0s)
- `--all-namespaces`, `-A`: Capture resources from all namespaces

### Import Snapshot

Import snapshot to Meshery:

```bash
kubectl meshsync-snapshot import [flags]
```

Flags:
- `--url`, `-u`: Meshery server URL (default: "http://localhost:9081")
- `--token`, `-t`: Meshery authentication token
- `--input`, `-i`: Input snapshot file path (default: "meshsync-snapshot.yaml")
- `--timeout`: Timeout for import operation (default: 30s)

### Cleanup

Remove MeshSync resources:

```bash
kubectl meshsync-snapshot cleanup [flags]
```

Flags:
- `--namespace`, `-n`: Namespace where MeshSync is deployed (default: "meshery")
- `--timeout`, `-t`: Timeout for cleanup operation (default: 1m0s)
- `--force`, `-f`: Force cleanup even if resources are still in use

## Examples

### Capture cluster state in a single namespace

```bash
# Deploy MeshSync
kubectl meshsync-snapshot deploy -n monitoring

# Capture state
kubectl meshsync-snapshot capture -n monitoring -o monitoring-snapshot.yaml

# Clean up when done
kubectl meshsync-snapshot cleanup -n monitoring
```

### Capture entire cluster and import to Meshery

```bash
# Deploy MeshSync
kubectl meshsync-snapshot deploy

# Capture all namespaces
kubectl meshsync-snapshot capture -A -o cluster-snapshot.yaml

# Import to Meshery
kubectl meshsync-snapshot import -i cluster-snapshot.yaml -u https://meshery.example.com -t your-token

# Clean up
kubectl meshsync-snapshot cleanup
```

## Development

### Prerequisites

- Go 1.22+
- Access to a Kubernetes cluster
- kubectl installed

### Building from source

```bash
# Clone the repository
git clone https://github.com/Prajwal-kp-18/kubectl-meshsync-snapshot.git
cd kubectl-meshsync-snapshot

# Build
go build -o kubectl-meshsync_snapshot cmd/kubectl-meshsync_snapshot/main.go

# Run locally
./kubectl-meshsync_snapshot help
```

### Running tests

```bash
# Run all tests
make test

# Format code
make fmt
```

### Creating releases

For maintainers:

1. Update version in Makefile
2. Tag the repo: `git tag -a v0.1.0 -m "Release v0.1.0"`
3. Push the tag: `git push origin v0.1.0`
4. The GitHub Actions workflow will create a release with binaries

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

- Prajwal-kp-18 