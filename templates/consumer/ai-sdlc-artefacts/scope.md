---
artefact: scope
status: draft
source_of_truth: true
updated_at: 2026-05-28
---

# Project scope (stub)

## Introduction

Greenfield product repository adopting the ai-sdlc agentic SDLC. Product vision and features are **to be defined** with the operator (stage 1). This stub satisfies epic traceability until scope is refined.

## Glossary

- **ai-sdlc** — Canonical agentic SDLC process repository (pinned by `ai-sdlc.version`).
- **Consumer project** — This repository (product code + `ai-sdlc-artefacts/`).

## In scope

- Adopt ai-sdlc process (EP-000): pin, artefacts layout, validator, CI.
- Define and deliver the product using pipeline stages 1–11 after adoption.

## Out of scope / deferred

- Product-specific features until captured in a refined scope (stage 1).
- Custom process forks inside `ai-sdlc/` (use pin + upstream PRs instead).
