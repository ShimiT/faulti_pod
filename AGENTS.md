# AGENTS.md

This file is written for AI coding agents working on the `faulti_pod` (module `github.com/your-org/faulty-app`) repository. It contains project-specific commands, conventions and pointers to make productive changes and PRs.

---

## Dev Environment Tips

- Required runtime:
  - Go 1.22 (see go.mod). Use a Go 1.22 toolchain.
- Clone:
  - git clone https://github.com/ShimiT/faulti_pod.git
  - cd faulti_pod
- Ensure module dependencies:
  - go mod download
  - When updating dependencies: go get ./... and then go mod tidy
- Build locally:
  - Using Makefile: `make build`
    - This runs: `docker build -t $(IMAGE):$(TAG) .`
    - Default variables (can be overridden on the command line):
      - IMAGE (default `shimit/faulti_pod`)
      - TAG (default `latest`)
      - Example overriding: `make build IMAGE=shimit/faulti_pod TAG=dev`
  - Alternatively build the Go binary locally (module-aware):
    - `go build ./...` or `go build -o bin/faulty-app ./cmd/...` (adjust path to your entry package if present)
- Docker/image:
  - Push to registry: `make push` (runs `docker push $(IMAGE):$(TAG)`)
  - You must be logged in to your Docker registry before `make push` (e.g. `docker login`).
- Helm and cluster deploy:
  - The Makefile has `deploy-cluster` which runs:
    - `helm upgrade --install $(RELEASE) $(CHART_DIR) --namespace $(NAMESPACE) --create-namespace --set image.repository=$(IMAGE) --set image.tag=$(TAG)`
  - Defaults:
    - RELEASE = `faulty`
    - CHART_DIR = `helm`
    - NAMESPACE = `default`
  - Example full flow (build -> push -> deploy):
    - `make deploy-e2e` (this runs `build`, `push`, then `deploy-cluster`)
  - Override namespace or release on the same command:
    - `make deploy-cluster NAMESPACE=staging RELEASE=faulty-staging IMAGE=myrepo/faulti_pod TAG=feature-branch`
- Files to verify before Docker build:
  - Confirm a `Dockerfile` exists at repo root (the Makefile's docker build uses `.`). If it does not, adapt the Makefile or add the Dockerfile.
  - Ensure the `helm/` directory exists and contains the chart (Makefile expects `helm/`).

---

## Testing Instructions

- Run all tests:
  - `go test ./...`
  - Verbose: `go test ./... -v`
- Run a specific package:
  - `go test ./pkg/name -v` (replace `./pkg/name` with package path)
- Run a single test function:
  - `go test ./... -run TestName -v`
- Race detector:
  - `go test -race ./...`
- Coverage:
  - `go test ./... -coverprofile=coverage.out`
  - View HTML: `go tool cover -html=coverage.out -o coverage.html`
- Common checks to run locally before a PR:
  - `go test ./...`
  - `go vet ./...`
  - `gofmt -l .` and `gofmt -w .` to format
  - `goimports` to fix imports (if used)
  - `go mod tidy` to keep go.mod clean
- CI normally expected to run the same commands above; ensure your branch passes them locally first.

---

## PR Instructions

- Title format:
  - Use Conventional Commit style for clarity. Examples:
    - feat(service): add new retry logic for pod controller
    - fix(deploy): ensure correct image tag passed to helm chart
    - docs: update README for helm values
  - Short, present-tense summary, optional `(scope)`.
- Description:
  - Provide a short summary, motivation, and list of changes.
  - Include commands to reproduce or test changes locally (e.g., `make build`, `go test ./...`, `make deploy-cluster
