---
# OKF v0.1
type: Pattern
title: Bulkhead
description: Partition resources (threads, connections, processes) so a failure
  or overload in one dependency cannot exhaust the whole system.
timestamp: 2026-07-27
tags: [resilience, isolation, capacity]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/bulkhead
    note: Primary upstream (Bulkhead pattern, Azure)
forces: [failure isolation, capacity fairness, operational complexity of pools]
when: multiple dependencies or workloads compete for shared workers, sockets,
  or memory and one slow path can starve the rest
when_not: a single dependency or very low concurrency where separate pools add
  tuning cost without measurable isolation benefit
kiss_default: one shared pool with timeouts first; introduce per-dependency
  limits only when a noisy neighbour is observed in production or load tests
quality: [availability, resilience]
related: [circuit-breaker, retry-and-timeouts, rate-limiting]
---

# Bulkhead

**Problem.** Shared worker or connection pools mean one slow or failing
dependency can exhaust capacity and take unrelated paths down with it.

**Options.**
- Shared pool + timeouts — simplest; acceptable until contention is proven.
- Bulkheads: separate pools/quotas per dependency or workload class.
- Process-level isolation (separate services/containers) — strongest, highest ops cost.

**Failure modes.** Over-partitioning (idle capacity while other pools starve);
limits guessed and never revisited; bulkheads without metrics hide which pool is full.
