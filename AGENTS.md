# AGENTS.md

## Project Overview
This repository is a Go service (module path github.com/your-org/faulty-app) built with Go 1.22. It is containerized (Docker) and packaged with Helm (helm/ chart directory). Primary tech stack: Go, Docker, Helm, and Kubernetes deployment via helm.

## Dev Environment Tips
- Required runtime/tooling:
  - Go 1.22 (as specified in go.mod)
  - Docker (for image build/push)
  - Helm (for deploy-cluster target)
  - kubectl / cluster access if you will deploy the chart to a cluster
- Install dependencies:
  - go mod download
  - (optional) go mod tidy
- Common local dev commands:
  - Build all packages: go build ./...
  - Run main packages locally (runs any main package found): go run ./...
  - Run a single package (from repo root): go run ./path/to/package
- Working directory: run commands from the repository root (where go.mod and Makefile live).
- Environment variables used by Makefile (can be overridden):
  - IMAGE (default shimit/faulti_pod)
  - TAG (default latest)
  - RELEASE (default faulty)
  - NAMESPACE (default default)
  - CHART_DIR (default helm)

## Build Commands
```bash
make build        # builds a Docker image using: docker build -t $(IMAGE):$(TAG) .
make push         # pushes the Docker image to the registry: docker push $(IMAGE):$(TAG)
make deploy-cluster
                  # runs helm upgrade --install $(RELEASE) $(CHART_DIR) --namespace $(NAMESPACE) --create-namespace --set image.repository=$(IMAGE) --set image.tag=$(TAG)
make deploy-e2e   # runs build, push, then deploy-cluster (combines build -> push -> deploy-cluster)
```

## Testing Instructions
- CI configuration: no CI workflow files detected in repository root (no .github/workflows//* found).
- Full test suite (Go's testing package): 
  - go test ./... -v
- Run a single package or test:
  - go test ./path/to/package -v
  - go test ./path/to/package -run TestName -v
- Testing framework: Go standard testing package (testing)
- IMPORTANT: Add or update tests for code you change, even if not asked.

## Code Style
- Formatting:
  - Use gofmt to format code. Apply formatting to the repo with:
    - gofmt -w .
  - Ensure module builds after formatting: go build ./...
- Basic static checking:
  - Use go vet for basic correctness checks:
    - go vet ./...
- Key conventions (follow idiomatic Go):
  - Keep imports grouped and ordered as standard (std lib first, then external).
  - Use short, clear names for variables; exported identifiers must be CamelCase.
  - Follow package-oriented layout (look for package main under cmd/ or main.go at repo root).
- Run linter/formatters before committing: run gofmt and go vet locally.

## PR
