# Contributing

## Scope

This repository defines the shared SDLC process for multiple projects. Keep changes process-focused and repository-agnostic unless explicitly approved.

## Workflow

1. Create a feature branch from `main`.
2. Update normative specification and related skills together.
3. Update docs (`README.md`, `CHANGELOG.md`) when behavior changes.
4. Open a PR with rationale and migration notes for consumer repositories.

## Release and adoption

1. Tag a release after approved changes (for example `v1.0.1`).
2. Consumer repositories update their `ai-sdlc.version` pin via a dedicated PR.
3. Do not make direct process edits in consumer repository copies after cutover.

See [README.md](README.md) (Consumption model) for `ai-sdlc.version` format and a minimal CI verification snippet.
