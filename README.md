# ai-sdlc — canonical process repository

This repository is the canonical source of truth for the shared **agentic SDLC process** used by multiple product repositories.

`ai-sdlc` describes the process only; project execution outputs are stored in each product repository under `ai-sdlc-artefacts/`.

---

## Repository layout

| Path | Role |
|------|------|
| **[specification/pipeline.spec.md](specification/pipeline.spec.md)** | **Normative process:** stage order, inputs/outputs, stage→skill mapping, Human-in-the-loop, delegated execution, artefact naming, traceability, **agent execution expectations** (single process, AC validation, etc.). |
| **[specification/skills/](specification/skills/)** | **Per-stage agent instructions** (`01-` … `11-` plus optional cross-cutting skills). Each skill defines workflow and artefact structure for its stage. |
| **[specification/skills/README.md](specification/skills/README.md)** | Index and **common behaviour** across skills. |
| **[proposals.md](proposals.md)** | Proposed pipeline improvements and measurement notes; a living document, not the normative process. |
| **[tools/validate/](tools/validate/)** | AC↔test coverage checker (`./tools/validate/validate` from repo root); see [VALIDATION.md](tools/validate/VALIDATION.md) and [README](tools/validate/README.md). |

---

## Consumption model (workspace-only)

Projects consume this repository as a separate workspace root. To keep reproducibility:

1. Each consumer project stores an `ai-sdlc.version` file with a pinned tag or commit SHA.
2. CI in consumer projects verifies that the pinned revision exists in this canonical repository.
3. Process changes are made here via PR; product repositories update only their pinned version.
