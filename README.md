# ai-sdlc — canonical process repository

This repository is the canonical source of truth for the shared **agentic SDLC process** used by multiple product repositories.

`ai-sdlc` describes the process only; project execution outputs are stored in each product repository under `ai-sdlc-artefacts/`.

---

## Repository layout

| Path | Role |
|------|------|
| **[specification/pipeline.spec.md](specification/pipeline.spec.md)** | **Normative process:** stage order, inputs/outputs, stage→skill mapping, **HOTL by default** (HITL at decision points), delegated execution, artefact naming, traceability, **agent execution expectations** (single process, AC validation, etc.). |
| **[specification/skills/](specification/skills/)** | **Per-stage agent instructions** (`01-` … `11-` plus optional cross-cutting skills). Each skill defines workflow and artefact structure for its stage. |
| **[specification/skills/README.md](specification/skills/README.md)** | Index and **common behaviour** across skills. |
| **[proposals.md](proposals.md)** | Proposed pipeline improvements and measurement notes; a living document, not the normative process. |
| **[tools/validate/](tools/validate/)** | AC↔test coverage checker (`./tools/validate/validate` from repo root); see [VALIDATION.md](tools/validate/VALIDATION.md) and [README](tools/validate/README.md). |

---

## Consumption model (workspace-only)

Projects consume this repository as a separate workspace root. To keep reproducibility:

1. Each consumer project stores an `ai-sdlc.version` file with a pinned tag or commit SHA.
2. CI in consumer projects verifies that the pinned revision exists in this canonical repository.
3. Process changes are made here via PR; product repositories update only their pinned version.

### Pin file example (`ai-sdlc.version`)

```text
v1.0.1
```

Or pin an exact commit:

```text
f0264d2063e2d881250e0db0542f0de3dfd9c413
```

### Consumer CI example (GitHub Actions)

Verify the pin exists in the canonical repository (set `AI_SDLC_REPO` to your fork or upstream, e.g. `ORG/ai-sdlc`):

```yaml
- name: Verify ai-sdlc pin
  env:
    AI_SDLC_REPO: ORG/ai-sdlc
  run: |
    PIN="$(tr -d '[:space:]' < ai-sdlc.version)"
    if git ls-remote "https://github.com/${AI_SDLC_REPO}.git" "refs/tags/${PIN}" | grep -q .; then
      echo "Tag ${PIN} found"
    elif git ls-remote "https://github.com/${AI_SDLC_REPO}.git" "${PIN}" | grep -q .; then
      echo "Commit ${PIN} found"
    else
      echo "Pin ${PIN} not found in ${AI_SDLC_REPO}" >&2
      exit 1
    fi
```
