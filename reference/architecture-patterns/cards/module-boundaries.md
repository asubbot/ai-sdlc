---
# OKF v0.1
type: Pattern
title: Module Boundaries
description: Split a codebase into cohesive modules with explicit, one-way
  dependencies so change stays local and wiring stays at the edge.
timestamp: 2026-07-24
tags: [modularity, boundaries, dependencies]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/microservices/model/domain-analysis
    note: Primary upstream (domain analysis and bounded contexts, Microsoft Learn)
  - url: https://martinfowler.com/bliki/BoundedContext.html
    note: Fowler, bounded context definition
forces: [cohesion vs coupling, change locality, testability in isolation, dependency direction]
when: the codebase has grown past one obvious module and unrelated concerns start
  importing each other
when_not: a small single-purpose tool where extra layers add indirection without
  isolating any real change axis
kiss_default: start with a flat package per concern and one-way imports; introduce
  formal boundaries (interfaces, dependency rules) only when a cycle or shared
  mutable state appears
quality: [maintainability, testability]
related: [authn-boundary]
---

# Module Boundaries

**Problem.** Without explicit boundaries, every package can import every other;
changes ripple, tests need the whole world, and cycles appear.

**Options.**
- Flat packages by concern with a one-way import rule — cheapest, often enough.
- Explicit boundaries: public interfaces per module, wiring only in a composition
  root (e.g. `cmd/`), dependency direction checked in review or by a script.
- Full bounded contexts with separate data models — for genuinely distinct domains.

**Failure modes.** Boundaries drawn by technical layer instead of domain (churn
crosses every module); anemic "utils" modules that everything imports; boundary
enforced only by convention and silently eroded.
