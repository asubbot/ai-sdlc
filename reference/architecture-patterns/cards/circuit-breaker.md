---
# OKF v0.1
type: Pattern
title: Circuit Breaker
description: Stop calling a dependency that is failing, fail fast while it is
  open, and probe periodically before closing again.
timestamp: 2026-07-24
tags: [resilience, fault-isolation]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
    note: Primary upstream (Circuit Breaker pattern, Azure)
  - url: https://martinfowler.com/bliki/CircuitBreaker.html
    note: Fowler, original description
forces: [failing-fast vs retry persistence, recovery detection, shared dependency protection]
when: a dependency can fail for long periods and continued calls waste resources,
  block threads, or delay recovery
when_not: faults are short transients already handled by retry with backoff, or
  call volume is so low that a breaker never accumulates a meaningful signal
kiss_default: timeout + bounded retry first; add a breaker only when a dependency
  outage measurably degrades the caller (thread/connection exhaustion, latency)
quality: [availability, resilience]
related: [retry-and-timeouts]
---

# Circuit Breaker

**Problem.** When a dependency is down, every call still burns a timeout;
under load this exhausts workers and prolongs the outage on both sides.

**Options.**
- Retry with backoff only — fine for low volume and short outages.
- Circuit breaker (closed → open → half-open) around the dependency client.
- Breaker + fallback (cached value, degraded mode) when a default response exists.

**Failure modes.** Thresholds guessed and never tuned (breaker flaps or never
trips); one breaker shared across unrelated endpoints (one bad route opens all);
no telemetry on state changes, so open circuits go unnoticed.
