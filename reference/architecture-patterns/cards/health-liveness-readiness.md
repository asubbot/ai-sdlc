---
# OKF v0.1
type: Pattern
title: Health, Liveness, and Readiness
description: Expose explicit probes so orchestrators and operators can tell
  whether the process is alive versus ready to receive traffic.
timestamp: 2026-07-27
tags: [operations, deploy, observability]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes
    note: Primary upstream (Kubernetes liveness, readiness, and startup probes)
  - url: https://learn.microsoft.com/en-us/dotnet/architecture/microservices/implement-resilient-applications/monitor-app-health
    note: Health monitoring in microservices (Microsoft Learn)
forces: [safe deploy/restart, avoid traffic to broken instances, probe side effects]
when: the process is managed by an orchestrator, load balancer, or supervisor
  that needs a contract for restart vs traffic removal
when_not: a short-lived CLI or local tool with no process manager that would
  act on the probe
kiss_default: a cheap liveness check (process up) first; add readiness that
  verifies critical dependencies only when bad instances would otherwise receive traffic
quality: [operability, availability]
related: [circuit-breaker, bulkhead]
---

# Health, Liveness, and Readiness

**Problem.** Without a clear alive/ready contract, orchestrators either kill
healthy processes or keep sending traffic to instances that cannot serve it.

**Options.**
- No probes — fine for unmanaged local tools.
- Liveness (restart if dead) + readiness (remove from load until ready).
- Deep dependency checks on every probe — usually too heavy; prefer shallow readiness.

**Failure modes.** Readiness that calls flaky dependencies and flaps; liveness
that does expensive work and self-DoS; probes that mutate state.
