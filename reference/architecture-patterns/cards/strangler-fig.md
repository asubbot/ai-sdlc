---
# OKF v0.1
type: Pattern
title: Strangler Fig
description: Incrementally replace a legacy system by routing slices of traffic
  to a new implementation until the old one can be retired.
timestamp: 2026-07-27
tags: [evolution, migration, legacy]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/strangler-fig
    note: Primary upstream (Strangler Fig pattern, Azure)
forces: [risk of big-bang rewrite, dual-run complexity, routing and parity]
when: a legacy component must be replaced while continuing to serve production
  traffic, and a full cutover is too risky
when_not: the system is small enough that a short, well-tested cutover is safer
  than prolonged dual maintenance
kiss_default: prefer a direct replacement with feature flags when scope is small;
  use strangler routing when the legacy surface is large or poorly understood
quality: [risk-reduction, maintainability]
related: [module-boundaries]
---

# Strangler Fig

**Problem.** Rewriting a large legacy surface in one step risks a long dark
period and a painful cutover; leaving it forever blocks improvement.

**Options.**
- Big-bang rewrite and cutover — only for small, well-tested scopes.
- Strangler fig: facade/router sends selected slices to the new system over time.
- Freeze legacy and build alongside without traffic migration — rarely finishes.

**Failure modes.** Facade that becomes permanent dual complexity; no exit
criteria for retiring legacy; parity gaps discovered only in production.
