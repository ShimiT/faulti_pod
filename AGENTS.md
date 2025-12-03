# AGENTS.md

This file is a concise guide for AI agents (and humans) working on the faulty-app repository (module: `github.com/your-org/faulty-app`). It focuses on concrete commands and project-specific conventions extracted from the repository (Makefile, go.mod, Helm chart directory).

---

## 1) Dev Environment Tips

Minimum required runtime:
- Go 1.22 (from go.mod)

Key project defaults (from Makefile):
- IMAGE (default): `shimit/faulti_pod`
- TAG (default): `latest`
- RELEASE (default): `faulty`
- NAMESPACE (default): `default`
- CHART_DIR: `helm`

Essential setup commands:
- Fetch Go dependencies:
  - go mod download
- Build Go packages / binary (local):
  - go build ./...
- Build Docker image (Makefile):
  - make build
  - This runs: `docker build -t $(IMAGE):$(TAG) .`
  - Override image/tag as environment variables, e.g.:
    - IMAGE=myrepo/myimage TAG=ci-123 make build
- Push Docker image to registry:
  - make push
  - This runs: `docker push $(IMAGE):$(TAG)`
- Deploy the Helm chart to a Kubernetes cluster:
  - make deploy-cluster
  - Under the hood:
    - `helm upgrade --install $(RELEASE) $(CHART_DIR) --namespace $(NAMESPACE) --create-namespace --set image.repository=$(IMAGE) --set image.tag=$(TAG)`
  - Override variables inline:
    - IMAGE=myrepo/myimage TAG=v1.2.3 NAMESPACE=staging RELEASE=myrelease make deploy-cluster
- Full e2e prepare + deploy target:
  - make deploy-e2e
  - This runs `build`, then `push`, then `deploy-cluster` in sequence.

Notes / environment prerequisites:
- Docker is required for `make build`/`make push`.
- Helm and kubectl are required for `make deploy-cluster` to succeed.
- The repository expects a Dockerfile in the project root (Docker build context is `.`). Ensure it exists/works before `make build`.

---

## 2) Testing Instructions

Unit tests and typical Go test workflow:
- Run all unit tests:
  - go test ./... -v
- Run a specific package:
  - go test ./pkg/somepkg -v
- Run a specific test function:
  - go test ./... -run TestMyThing -v
- Collect coverage:
  - go test ./... -coverprofile=coverage.out
  - go tool cover -html=coverage.out -o coverage.html

E2E / integration:
- The Makefile target `deploy-e2e` automates build → push → helm deploy. Use this when preparing integration/e2e runs that require the image to be in
