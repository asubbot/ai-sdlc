# Consumer project templates

Files in this directory are **copied to the product repository root** during greenfield bootstrap ([00-project-bootstrap.skill.md](../../specification/skills/00-project-bootstrap.skill.md)).

| Template file | Product path |
|---------------|----------------|
| `AGENTS.md` | `AGENTS.md` |
| `.gitignore` | `.gitignore` |
| `Makefile` | `Makefile` |
| `ai-sdlc.version` | `ai-sdlc.version` |
| `README.project.md` | `README.md` |
| `ai-sdlc-artefacts/**` | `ai-sdlc-artefacts/**` |
| `.github/workflows/ai-sdlc.yml` | `.github/workflows/ai-sdlc.yml` |

## Layout

- **Default (layout A):** `ai-sdlc/` is a local clone (listed in `.gitignore`). Process version is recorded in `ai-sdlc.version` at the product root.
- **Layout B:** `ai-sdlc/` is a git submodule; remove the `ai-sdlc/` line from `.gitignore`.

## Rules

1. Do **not** edit normative process files inside the vendored `ai-sdlc/` copy for durable changes — open PRs in the canonical ai-sdlc repository and bump `ai-sdlc.version` in this project.
2. Pipeline outputs live under `ai-sdlc-artefacts/` at the product root, not inside `ai-sdlc/`.
3. Open the IDE workspace at the **product root**, not at `ai-sdlc/`.

See [ai-sdlc README](../../README.md) (section **Starting a new project**) for the full walkthrough.
