---
artefact: ep-acceptance-criteria
epic_id: EP-000
status: draft
source_of_truth: true
updated_at: 2026-05-28
---

# EP-000: Acceptance Criteria

## AC Index

| AC | REQ | Summary |
|----|-----|---------|
| [AC-00.001](#ac-00001) | REQ-00.001 | Pin file present and well-formed |
| [AC-00.002](#ac-00002) | REQ-00.002 | Consumer layout and artefact files |
| [AC-00.003](#ac-00003) | REQ-00.003 | Validator builds |
| [AC-00.004](#ac-00004) | REQ-00.004 | CI verifies pin and runs gates |
| [AC-00.005](#ac-00005) | REQ-00.005 | Project-wide make validate gate |
| [AC-00.006](#ac-00006) | REQ-00.006 | make check product gate |

## Scenarios

### AC-00.001 Pin file (Trace: REQ-00.001)

Given a clone of the consumer repository
When the operator reads `ai-sdlc.version`
Then it contains a non-empty tag or 40-character commit SHA without whitespace
And WHEN `ai-sdlc/` is checked out locally at a commit pin, THEN `git -C ai-sdlc rev-parse HEAD` matches the pin

### AC-00.002 Consumer layout (Trace: REQ-00.002)

Given the repository root
When the operator inspects the consumer layout
Then `AGENTS.md`, `.gitignore`, `Makefile`, `ai-sdlc.version`, `.golangci.yml`, `scripts/check-module-boundaries.sh`, and `.github/workflows/ci.yml` exist at the root, `.gitignore` lists `ai-sdlc/`, `.github/workflows/ai-sdlc.yml` does not exist, and under `ai-sdlc-artefacts/` exist `scope.md`, `strategy.md`, and `epics/EP-000/ep-scope.md`, `epics/EP-000/ep-requirements.md`, `epics/EP-000/ep-acceptance-criteria.md`, and `epics/EP-000/ep-context.md`

### AC-00.003 Validator build (Trace: REQ-00.003)

Given `ai-sdlc/` is present at the pinned revision
When the operator runs `make build` with no existing `bin/validate`
Then `bin/validate` is created and executes without error

### AC-00.004 Product CI (Trace: REQ-00.004)

Given a push to the default branch
When GitHub Actions runs `.github/workflows/ci.yml`
Then the pin is verified, `ai-sdlc` is checked out at the pin, and steps run `make build`, `make check`, and `make validate`

Bootstrap smoke in `tests/bootstrap_test.go` asserts `ci.yml` defines named steps (Verify ai-sdlc pin, Checkout ai-sdlc at pin, Build validate binary, Product check gate, Product validate gate) and required commands including `make build`, `make check`, and `make validate`.

### AC-00.005 Project make validate gate (Trace: REQ-00.005)

Given `bin/validate` is available at the repository root
When the operator runs `make validate` with no extra goals
Then AC coverage passes for all epics and pipeline/structure pass for each epic (structure may emit warnings for optional stages not yet created)

### AC-00.006 make check gate (Trace: REQ-00.006)

Given a product `go.mod` exists at the repository root
When the operator runs `make check`
Then `go fmt`, `go vet`, `govulncheck`, `golangci-lint`, and race-enabled tests with coverage on `tests/` complete without errors
And WHEN `go test ./tests/...` runs as part of `make check`, THE environment variable `BOOTSTRAP_MAKE_CHECK_INVOKED` is set so bootstrap tests do not recurse into `make check`
