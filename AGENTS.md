# AGENTS.md

This repository contains a Go application (module: `github.com/your-org/faulty-app`) and Helm charts for deploying it. Use the guidance below when authoring code, tests, PRs, or automations for this project.

---

## 1) Dev Environment Tips

Requirements
- Go toolchain: go 1.22 (see go.mod)
- Docker (for local image build/push)
- Helm (for local cluster deploys)

Quick setup
- Fetch dependencies:
  ```bash
  go mod download
  go mod tidy
  ```
- Build the binary locally:
  ```bash
  go build ./...
  ```
- Run the application locally (if a main package exists):
  ```bash
  # run from repo root
  go run ./...          # or go run ./cmd/yourcmd if using cmd/ layout
  ```

Docker / image workflows (Makefile targets)
- Build the image:
  ```bash
  # uses IMAGE and TAG; defaults are in Makefile
  make build
  # or override
  make build IMAGE=shimit/faulti_pod TAG=1.2.3
  ```
  Under the hood:
  ```bash
  docker build -t $(IMAGE):$(TAG) .
  ```
- Push the image:
  ```bash
  make push IMAGE=shimit/faulti_pod TAG=1.2.3
  # underlying: docker push $(IMAGE):$(TAG)
  ```
- Deploy to a Kubernetes cluster with Helm:
  ```bash
  # default variables in Makefile: IMAGE, TAG, RELEASE=faulty, NAMESPACE=default, CHART_DIR=helm
  make deploy-cluster IMAGE=shimit/faulti_pod TAG=1.2.3 RELEASE=myrelease NAMESPACE=staging
  ```
  Under the hood the Makefile runs:
  ```bash
  helm upgrade --install $(RELEASE) $(CHART_DIR) \
    --namespace $(NAMESPACE) --create-namespace \
    --set image.repository=$(IMAGE) --set image.tag=$(TAG)
  ```
- Full E2E deploy (build -> push -> helm deploy):
  ```bash
  make deploy-e2e IMAGE=shimit/faulti_pod TAG=ci-123
  ```

Environment variable tips
- Override Makefile variables on the make command line (see examples above).
- Ensure your local Docker daemon is logged into the registry you plan to push to.

---

## 2) Testing Instructions

Framework: standard Go testing (go test).

Run all tests
```
