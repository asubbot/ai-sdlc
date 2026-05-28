---
artefact: ep-requirements
epic_id: EP-000
status: draft
source_of_truth: true
updated_at: 2026-05-28
---

# EP-000: Requirements

## Introduction

Requirements to adopt ai-sdlc in this consumer repository.

## Glossary

| Term | Definition |
|------|------------|
| Consumer layout | Product root with `AGENTS.md`, `.gitignore`, `Makefile`, `ai-sdlc.version`, `ai-sdlc-artefacts/`, and local `ai-sdlc/` checkout (not committed; listed in `.gitignore`) |

## Requirements

### REQ-00.001 — Process pin
THE consumer repository SHALL record the canonical ai-sdlc revision in `ai-sdlc.version` at the repository root.

### REQ-00.002 — Consumer layout and artefacts
THE consumer repository SHALL provide `AGENTS.md`, `.gitignore`, and `Makefile` at the repository root, SHALL list `ai-sdlc/` in `.gitignore`, and SHALL store pipeline outputs under `ai-sdlc-artefacts/` at the repository root, not inside `ai-sdlc/`.

### REQ-00.003 — Validator invocable
WHEN the operator runs `make build` from the repository root, THE system SHALL produce `bin/validate` from `ai-sdlc/tools/validate`.

### REQ-00.004 — CI pin verification
WHEN CI runs on the default branch, THE workflow SHALL verify that `ai-sdlc.version` exists in the configured upstream repository, SHALL checkout `ai-sdlc` at that pin, SHALL run `make build`, `make check`, and project `make validate`.

### REQ-00.005 — Project validate gate
WHEN the operator runs `make validate` from the repository root with no extra goals, THE system SHALL run AC coverage for all epics and pipeline/structure checks for each epic without errors (structure warnings for not-yet-created optional stages are acceptable).

### REQ-00.006 — Product check gate
WHEN the operator runs `make check` from the repository root, THE system SHALL run `go vet` on the product module and bootstrap tests under `tests/` without errors.

## REQ Index

| ID | Type | Summary |
|----|------|---------|
| REQ-00.001 | FR | Process pin file |
| REQ-00.002 | FR | Consumer root files and artefact directory |
| REQ-00.003 | FR | Validator build |
| REQ-00.004 | FR | CI pin, checkout, and gates |
| REQ-00.005 | FR | Project-wide make validate gate |
| REQ-00.006 | FR | make check product gate |
