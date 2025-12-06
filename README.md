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
brew tap suminhong/homebrew-tap
brew install mk
```

### From Source
```bash
git clone https://github.com/suminhong/multi-kubectl.git
cd multi-kubectl
make build
sudo make install
```

## Release Process (for Maintainers)
1. Tag a new version: `git tag v0.0.1 && git push origin v0.0.1`
2. Get the SHA256 of the tarball:
   ```bash
   curl -L https://github.com/suminhong/multi-kubectl/archive/v0.0.1.tar.gz | shasum -a 256
   ```
3. Update `mk.rb` in `suminhong/homebrew-tap` with the new URL and SHA256.
