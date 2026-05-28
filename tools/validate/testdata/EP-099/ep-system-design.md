---
artefact: ep-system-design
epic_id: EP-099
status: approved
source_of_truth: true
updated_at: 2026-05-21
---

# EP-099: System Design

## Overview

Extends the validate CLI tool with new subcommands.

## Architecture diagram

The tool uses a subcommand dispatch pattern.

## Module boundaries

All validators are in package main under ai-sdlc/tools/validate/.

## Components

- req_ac_trace.go — REQ to AC traceability
- artefact_structure.go — structure validator
- ears_lint.go — EARS linter
- pipeline_state.go — pipeline gate checker

## Data models

Each validator returns a structured result with errors and warnings.

## Traceability table

| REQ | Component | Notes |
|-----|-----------|-------|
| REQ-99.001 | helpers_test.go | Fixture loading |
| REQ-99.003 | ears_lint.go | EARS patterns |
| REQ-99.004 | artefact_structure.go | Structure checks |
