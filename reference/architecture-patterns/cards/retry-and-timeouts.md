---
# OKF v0.1
type: Pattern
title: Retry and Timeouts
description: Bound every remote call with a timeout and retry transient faults
  with backoff, without amplifying load or duplicating effects.
timestamp: 2026-07-24
tags: [resilience, transient-faults, backoff]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/retry
    note: Primary upstream (Retry pattern, Azure)
  - url: https://learn.microsoft.com/en-us/azure/architecture/best-practices/transient-faults
    note: Transient fault handling best practices (Microsoft Learn)
forces: [transient vs permanent faults, retry storms, caller latency budget, duplicate side effects]
when: any call crosses a process or network boundary and can fail transiently
when_not: failures are deterministic (bad input, auth) — retrying repeats the same
  error and hides the root cause
kiss_default: single timeout + small fixed retry count with backoff and jitter;
  add budgets/policies per dependency only when call volume makes storms plausible
quality: [reliability, availability]
related: [circuit-breaker, idempotency]
---

# Retry and Timeouts

**Problem.** Remote calls fail transiently or hang; without timeouts callers
block indefinitely, and naive retries turn a blip into a self-inflicted outage.

**Options.**
- Timeout only, surface the error — acceptable when the operation is user-retryable.
- Timeout + bounded retries with exponential backoff and jitter — the default.
- Retry budget / policy per dependency — for high-volume services.

**Failure modes.** Retrying non-idempotent operations (duplicates); retry storms
when many callers back off in sync (no jitter); timeouts longer than the caller's
own deadline; retrying permanent errors and masking bugs.
