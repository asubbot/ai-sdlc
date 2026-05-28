# Changelog

All notable changes to this repository are documented in this file.

## Unreleased

_No changes yet._

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
