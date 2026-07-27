---
# OKF v0.1
type: Pattern
title: Publisher–Subscriber
description: Publishers emit events to a broker topic; subscribers consume
  independently so producers need not know their consumers.
timestamp: 2026-07-27
tags: [messaging, decoupling, events]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/publisher-subscriber
    note: Primary upstream (Publisher-Subscriber pattern, Azure)
forces: [producer-consumer decoupling, fan-out, eventual consistency, broker ops]
when: multiple consumers need the same facts, or producers and consumers must
  evolve and scale independently
when_not: there is exactly one consumer that must reply inline; a direct call
  or simple queue is simpler and clearer
kiss_default: direct call or single queue first; introduce pub/sub when a second
  independent consumer appears or fan-out becomes a real requirement
quality: [scalability, decoupling]
related: [sync-vs-async, transactional-outbox, dead-letter-queue, idempotency]
---

# Publisher–Subscriber

**Problem.** Point-to-point calls couple producers to every consumer; adding a
new listener means changing the producer and often the deployment unit.

**Options.**
- Direct call / single consumer queue — simplest for one listener.
- Publisher–subscriber via a topic/broker — fan-out with independent subscribers.
- In-process event bus — fine inside one process; does not survive restarts or scale-out.

**Failure modes.** Implicit contracts with no schema/versioning; at-least-once
delivery without idempotent subscribers; topics used where request/response was needed.
