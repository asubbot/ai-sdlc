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
| Consumer layout | Product root with `ai-sdlc/`, `ai-sdlc-artefacts/`, and `ai-sdlc.version` |

## Requirements

### REQ-00.001 — Process pin
THE consumer repository SHALL record the canonical ai-sdlc revision in `ai-sdlc.version` at the repository root.

### REQ-00.002 — Artefact directory
THE consumer repository SHALL store pipeline outputs under `ai-sdlc-artefacts/` at the repository root, not inside `ai-sdlc/`.

### REQ-00.003 — Validator invocable
WHEN the operator runs `make build` from the repository root, THE system SHALL produce `bin/validate` from `ai-sdlc/tools/validate`.

### REQ-00.004 — CI pin verification
WHEN CI runs on the default branch, THE workflow SHALL verify that `ai-sdlc.version` exists in the configured upstream repository and SHALL checkout `ai-sdlc` at that pin before running validate tool tests.

## REQ Index

| ID | Type | Summary |
|----|------|---------|
| REQ-00.001 | FR | Process pin file |
| REQ-00.002 | FR | Artefact directory layout |
| REQ-00.003 | FR | Validator build |
| REQ-00.004 | FR | CI pin and checkout |
