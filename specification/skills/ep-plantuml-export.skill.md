---
name: ep-plantuml-export.skill
description: >-
  Render C4 PlantUML sources under an epic diagrams/ folder to PNG and verify
  markdown embed paths. Use after writing or updating *.puml in stages 4, 6, or
  ep-C4-component; when PNG is missing; or when the user asks to regenerate
  architecture diagrams.
---

# Epic PlantUML export (C4 PNG)

**Pipeline:** Cross-cutting helper for stages **4** ([04-requirements.skill.md](04-requirements.skill.md)), **6** ([06-system-design.skill.md](06-system-design.skill.md)), and optional [ep-C4-component.skill.md](ep-C4-component.skill.md).

**Goal:** After saving or updating a `.puml` file, **automatically render** the matching `.png` and ensure epic markdown embeds the PNG (not source-only).

---

## Prerequisites

- **`plantuml` on PATH** (Homebrew, jar, or IDE CLI). Verify: `plantuml -version`.
- Network access for C4-PlantUML `!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/...` on first render (cached locally afterward).

If `plantuml` is missing: stop and tell the operator to install it; do **not** leave `.puml`-only artefacts when the tool is required by the stage skill.

---

## Naming conventions

| Diagram | Source | PNG output | Used in |
|---------|--------|------------|---------|
| C4 C1 (context) | `diagrams/c4-context.puml` | `diagrams/c4-context.png` | `ep-requirements.md` |
| C4 C2 (containers) | `diagrams/c4-container.puml` | `diagrams/c4-container.png` | `ep-system-design.md` |
| C4 C3 (Go components) | `diagrams/c4-component-go.puml` | `diagrams/c4-component-go.png` | `ep-system-design.md` |

**Stable filenames:** First line of each file MUST be `@startuml <stem>` where `<stem>` equals the basename without extension (e.g. `@startuml c4-context`). PlantUML writes `<stem>.png` next to the source when run from `diagrams/`.

Do **not** use spaces or epic ids in `@startuml` names (avoids `EP034_C4_Context.png` mismatches).

---

## Workflow (mandatory after .puml write/update)

1. **Working directory** — `cd ai-sdlc-artefacts/epics/<epic-id>/` (epic root, not repo root).
2. **Render** — For each changed `.puml` under `diagrams/`:
   ```bash
   plantuml -tpng diagrams/<stem>.puml
   ```
   Or from `diagrams/`:
   ```bash
   cd diagrams && plantuml -tpng <stem>.puml && cd ..
   ```
3. **Verify** — After each render:
   - `plantuml` exit code is `0`.
   - Expected PNG exists (e.g. `test -f diagrams/c4-context.png`).
   - PNG is newer than or equal to the `.puml` mtime.
4. **Embed in markdown** — In the target artefact, include **before** the Source line:
   ```markdown
   <p align="center"><img src="diagrams/c4-context.png" alt="C4 C1 — System Context" /></p>
   ```
   Adjust `alt` text per epic/diagram. Keep the Source line:
   ```markdown
   **Source:** [diagrams/c4-context.puml](diagrams/c4-context.puml). Regenerate PNG: `plantuml -tpng diagrams/c4-context.puml` from this directory.
   ```
5. **Report** — On success: `PLANTUML_EXPORT_COMPLETE: ai-sdlc-artefacts/epics/<epic-id>/diagrams [<list of PNG paths>]`.

---

## Batch render (all epic diagrams)

From epic root:

```bash
for f in diagrams/*.puml; do plantuml -tpng "$f"; done
```

Use when backfilling PNGs for an existing epic or after copying `.puml` from another epic.

---

## Failure handling

| Condition | Action |
|-----------|--------|
| `plantuml` not found | Stop; document in chat; do not mark stage Done when PNG is required |
| Non-zero exit / parse error | Fix `.puml`; do not commit broken PNG; report stderr |
| PNG name mismatch | Fix `@startuml` stem or rename; re-run render |
| Missing embed in markdown | Add `<img>` block per § Workflow step 4 |

---

## Done when

- [ ] Every `.puml` listed by the calling stage has a matching `.png` in `diagrams/`.
- [ ] Referencing markdown files embed the PNG with a working relative path.
- [ ] `@startuml` stem matches the documented PNG basename (`c4-context`, `c4-container`, `c4-component-go`).
- [ ] Regeneration command in the Source line matches the paths above.

---

**Called from:** stage 4 (after C1 `.puml`), stage 6 (after C2 `.puml`), [ep-C4-component.skill.md](ep-C4-component.skill.md) (after C3 `.puml`).
