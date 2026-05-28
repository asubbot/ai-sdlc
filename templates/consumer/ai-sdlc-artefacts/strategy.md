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

- **Product gates (this repo):** `make check` (code quality), `make validate` (artefacts under `ai-sdlc-artefacts/`), enforced in `.github/workflows/ci.yml` after pin checkout.
- **Process regression (canonical ai-sdlc repo only):** `bootstrap-smoke.sh` and validator unit tests — reference when bumping `ai-sdlc.version`.
- **Product tests:** Bootstrap tests in `tests/` today; unit / integration / e2e when `cmd/`, `internal/`, or additional test dirs exist (EP-001+; uncomment Makefile targets).

### 2.2 AC coverage

- Run `make validate` from repo root for the **project gate** (AC coverage for all epics, then pipeline and structure per epic) before treating adoption complete.
- Run `./bin/validate EP-XXX` (or `make validate EP-XXX`) for a single-epic AC report when debugging.
- Use **Deferred** in `ep-acceptance-criteria.md` for bootstrap-only ACs without product tests (see [VALIDATION.md](../ai-sdlc/tools/validate/VALIDATION.md)).
