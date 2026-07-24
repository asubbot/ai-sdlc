---
# OKF v0.1
type: Pattern
title: Transactional Outbox
description: Atomically persist business state and an outgoing event in one
  local transaction; a separate relay delivers events.
timestamp: 2026-07-24
tags: [reliability, messaging, consistency]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://microservices.io/patterns/data/transactional-outbox.html
    note: Primary upstream (microservices.io; canonical pattern description)
  - url: https://learn.microsoft.com/en-us/azure/architecture/databases/guide/transactional-outbox-cosmos
    note: Transactional outbox implementation guide (Microsoft Learn)
forces: [atomicity of state+event, at-least-once delivery, no distributed transactions]
when: an event must be delivered if and only if the local transaction committed
when_not: single process where losing an event is acceptable, or the broker
  offers end-to-end transactional publish
kiss_default: start with direct call + retry; adopt outbox only when event loss
  is genuinely unacceptable
quality: [reliability, consistency]
related: [idempotency, retry-and-timeouts]
---

# Transactional Outbox

**Problem.** Writing to the database and publishing an event are two systems;
a crash between them loses the event or creates a phantom one.

**Options.**
- Direct publish after commit + retry — simplest; loses events on crash window.
- Transactional outbox — event row committed with the data; relay polls/streams it.
- Distributed transaction (2PC) — atomic but operationally heavy; avoid.

**Failure modes.** Relay duplicates deliveries (pair with idempotency);
outbox table growth without cleanup; ordering only per-aggregate.
