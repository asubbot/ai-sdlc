# Project Instructions — canonical ai-sdlc

Instructions for AI agents working in the **canonical ai-sdlc** repository (shared agentic SDLC process). **Consumer product repositories** maintain their own root `AGENTS.md` for product code, `make check`, and `ai-sdlc-artefacts/` — see [specification/pipeline.spec.md](specification/pipeline.spec.md) (Agent execution expectations).

## How to use this file

- This file applies when the workspace root **is** this repository (the canonical ai-sdlc checkout itself).
- **New consumer product repositories** use [templates/consumer/AGENTS.md](templates/consumer/AGENTS.md) and [00-project-bootstrap.skill.md](specification/skills/00-project-bootstrap.skill.md) — not the rules below about forbidding `ai-sdlc-artefacts/` at the repo root.
- The normative SDLC pipeline lives in **[specification/pipeline.spec.md](specification/pipeline.spec.md)** and the stage skills under **[specification/skills/](specification/skills/)**.
- Use **[README.md](README.md)** as the directory index. **[proposals.md](proposals.md)** is a living ideas document, not normative process.

## Cooperation with the user

- Work **with** the user. Ask them to choose when a decision has meaningful trade-offs in design, naming, artefact placement, scope, or interpretation. For small local implementation details, choose the simplest option consistent with the repo.
- Match the user's chat language unless they ask otherwise.

## Principles

- **KISS** — prefer the smallest change that solves the problem; avoid unnecessary abstraction and scope creep.
- **Fail fast** — detect invalid state and errors early; do not swallow failures without a clear, documented reason.
- Do not add `//nolint:gocyclo`.

## Language

- All code comments, commit messages, and normative process text in this repository must be in English.

## Research / Docs-first

- For third-party libraries, APIs, and platforms, use official documentation first, keeping the user's keywords where useful. Prefer official docs over issues or blogs; fall back only if official docs are insufficient.

## Safety

- A user task is permission to edit files within the requested scope. Do not make unrelated changes outside that scope without explicit approval.
- Do not commit without explicit permission. Commit messages must be in English; when helpful, reference the skill or plan step.
- Merge is not push. Do not push to a remote unless the user explicitly asks for that remote update in the current request or otherwise clearly authorizes it for this step.
- When asked to merge, use a real merge commit; do not fast-forward.
- Never commit real tokens, passwords, private keys, or secrets. Use placeholders and documented configuration patterns.
- Do not weaken security or reliability for convenience without an explicit trade-off discussion with the owner.

## This repository (canonical ai-sdlc)

- **Process edits:** Update [specification/pipeline.spec.md](specification/pipeline.spec.md) and related skills together; follow [CONTRIBUTING.md](CONTRIBUTING.md). Keep changes process-focused and repository-agnostic unless the owner approves otherwise.
- **Validator:** Implementation and docs live under **[tools/validate/](tools/validate/)**. After non-trivial validator changes, run `cd tools/validate && go test ./...` (and `go vet ./...` when appropriate).
- **Bootstrap templates:** After changes to [templates/consumer/](templates/consumer/) or [00-project-bootstrap.skill.md](specification/skills/00-project-bootstrap.skill.md), run `./tools/scripts/bootstrap-smoke.sh` in addition to validator tests.
- **Not product source here:** Do not treat `cmd/`, `internal/`, `tests/`, `config.json`, or a product `Makefile` as in-scope product code in this repository — they belong in consumer repos.
- **SDLC artefacts:** Do not create or update `ai-sdlc-artefacts/` in this repository except under `tools/validate/testdata/` for validator fixtures.
- **Optional editor tooling:** Sourcerer MCP may be used for semantic navigation, but not instead of required repository docs.

## Out of scope for this file

Product layout, `make check`, `config.json` rules, pinning `ai-sdlc.version`, and epic delivery artefacts are documented for **consumer** workspaces ([README.md](README.md) consumption model; consumer root `AGENTS.md`).
