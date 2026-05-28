---
artefact: ep-requirements
epic_id: EP-099
status: approved
source_of_truth: true
updated_at: 2026-05-20
---

# EP-099: Requirements

## Introduction

Requirements for test fixture epic.

## Glossary

| Term | Definition |
|------|-----------|
| Validator | A check that verifies artefact quality |

## C4 Context (C1)

System context is the validate CLI tool.

## Requirements

### REQ-99.001 — Fixture file loading
WHEN the validate tool runs, THE system SHALL load fixture files from the testdata directory.

### REQ-99.002 — Golden file comparison
WHEN the validate tool produces JSON output, THE system SHALL compare it against golden files.

### REQ-99.003 — EARS compliance check
WHEN linting requirements, THE system SHALL validate that each requirement follows one of the six EARS patterns.

### REQ-99.004 — Structure validation
WHEN an artefact is saved, THE system SHALL verify required sections are present.

### REQ-99.005 — Pipeline gate checking
WHEN the orchestrator requests gate status, THE system SHALL parse YAML front matter and gate summaries.

### REQ-99.006 — Broken link detection
WHEN validating artefact structure, THE system SHALL verify that all markdown links point to existing files.

## REQ Index

| ID | Type | Summary |
|----|------|---------|
| REQ-99.001 | FR | Fixture file loading |
| REQ-99.002 | FR | Golden file comparison |
| REQ-99.003 | FR | EARS compliance check |
| REQ-99.004 | FR | Structure validation |
| REQ-99.005 | FR | Pipeline gate checking |
| REQ-99.006 | FR | Broken link detection |
