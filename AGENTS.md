# Project Instructions

Instructions for AI agents working in a **git-backed software repository**.

## How to use this file
- This file contains persistent agent instructions for this repository.
- The normative SDLC pipeline lives in **[ai-sdlc/specification/pipeline.spec.md](ai-sdlc/specification/pipeline.spec.md)** and the stage skills it maps to.
- Use **[ai-sdlc/README.md](ai-sdlc/README.md)** as the directory index for `ai-sdlc/`.

## Cooperation with the user
- Work **with** the user. Ask them to choose when a decision has meaningful trade-offs in design, naming, artefact placement, scope, or interpretation. For small local implementation details, choose the simplest option consistent with the repo.
- Match the user's chat language unless they ask otherwise.

## Principles
- **KISS** — prefer the smallest change that solves the problem; avoid unnecessary abstraction and scope creep.
- **Fail fast** — detect invalid state and errors early; do not swallow failures without a clear, documented reason.
- Do not add `//nolint:gocyclo`.

## Language
- All code comments, UI/user-facing messages in the product, and commit messages must be in English.

## Research / Docs-first
- For third-party libraries, APIs, and platforms, use official documentation first, keeping the user's keywords where useful. Prefer official docs over issues or blogs; fall back only if official docs are insufficient.

## Safety
- A user task is permission to edit files within the requested scope. Do not make unrelated product-source or build-configuration changes without explicit approval.
- Write or update delivery-process artefacts (requirements, design, plans, reviews, etc.) only through the repository's defined SDLC process and skills.
- Do not commit without explicit permission. Commit messages must be in English; when helpful, reference the skill or plan step.
- Merge is not push. Do not push to a remote unless the user explicitly asks for that remote update in the current request or otherwise clearly authorizes it for this step.
- When asked to merge, use a real merge commit; do not fast-forward.
- Never commit real tokens, passwords, private keys, or secrets. Use placeholders and documented configuration patterns.
- Do not weaken security or reliability for convenience without an explicit trade-off discussion with the owner.

## This Repository (`ai-sdlc`)
- **Product code layout:** treat **`cmd/`**, **`internal/`**, **`tests/`**, and project **build files** (e.g. `Makefile`, Go module files) as product source unless the owner narrows scope further.
- **Agentic SDLC:** definitions live under **`ai-sdlc/`**. Use **[ai-sdlc/README.md](ai-sdlc/README.md)** as the **directory index** (what each path is for). The **step-by-step pipeline** (stages, delegation, artefact rules, AC validation commands) is only in **[ai-sdlc/specification/pipeline.spec.md](ai-sdlc/specification/pipeline.spec.md)** and the **`*.skill.md`** files it maps to.
- **Pipeline outputs:** documents produced by the SDLC (scope, strategy, epic files, etc.) go under **`ai-sdlc-artefacts/`** at the repo root, not inside `ai-sdlc/`.
- **Checks after substantive code edits:** when you may run commands, run **`make check`** after non-trivial changes unless the owner says otherwise.
- **Explicit JSON configuration:** product **`config.json`** must list every documented top-level key exactly once. Optional product blocks are disabled with JSON `null`, not by omitting the key. Unknown top-level keys, missing keys, invalid values, or structural drift must fail config load.
- **Optional editor tooling:** Sourcerer MCP may be used for semantic navigation, but not instead of required repository docs.
- **Sensitive domains in this codebase:** treat **config paths**, **SSH**, **Telegram**, and **LLM logs** as high-sensitivity when touching related code; follow epic requirements (e.g. redaction, allowlists) where they apply.
