# Project Instructions

Instructions for AI agents working in this **consumer product repository**.

## Workspace root

- Open Cursor (or equivalent) at the **product repository root** — the directory that contains `ai-sdlc-artefacts/` and `ai-sdlc.version`.
- If you only see `ai-sdlc/specification/` at the workspace root, you are inside the process clone; **stop** and reopen the parent product folder, or run [00-project-bootstrap.skill.md](ai-sdlc/specification/skills/00-project-bootstrap.skill.md) from a parent that will become the product root.

## Normative pipeline

- **Process:** [ai-sdlc/specification/pipeline.spec.md](ai-sdlc/specification/pipeline.spec.md) and stage skills under [ai-sdlc/specification/skills/](ai-sdlc/specification/skills/).
- **Index:** [ai-sdlc/README.md](ai-sdlc/README.md).
- **Outputs:** under `ai-sdlc-artefacts/` (not inside `ai-sdlc/`).

## Greenfield

If `AGENTS.md` (this file) or `ai-sdlc-artefacts/scope.md` is missing but `ai-sdlc/` exists, run [00-project-bootstrap.skill.md](ai-sdlc/specification/skills/00-project-bootstrap.skill.md) before stage 1.

Post-bootstrap order: bootstrap skeleton → stages **1–2** (scope, strategy) → **EP-000** (adopt SDLC) → refine scope/strategy → product epics (e.g. **EP-001**).

## Principles

- **KISS** — smallest change that solves the problem.
- **Fail fast** — detect invalid state early; do not swallow errors without reason.
- **Explicit configuration** — when this product uses a config file, document every top-level key; optional blocks use JSON `null`, not omission (add product-specific rules here when applicable).

## Language

- Code comments, commit messages, and in-repo technical docs: **English**.
- Chat language: match the operator unless they ask otherwise.

## Cooperation

- Ask the operator to choose when several valid options exist (design, naming, scope, approach).
- Do not commit without explicit allowance. **Merge ≠ push** unless explicitly requested.
- Use real merge commits, not fast-forward, when merging branches.
- Never commit secrets; use placeholders and documented patterns.

## Checks

- After substantive code edits, run **`make check`** when allowed.
- For AC and pipeline gates: **`make validate`** or **`./bin/validate`** from repo root (build with `make build`).

## Process pin

- Do **not** treat edits under `ai-sdlc/` as the durable process source — bump `ai-sdlc.version` and update the `ai-sdlc/` checkout to match.
