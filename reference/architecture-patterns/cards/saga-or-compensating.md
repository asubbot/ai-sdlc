---
# OKF v0.1
type: Pattern
title: Saga (Compensating Transactions)
description: Coordinate a multi-step business process across resources with a
  sequence of local transactions and compensating actions on failure.
timestamp: 2026-07-27
tags: [consistency, orchestration, compensating-transactions]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
    note: Primary upstream (Saga pattern, Azure)
forces: [cross-resource consistency, compensation complexity, observability of long flows]
when: a business outcome needs several local commits (services or resources)
  and a distributed transaction is unavailable or undesirable
when_not: a single local transaction (or outbox) already covers the outcome;
  inventing compensations for a one-step flow is pure overhead
kiss_default: keep the workflow in one transaction or one service first; adopt
  a saga only when the steps genuinely span independent commit boundaries
quality: [consistency, reliability]
related: [transactional-outbox, idempotency, sync-vs-async]
---

# Saga (Compensating Transactions)

**Problem.** A business process spans multiple commit boundaries; a mid-flight
failure leaves partial state that must be undone or completed explicitly.

**Options.**
- Single local transaction / outbox — prefer when possible.
- Saga (choreography or orchestration) with compensating steps.
- Distributed transaction (2PC) — usually avoided for ops and availability cost.

**Failure modes.** Compensations that are not idempotent; missing compensation
for a new step; orchestration without timeouts leaving stuck sagas forever.
