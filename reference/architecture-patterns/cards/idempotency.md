---
# OKF v0.1
type: Pattern
title: Idempotency
description: Design operations so that processing the same request or message
  more than once has the same effect as processing it once.
timestamp: 2026-07-24
tags: [consistency, retries, deduplication]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/
    note: Primary upstream (AWS Builders' Library)
  - url: https://microservices.io/patterns/communication-style/idempotent-consumer.html
    note: Idempotent consumer pattern (microservices.io)
forces: [at-least-once delivery, retry safety, dedup storage cost, key lifetime]
when: an operation with side effects can be retried or delivered more than once
  (queues, retries, user double-submit)
when_not: the operation is naturally idempotent (pure read, absolute overwrite)
  — extra dedup machinery adds state for nothing
kiss_default: prefer naturally idempotent designs (set/upsert semantics, natural
  keys); add explicit idempotency keys + dedup store only for genuine
  create/charge-style effects
quality: [consistency, reliability]
related: [retry-and-timeouts, transactional-outbox, sync-vs-async]
---

# Idempotency

**Problem.** Retries and at-least-once delivery mean the same command arrives
twice; without idempotency this creates duplicates or double side effects.

**Options.**
- Naturally idempotent operations: upsert by natural key, absolute state writes.
- Idempotency key supplied by the caller + dedup table with response replay.
- Consumer-side dedup: processed-message ids stored transactionally with effects.

**Failure modes.** Dedup record written in a different transaction than the
effect (crash window reopens duplicates); unbounded key store without TTL;
keys scoped too broadly so distinct requests collide.
