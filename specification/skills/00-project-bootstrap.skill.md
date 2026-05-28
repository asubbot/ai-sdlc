---
name: project-bootstrap.skill
description: >-
  Greenfield consumer onboarding: materialize product repo layout (AGENTS.md,
  ai-sdlc.version, ai-sdlc-artefacts stubs, Makefile, CI) when adopting ai-sdlc.
  Use when starting a new project, bootstrapping ai-sdlc, or when only ai-sdlc/
  exists without consumer AGENTS.md or scope.md.
---

# Bootstrap: Consumer project (greenfield)

**Pipeline:** Pre-stage — see [pipeline.spec.md](../pipeline.spec.md) (Consumer onboarding).

**Output:** Consumer repository skeleton at the product root (not inside `ai-sdlc/`).

---

## Core Principles

1. **Product root workspace** — All paths below are relative to the **consumer product repository root** (the directory that will contain `ai-sdlc-artefacts/` and `ai-sdlc.version`). If the workspace root is `ai-sdlc/` only, instruct the operator to reopen the parent folder or create the product root first.
2. **Do not apply canonical ai-sdlc AGENTS rules** — In a consumer repo, creating `ai-sdlc-artefacts/` is required. Ignore canonical-repo bans on artefacts.
3. **Templates source** — Copy or adapt from [`templates/consumer/`](../../templates/consumer/) in the pinned `ai-sdlc/` tree unless the operator provides a fork path.
4. **Process clone** — Nested clone at `ai-sdlc/` listed in `.gitignore`; durable process reference is `ai-sdlc.version` only (not the vendored tree in product git).
5. **Bootstrap ends at skeleton** — Full **EP-000** pipeline (stages 3–11) runs **after** stages 1–2, in separate orchestration sessions if needed.

---

## Orchestrator brief

- **Required input:** `ai-sdlc/` directory present (local clone) with `specification/pipeline.spec.md`.
- **Greenfield detection:** Missing consumer root `AGENTS.md` **or** missing `ai-sdlc-artefacts/scope.md` while adopting a new product.
- **Gate check:** Abort if `ai-sdlc/specification/pipeline.spec.md` is missing.
- **Output signal:** `BOOTSTRAP_COMPLETE: skeleton ready [<paths created>]`
- **Validation after:** `make build` succeeds when `ai-sdlc/` is present; optional `make validate structure EP-000` when EP-000 stubs exist.

---

## Workflow

Follow this order:

1. **Verify process tree** — Confirm `ai-sdlc/specification/pipeline.spec.md` exists. If not, stop with a clear error (clone canonical ai-sdlc into `ai-sdlc/`).

2. **Clone if missing** — If `ai-sdlc/` is absent, run `git clone <AI_SDLC_REPO_URL> ai-sdlc` and checkout the revision from `ai-sdlc.version` (or ask operator for URL and pin). Record clone instructions in the product README if clone cannot run in-session.

3. **Materialize skeleton** — From `ai-sdlc/templates/consumer/`, copy to product root:

   | Template | Product path |
   |----------|----------------|
   | `AGENTS.md` | `AGENTS.md` |
   | `.gitignore` | `.gitignore` |
   | `Makefile` | `Makefile` |
   | `ai-sdlc.version` | `ai-sdlc.version` |
   | `README.project.md` | `README.md` |
   | `ai-sdlc-artefacts/**` | `ai-sdlc-artefacts/**` |
   | `.github/workflows/ai-sdlc.yml` | `.github/workflows/ai-sdlc.yml` |

   Do not overwrite existing operator content without HITL.

4. **Write pin** — Set `ai-sdlc.version` to `git -C ai-sdlc rev-parse HEAD` or an agreed release tag (e.g. `v1.0.1`). One line, no extra whitespace.

5. **Build validator** — Run `make build` from product root. Fail with actionable message if `ai-sdlc/tools/validate` is missing.

6. **Git** — If no `.git/`, suggest `git init` (HITL if remote/origin policy matters).

7. **Signal completion** — Emit `BOOTSTRAP_COMPLETE: skeleton ready [<paths>]`.

8. **Hand off (normative order)** — Tell the operator:

   1. **Stage 1** — [01-scope-analysis.skill.md](01-scope-analysis.skill.md) (refine `scope.md`; stubs are starting points).
   2. **Stage 2** — [02-strategy-analysis.skill.md](02-strategy-analysis.skill.md) (refine `strategy.md`).
   3. **EP-000** — Stages 3→11 for adoption epic (mandatory); stage 3 requires scope + strategy from steps 1–2.
   4. **EP-001+** — First product epic after scope/strategy reflect the real product.

---

## EP-000 orchestration (after stages 1–2)

Not part of bootstrap **Done when**; run as follow-on work:

- Branch: `epic/EP-000-adopt-sdlc` (or stage 3 variant B naming).
- Stages 3→8: HOTL; use template stubs under `ai-sdlc-artefacts/epics/EP-000/` as baseline.
- Stages **7** and **10**: **delegate** per [pipeline.spec.md](../pipeline.spec.md) §3.
- Stage **9**: May verify scaffold only (Makefile, CI, `make build`); no product code required.
- **Deferred** ACs allowed per [VALIDATION.md](../../tools/validate/VALIDATION.md) when product `tests/` / `internal/` / `cmd/` are absent.

---

## Done when (bootstrap skill only)

- [ ] Consumer `AGENTS.md` exists at product root.
- [ ] `.gitignore` exists and includes `ai-sdlc/`.
- [ ] `ai-sdlc.version` contains a pin matching the `ai-sdlc/` checkout.
- [ ] `ai-sdlc-artefacts/README.md` exists.
- [ ] `Makefile` and `.github/workflows/ai-sdlc.yml` exist.
- [ ] `make build` produces `bin/validate` when `ai-sdlc/` is present.
- [ ] Operator briefed on next steps (stages 1–2, then EP-000).

**Not required for bootstrap Done when:** complete EP-000 through stage 11; refined product scope; product source code.

---

## Maintainer verification (canonical repository only)

After changes to [`templates/consumer/`](../../templates/consumer/) or this skill, maintainers in the **canonical ai-sdlc** repository SHOULD run:

```bash
./tools/scripts/bootstrap-smoke.sh
```

See [tools/scripts/README.md](../../tools/scripts/README.md). This is **not** part of consumer bootstrap Done when; it regression-tests that templates still materialize and `make build` / `validate pipeline|structure EP-000` succeed.

---

## First message to agent (operator hint)

> New consumer project. Run project bootstrap, then stage 1 scope analysis — ask me what we are building.
