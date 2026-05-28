# Changelog

All notable changes to this repository are documented in this file.

## Unreleased

### Documentation and skills

- README: describe HOTL-by-default (HITL at decision points) instead of HITL as the primary model.
- Skills 04–06, 08–09: align `ep-context.md` ownership with pipeline §4 (subagent reports `context_delta`; orchestrator applies; solo mode may update directly).
- Skills 01–02: add Orchestrator brief (subagent mode) sections.
- pipeline.spec.md: fix artefact path wording; skill 02: fix duplicate Core Principles numbering.
- README: `ai-sdlc.version` pin example and consumer CI verification snippet; CONTRIBUTING cross-link.

### Validator and CI

- `validate structure`: check `ep-context.md`, review artefacts (gate sections), and optional `ep-code-review.md` / `ep-audit-report.md`.
- `validate pipeline`: warn on unchecked plan tasks when stage 10 exists; error when stage 11 exists with unchecked tasks.
- VALIDATION.md: document hard vs soft enforcement model.
- CI: trigger workflow on `AGENTS.md` and `CONTRIBUTING.md` changes.

## v1.0.0 - 2026-05-28

- Extract canonical `ai-sdlc` repository with merged history from `fireman` and `PersonalAssistant`.
- Keep `fireman` content as canonical merge resolution while preserving both lineages.
- Document workspace-only consumption contract with version pinning (`ai-sdlc.version`).
