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

## Scope (features/capabilities)

- Consumer `AGENTS.md`, `.gitignore` (layout A), `Makefile`, and `ai-sdlc-artefacts/` layout.
- `ai-sdlc.version` records the process revision; CI verifies the pin and checks out `ai-sdlc` at that revision.
- `make build` and `make validate` work from the product repository root.
- Project-level [scope.md](../../scope.md) and [strategy.md](../../strategy.md) stubs exist for pipeline gates.

## Success criteria

- Operator can run `make validate pipeline EP-000` without errors when epic artefacts are complete through stage 7 gate pass.
- CI workflow `ai-sdlc.yml` passes on the default branch.

## Traceability

- **Scope:** Adopting ai-sdlc and artefact layout from [scope.md](../../scope.md).
- **Strategy:** First increment in [strategy.md](../../strategy.md) (EP-000 adoption).
