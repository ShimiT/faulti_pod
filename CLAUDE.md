# CLAUDE.md

## Project overview
Faulty-app is a small Go service packaged as a Docker image and deployed with Helm. The repository is a single Go module (go 1.22) and contains a Helm chart at `helm/` plus a Makefile that builds, pushes, and deploys the image.

## Common commands
(From the repo Makefile and Go toolchain)

- Build container image:
  - `make build` (uses IMAGE/TAG env vars: defaults to `shimit/faulti_pod:latest`)
  - Equivalent explicit: `docker build -t shimit/faulti_pod:latest .`
- Push image:
  - `make push`
- Deploy to cluster (Helm):
  - `make deploy-cluster`
  - Customizable: `make deploy-cluster IMAGE=your/repo TAG=v1 NAMESPACE=prod RELEASE=your-release`
- Full e2e deploy (build → push → helm upgrade/install):
  - `make deploy-e2e`
- Build/run locally with Go:
  - `go build ./...`
  - Run (if a main package exists): `go run ./...`
- Tests:
  - `go test ./... -v`
  - Coverage: `go test ./... -coverprofile=coverage.out`
- Lint / static checks (recommended):
  - `gofmt -s -w .` (format)
  - `go vet ./...`
  - (Optional) `golangci-lint run` if you add it to CI.

## Code style guidelines
There is no project-specific linter config in the repo; follow standard Go conventions:
- Target Go 1.22 (as declared in `go.mod`).
- Formatting: run `gofmt -s -w .` before committing.
- Use `go vet` to catch suspicious constructs.
- Prefer idiomatic Go naming and error handling (return errors, wrap with fmt.Errorf/%w).
- Keep package APIs small and well-documented with comments for exported symbols.
- If you add a linter, prefer `golangci-lint` with common enabled linters (govet, staticcheck, gosimple, gofmt/gci, errcheck).

## Testing instructions
- Unit tests: standard Go testing framework.
  - Run all tests: `go test ./... -v`
  - Generate coverage: `go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out`
- Add tests next to the packages they exercise (`*_test.go` files).
- CI should run `gofmt -l`/`gofmt -s -w`, `go vet ./...`, `go test ./... -race` where applicable.

## Repository etiquette
- Branch naming:
  - feature branches: `feature/<short-desc>` or `feat/<short-desc>`
  - bugfix branches: `fix/<short-desc>`
  - chores/refactors: `chore/<short-desc>` or `refactor/<short-desc>`
- Commit messages:
  - Follow a Conventional Commits-like pattern: `<type>(scope?): short summary`
    - Examples: `feat(api): add /health endpoint`, `fix(build): correct Dockerfile path`
  - Include a longer description in the body when needed and reference issue IDs: `Refs: #123`
- PRs:
  - Small, focused PRs per logical change.
  - Include:
    - What changed and why.
    - How to test locally (commands).
    - Link to issue if present.
  - Assign reviewers and wait for at least one approval + passing tests.
  - Update changelog or release notes when the PR affects behavior or public APIs.
- CI expectations before merge:
  - All tests pass.
  - Code formatted (`gofmt`) and vetted (`go vet`).
  - Docker image builds successfully if the change affects Docker/Helm.

## Key patterns & conventions in this codebase
- Single Go module declared in `go.mod` (module path: `github.com/your-org/faulty-app`).
- Container-first deployment:
  - Docker image built from repository root; image defaults controlled by Makefile `IMAGE` and `TAG`.
  - Helm chart in `helm/` accepts `image.repository` and `image.tag` values; `make deploy-cluster` wires these values.
- Makefile orchestration:
  - `build`, `push`, `deploy-cluster`, `deploy-e2e` are the canonical developer flows.
- Keep release name and namespace configurable:
  - Makefile variables: `RELEASE`, `NAMESPACE`, `IMAGE`, `TAG`, `CHART_DIR`.
- Tests and tools should be run at package root to operate across all packages: use `./...`.
- Keep container image determinism:
  - Ensure Dockerfile (root) builds reproducibly and the Helm values are updated via the Makefile when bumping image versions.

If you add tooling (linters, CI, image signing), add configs to the repo and document their usage here so Claude and contributors can follow the same workflow.
