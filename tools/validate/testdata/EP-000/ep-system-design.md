---
artefact: ep-system-design
epic_id: EP-000
status: approved
source_of_truth: true
updated_at: 2026-05-28
---

# EP-000: System Design

## Overview

Consumer repository layout: product root, `ai-sdlc/` process clone, `ai-sdlc-artefacts/`, `bin/validate`.

## Components

- `ai-sdlc.version` — process pin
- `ai-sdlc-artefacts/` — pipeline outputs
- `bin/validate` — validator CLI

## Traceability

| REQ | Component | Notes |
|-----|-----------|-------|
| REQ-00.001 | ai-sdlc.version | Pin file |
| REQ-00.002 | ai-sdlc-artefacts/ | Artefact root |
