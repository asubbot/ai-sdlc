---
# OKF v0.1
type: Pattern
title: Sync vs Async Communication
description: Choose between synchronous request/response and asynchronous
  messaging or request-reply for slow, unreliable, or decoupled work.
timestamp: 2026-07-24
tags: [messaging, decoupling, latency]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/async-request-reply
    note: Primary upstream (Asynchronous Request-Reply pattern, Azure)
  - url: https://learn.microsoft.com/en-us/azure/architecture/guide/technology-choices/messaging
    note: Asynchronous messaging options overview (Microsoft Learn)
forces: [caller latency budget, failure isolation, ordering, operational complexity of queues]
when: a caller must not block on slow or unreliable downstream work, or producers
  and consumers evolve/scale independently
when_not: the caller needs the result to proceed and the callee is fast and
  reliable; async here only adds queues, retries, and eventual-consistency cost
kiss_default: synchronous call with timeout and retry first; go async only when
  latency budget, burst absorption, or failure isolation demands it
quality: [responsiveness, resilience]
related: [retry-and-timeouts, idempotency, transactional-outbox, publisher-subscriber]
---

# Sync vs Async Communication

**Problem.** Blocking a caller on slow or flaky work couples availability and
latency of two components; but queues bring duplication, ordering, and ops cost.

**Options.**
- Sync request/response with timeout + retry — simplest, fine for fast reliable calls.
- Async request-reply: accept request, return status endpoint/callback — for slow work
  behind an API.
- Fire-and-forget via queue/worker — for work whose result the caller never needs inline.

**Failure modes.** Hidden coupling via synchronous fan-out chains; async adopted
without idempotent consumers (duplicates corrupt state); unbounded queues masking
a consumer that cannot keep up.
