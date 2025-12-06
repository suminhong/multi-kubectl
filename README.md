# Multi-Kubectl (mk)

`mk` is a CLI tool that wraps `kubectl` to execute commands across multiple Kubernetes contexts simultaneously.

## Usage

```bash
mk [kubectl-args] [--context <filter>]
```

### Examples

Run `get pods` on all contexts:
```bash
mk get pods -n devops
```

Run on contexts containing "prod":
```bash
mk get pods --context prod
```

Run on contexts containing "prod" OR "stage":
```bash
mk get pods --context prod,stage
```

## Installation

### Homebrew
You can install `mk` using your custom tap:

```bash
brew install suminhong/tap/mk
```

### From Source
```bash
git clone https://github.com/suminhong/multi-kubectl.git
cd multi-kubectl
make build
sudo make install
```

## Release Process

Releases are automated using GitHub Actions and Goreleaser. Simply push to the `main` branch to trigger a new release.
