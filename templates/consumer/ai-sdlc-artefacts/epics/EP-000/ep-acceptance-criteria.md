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
| [AC-00.001](#ac-00001) | REQ-00.001 | Pin file present |
| [AC-00.002](#ac-00002) | REQ-00.002 | Artefacts directory layout |
| [AC-00.003](#ac-00003) | REQ-00.003 | Validator builds |
| [AC-00.004](#ac-00004) | REQ-00.004 | CI verifies pin and tests validate tool |

## Scenarios

### AC-00.001 Pin file (Trace: REQ-00.001)

Given a clone of the consumer repository
When the operator reads `ai-sdlc.version`
Then it contains a non-empty tag or commit SHA

### AC-00.002 Artefact layout (Trace: REQ-00.002)

Given the repository root
When the operator lists `ai-sdlc-artefacts/`
Then `scope.md`, `strategy.md`, and `epics/EP-000/` exist

### AC-00.003 Validator build (Trace: REQ-00.003)

Given `ai-sdlc/` is present at the pinned revision
When the operator runs `make build`
Then `bin/validate` exists and executes without error

**Deferred:** No unit test trace in product `tests/` until product code exists (bootstrap gate; verified by CI and manual `make build`).

### AC-00.004 CI adoption (Trace: REQ-00.004)

Given a push to the default branch
When GitHub Actions runs `ai-sdlc.yml`
Then the pin is verified and `go test` in `ai-sdlc/tools/validate` succeeds

**Deferred:** CI-only; no duplicated unit test in product tree for this AC.
