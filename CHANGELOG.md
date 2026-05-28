# Changelog

All notable changes to this repository are documented in this file.

## Unreleased

### Greenfield consumer bootstrap

- [00-project-bootstrap.skill.md](specification/skills/00-project-bootstrap.skill.md) — pre-pipeline skeleton for new product repositories.
- [templates/consumer/](templates/consumer/) — `AGENTS.md`, `.gitignore`, `Makefile`, `ai-sdlc.version`, `.golangci.yml`, `scripts/check-module-boundaries.sh`, product CI (`.github/workflows/ci.yml`), aggressive `make check`, `ai-sdlc-artefacts/` stubs including **EP-000** adoption epic.
- Consumer template: removed `.github/workflows/ai-sdlc.yml`; product gates (`make check`, `make validate`) run in `ci.yml` only. Process regression stays in canonical repo (`bootstrap-smoke.sh`).
- [pipeline.spec.md](specification/pipeline.spec.md) — **Consumer onboarding (greenfield)**; gitignored nested `ai-sdlc/` clone; normative order bootstrap → stages 1–2 → EP-000 → product epics.
- [README.md](README.md) — **Starting a new project** walkthrough and target layout summary.
- Validator testdata: [tools/validate/testdata/EP-000/](tools/validate/testdata/EP-000/), project-level [scope.md](tools/validate/testdata/scope.md) / [strategy.md](tools/validate/testdata/strategy.md); CI smoke for EP-000 and EP-099.
- Skills: gate in [01-scope-analysis.skill.md](specification/skills/01-scope-analysis.skill.md); consumer `bin/validate` note in [09-task-execution.skill.md](specification/skills/09-task-execution.skill.md).
- [tools/scripts/bootstrap-smoke.sh](tools/scripts/bootstrap-smoke.sh) — automated greenfield bootstrap smoke; CI replaces partial Makefile-only consumer smoke.
- Consumer onboarding: **single process layout** only (gitignored nested `ai-sdlc/` clone); submodule layout removed from normative docs and bootstrap skill.

### Upgrading to Unreleased (consumer greenfield)

| Topic | Action |
|-------|--------|
| New projects | Follow README **Starting a new project** + bootstrap skill |
| Existing repos | No change required; bootstrap is greenfield-only |
| EP-000 | Optional for mature repos adopting explicit SDLC infra epic |
| Process clone | Add `ai-sdlc/` to `.gitignore`; clone locally and in CI at pin in `ai-sdlc.version` (see template workflow) |

## v1.0.1 - 2026-05-28

### Documentation and skills

- README: describe HOTL-by-default (HITL at decision points) instead of HITL as the primary model.
- Skills 04–06, 08–09: align `ep-context.md` ownership with pipeline §4 (subagent reports `context_delta`; orchestrator applies; solo mode may update directly).
- Skills 01–02: add Orchestrator brief (subagent mode) sections.
- pipeline.spec.md: fix artefact path wording; skill 02: fix duplicate Core Principles numbering.
- README: `ai-sdlc.version` pin example and consumer CI verification snippet; CONTRIBUTING cross-link.
- pipeline.spec.md §4.5 and VALIDATION.md: aligned enforcement model tables; CI-minimal decision record documented.
- skills/README.md: `ep-context.md` required vs recommended sections; HITL decision record CI subset.

### Validator and CI

- `validate structure`: check `ep-context.md` (including **Open Questions**), review artefacts (gate sections), and optional `ep-code-review.md` / `ep-audit-report.md`.
- `validate pipeline`: warn on unchecked plan tasks when stage 10 exists; error when stage 11 exists with unchecked tasks.
- VALIDATION.md: document hard vs soft enforcement model; artefact structure section; architecture source list.
- CI: trigger workflow on `AGENTS.md` and `CONTRIBUTING.md` changes; smoke `validate pipeline` / `structure` against test fixtures in consumer layout.

### Upgrading from v1.0.0

Consumer repositories updating their `ai-sdlc.version` pin should:

| Topic | Action in consumer repo |
|-------|-------------------------|
| Validator path | Invoke `./tools/validate/validate` from the product repo root; build with `cd tools/validate && go build -o validate .` |
| HOTL default | Agents proceed on routine work; HITL only at decision points in [pipeline.spec.md](specification/pipeline.spec.md) |
| `validate pipeline` | Stage 11 without `ep-code-review.md` fails; unchecked plan tasks when `ep-audit-report.md` exists fail; unchecked tasks when only `ep-code-review.md` exists are **warnings** (exit 0) |
| `validate structure` | `ep-context.md` requires `Purpose`, `Current Scope`, `Open Questions`, `Links`; review artefacts need Current Gate Summary and Review iteration sections |
| Gate override | Review artefact must include `Decision needed:` and `Operator choice:` (full decision record remains process guidance) |
| AC trace | Top-level `Test*` without an AC trace line fails `validate` / `validate ac` |

## v1.0.0 - 2026-05-28

- Extract canonical `ai-sdlc` repository with merged history from `fireman` and `PersonalAssistant`.
- Keep `fireman` content as canonical merge resolution while preserving both lineages.
- Document workspace-only consumption contract with version pinning (`ai-sdlc.version`).
