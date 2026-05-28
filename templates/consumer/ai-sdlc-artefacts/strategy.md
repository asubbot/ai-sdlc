---
artefact: strategy
status: draft
source_of_truth: true
updated_at: 2026-05-28
---

# Delivery and test strategy (stub)

## Introduction

Minimal strategy for greenfield adoption. Aligns with [scope.md](scope.md). Refine in stage 2 when the product stack and increments are known.

## 1. Delivery strategy

### 1.1 Increments

1. **Adopt SDLC (EP-000)** — Bootstrap layout, pin, validator, CI, artefact conventions.
2. **Product MVP (EP-001+)** — To be planned after scope refinement.

### 1.2 Success criteria

- EP-000 completes with validator and CI smoke passing.
- Stages 1–2 produce agreed scope and strategy before product epics.

## 2. Test strategy

### 2.1 Levels

- **Process / infra:** `make validate`, CI pin verification (EP-000).
- **Product:** Unit / integration tests when `cmd/`, `internal/`, or `tests/` exist (EP-001+).

### 2.2 AC coverage

- Run `./bin/validate EP-XXX` from repo root before treating an epic as AC-complete.
- Use **Deferred** in `ep-acceptance-criteria.md` for bootstrap-only ACs without product tests (see [VALIDATION.md](../ai-sdlc/tools/validate/VALIDATION.md)).
