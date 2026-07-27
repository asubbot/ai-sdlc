---
# OKF v0.1
type: Pattern
title: Rate Limiting
description: Cap how many requests a client or dependency may make in a time
  window to protect capacity and control cost.
timestamp: 2026-07-27
tags: [resilience, capacity, fairness]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/rate-limiting-pattern
    note: Primary upstream (Rate Limiting pattern, Azure)
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/throttling
    note: Throttling pattern (related Azure guidance)
forces: [protect capacity vs reject legitimate load, fairness across clients, client backoff UX]
when: inbound traffic or outbound calls can burst beyond what you can afford
  (capacity, cost, or partner quotas)
when_not: traffic is inherently low and bounded; a hard limit only adds
  complexity and false 429s
kiss_default: rely on upstream quotas and sensible client backoff first; add an
  explicit limiter at the edge or for expensive outbound calls when bursts hurt
quality: [availability, cost-control]
related: [bulkhead, retry-and-timeouts]
---

# Rate Limiting

**Problem.** Bursts (legitimate or abusive) can exhaust capacity, inflate cost,
or trip partner quotas; unbounded clients then amplify the damage with retries.

**Options.**
- No app-level limit — acceptable for low, trusted traffic.
- Token/leaky bucket at the edge or per client key — usual default.
- Adaptive throttling under measured overload — for mature platforms.

**Failure modes.** Global limit that punishes all tenants for one noisy client;
limit without clear 429 + Retry-After contract; retries that ignore the limit
and create a feedback loop.
