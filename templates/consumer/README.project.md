# Project name

Short description of the product (replace during stage 1 — scope analysis).

## Agentic SDLC

This project uses [ai-sdlc](https://github.com/asubbot/ai-sdlc) (nested under `ai-sdlc/`). Process pin: [`ai-sdlc.version`](ai-sdlc.version).

- **Artefacts:** [`ai-sdlc-artefacts/`](ai-sdlc-artefacts/)
- **Bootstrap / new project:** see [ai-sdlc/README.md](ai-sdlc/README.md) — *Starting a new project*

## Commands

```bash
make build      # build bin/validate from ai-sdlc/tools/validate
make validate   # project gate: AC (all epics) + pipeline/structure per epic
make check      # go vet + bootstrap smoke tests in tests/
```
