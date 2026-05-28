---
artefact: ep-scope
epic_id: EP-000
status: draft
source_of_truth: true
updated_at: 2026-05-28
---

# Epic scope — EP-000 Adopt ai-sdlc

| Field | Content |
|-------|---------|
| **ID** | EP-000 |
| **Status** | NEW |
| **Title** | Adopt ai-sdlc |
| **Description** | Wire this repository to the pinned ai-sdlc process: consumer layout, validator, CI, and artefact conventions. |
| **First version date** | 2026-05-28 |

## Glossary

- **Pin** — `ai-sdlc.version` file referencing a tag or commit in the canonical process repository.
- **Validator** — `bin/validate` built from `ai-sdlc/tools/validate`.
- **Project validate gate** — `make validate` with no extra goals: AC coverage for all epics, then pipeline and structure for each epic under `ai-sdlc-artefacts/epics/`.

## Scope (features/capabilities)

- Consumer `AGENTS.md`, `.gitignore` (includes `ai-sdlc/`), `Makefile`, `.golangci.yml`, `scripts/check-module-boundaries.sh`, and `ai-sdlc-artefacts/` layout.
- `ai-sdlc.version` records the process revision; product CI (`.github/workflows/ci.yml`) verifies the pin and checks out `ai-sdlc` at that revision.
- Product gates `make build`, `make validate`, and `make check` work from the product repository root.
- Project-level [scope.md](../../scope.md) and [strategy.md](../../strategy.md) stubs exist for pipeline gates.

## Success criteria

- Operator can run `make validate` from the repository root without errors (project validate gate; structure warnings for not-yet-created optional stages are acceptable).
- Operator can run `make check` from the repository root without errors (fmt, vet, vuln, lint, race tests and coverage on `tests/`).
- CI workflow `ci.yml` passes on the default branch (`make build`, `make check`, and `make validate`).

## Traceability

- **Scope:** Adopting ai-sdlc and artefact layout from [scope.md](../../scope.md).
- **Strategy:** First increment in [strategy.md](../../strategy.md) (EP-000 adoption).
