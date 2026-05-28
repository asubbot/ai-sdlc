# Project name

Short description of the product (replace during stage 1 — scope analysis).

## Agentic SDLC

This project uses [ai-sdlc](https://github.com/asubbot/ai-sdlc) (nested under `ai-sdlc/`). Process pin: [`ai-sdlc.version`](ai-sdlc.version).

- **Artefacts:** [`ai-sdlc-artefacts/`](ai-sdlc-artefacts/)
- **Bootstrap / new project:** see [ai-sdlc/README.md](ai-sdlc/README.md) — *Starting a new project*

## Commands

Product gates (run locally and in `.github/workflows/ci.yml`):

```bash
make build      # build bin/validate from ai-sdlc/tools/validate (requires ai-sdlc/ at pin)
make check      # fmt, vet, govulncheck, golangci-lint, race tests + coverage on tests/
make validate   # artefact gate: AC (all epics) + pipeline/structure per epic
```

First `make check` may download linter and vulnerability DB modules. Uncomment additional targets in the `Makefile` when `cmd/`, integration, or e2e tests exist.
