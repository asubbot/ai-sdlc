---
# OKF v0.1
type: Pattern
title: Caching
description: Serve hot reads from a cache (cache-aside or read-through) to cut
  latency and load, with explicit staleness and invalidation policy.
timestamp: 2026-07-24
tags: [performance, caching, staleness]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside
    note: Primary upstream (Cache-Aside pattern, Azure)
  - url: https://learn.microsoft.com/en-us/azure/architecture/best-practices/caching
    note: Caching guidance (Microsoft Learn)
forces: [read latency vs staleness, invalidation correctness, memory cost, stampede on miss]
when: reads are hot, the source is comparatively slow or expensive, and bounded
  staleness is acceptable
when_not: data must always be current (auth decisions, balances) or read volume
  is too low for hit rate to matter
kiss_default: no cache; if needed, in-process TTL cache for one hot spot before
  any shared cache infrastructure
quality: [performance, scalability]
related: [module-boundaries]
---

# Caching

**Problem.** Hot reads hit a slow or expensive source on every request; naive
caching then trades correctness for speed via stale or inconsistent data.

**Options.**
- No cache — correct by default; measure before optimizing.
- In-process TTL cache per hot spot — cheap, node-local staleness only.
- Cache-aside with shared cache — scales across instances; needs invalidation
  policy and stampede protection.

**Failure modes.** Invalidation forgotten on the write path (permanently stale);
cache stampede when a hot key expires; treating the cache as source of truth;
unbounded memory without eviction.
