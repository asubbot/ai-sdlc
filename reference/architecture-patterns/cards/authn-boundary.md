---
# OKF v0.1
type: Pattern
title: Authentication Boundary
description: Concentrate authentication and authorization at an explicit trust
  boundary (gatekeeper/edge), keeping inner components free of ad-hoc checks.
timestamp: 2026-07-24
tags: [security, authentication, trust-boundary]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/gatekeeper
    note: Primary upstream (Gatekeeper pattern, Azure)
  - url: https://learn.microsoft.com/en-us/azure/architecture/patterns/federated-identity
    note: Federated Identity pattern (Azure)
forces: [single enforcement point vs defense in depth, least privilege, auditability]
when: multiple entry points or handlers would otherwise each implement their own
  authentication/authorization checks
when_not: a single-user local tool with no untrusted input surface — an auth
  layer adds ceremony without a real trust boundary
kiss_default: one middleware/interceptor at the single entry point with an
  explicit allowlist; delegate identity to an existing provider instead of
  storing credentials
quality: [security, auditability]
related: [module-boundaries]
---

# Authentication Boundary

**Problem.** Scattered per-handler auth checks drift apart; one forgotten check
becomes an unauthenticated endpoint, and audits cannot say where trust starts.

**Options.**
- Edge middleware/gatekeeper: authenticate and authorize once, pass verified
  identity inward; inner code trusts the boundary.
- Federated identity: delegate authentication to an external provider (tokens),
  keep no credentials in the product.
- Defense in depth: boundary plus explicit checks around the most sensitive
  operations only.

**Failure modes.** Bypass paths that skip the boundary (debug ports, admin CLIs);
authorization decisions duplicated inside after "temporary" exceptions; trusting
unvalidated identity claims forwarded across the boundary.
