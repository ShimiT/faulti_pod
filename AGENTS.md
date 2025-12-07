# AGENTS.md

## Project Overview
This repository is a Go service (module github.com/your-org/faulty-app) packaged as a Docker image and deployable to Kubernetes via a Helm chart. Primary tech stack: Go 1.22, Docker, Helm (charts under helm/), and Kubernetes deployment automation driven from the Makefile.

## Dev Environment Tips
- Required tools
  - Go 1.22 (module-aware)
  - Docker (for image builds)
  - Helm 3 (for deploys under helm/)
  - kubectl (to interact with clusters)
- Common commands
  - Install dependencies: go mod download
  - Build the Go binaries (local): go build ./...
  - Run the program locally (if main packages are in the repo root or packages): go build -o bin/faulty-app ./... && ./bin/faulty-app
  - Run tests locally: go test ./... -v
- Using the Makefile (recommended for image/chart flows)
  - Build image locally: make build
  - Push image to registry: make push
  - Deploy to cluster using Helm: make deploy-cluster
  - Full e2e flow (build → push → helm deploy): make deploy-e2e
- Environment variables used by Makefile
  - IMAGE (default: shimit/faulti_pod)
  - TAG (default: latest)
  - RELEASE (default: faulty)
  - NAMESPACE (default: default)
  - CHART_DIR (default: helm)

## Build Commands
```bash
# From Makefile
make build          # docker build -t $(IMAGE):$(TAG) .
make push           # docker push $(IMAGE
