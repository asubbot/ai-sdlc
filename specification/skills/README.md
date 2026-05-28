# specification/skills — Agent instructions

Agent instructions for the SDLC pipeline. **One numbered skill per pipeline stage (1–11)**, plus **optional cross-cutting** skills (C4 C3 diagram, project comparison, threat model from code). Paths in skills use **ai-sdlc-artefacts/** (project root: `scope.md`, `strategy.md`, optional `analytics/`) and **ai-sdlc-artefacts/epics/<epic-id>/** for epic artefacts: ep-scope, ep-context, ep-requirements, ep-acceptance-criteria, ep-system-design, ep-system-design-review, ep-implementation-plan, **ep-code-review** (when saved; **§2.2** uses per-iteration sections), ep-audit-report. Story-level paths are not used by the pipeline.

**Common behaviour:** The agent works in cooperation with the user. When several valid choices exist (output format, path, scope, or interpretation of the request), present them (e.g. A / B) and ask the user which they prefer. Do not proceed until the user has chosen, except when [pipeline.spec.md](../pipeline.spec.md) §4.5 autonomous mode is explicitly enabled and the orchestrator acts as the approver on the user's behalf.

**Mandatory delegation:** Stages **7** (system design review) and **10** (code review) MUST run via a **subagent** or equivalent fresh session—see [pipeline.spec.md](../pipeline.spec.md) §3. Each **§2.1** or **§2.2** iteration after material edits requires a **new** delegated review run.

**Subagent orchestration:** All stages (3–11) SHOULD run in separate subagent sessions for fresh context—see [pipeline.spec.md](../pipeline.spec.md) §4. Stage 9 supports per-task subagent isolation (§4.4).

---

## Common token-optimized context rules

Use the three lightweight context layers defined in [pipeline.spec.md](../pipeline.spec.md):

1. **YAML front matter** — cheap machine-readable file metadata.
2. **Current Gate Summary** — current state of saved review gates.
3. **`ep-context.md`** — compact semantic context for an epic.

### YAML front matter

Pipeline artefacts SHOULD begin with YAML front matter. If front matter is absent, agents MUST fall back to reading the Markdown body.

Source artefact example:

```yaml
---
artefact: ep-scope
epic_id: EP-XXX
status: draft
source_of_truth: true
updated_at: YYYY-MM-DD
---
```

Compact context example:

```yaml
---
artefact: ep-context
epic_id: EP-XXX
status: draft
source_of_truth: false
updated_at: YYYY-MM-DD
---
```

Allowed `status` values: `draft`, `approved`, `superseded`, `archived`.

Full source artefacts such as `ep-scope.md`, `ep-requirements.md`, `ep-acceptance-criteria.md`, and `ep-system-design.md` use `source_of_truth: true`. `ep-context.md` uses `source_of_truth: false`. YAML front matter MUST be updated in the same save operation as the body content it reflects.

### `ep-context.md`

`ep-context.md` is a compact epic handoff file, not a source of truth. If it conflicts with a source artefact, the source artefact wins. If `ep-context.md` is missing for an epic, agents SHOULD create it on the first approved epic artefact write/update during stages 3–11.

Recommended sections:

```markdown
# Epic Context — EP-XXX

## Purpose
## Current Scope
## Key Requirements
## Acceptance Signals
## Design Decisions
## Interfaces / Contracts
## Current Gate Summary
## Open Questions
## Links
```

Keep it short, ideally 100–150 lines. If `ep-context.md` has `updated_at` older than any source artefact it summarizes, treat it as stale and open the changed source artefact before relying on the context.

### Current Gate Summary

Saved review artefacts (`ep-system-design-review.md`, `ep-code-review.md`) SHOULD keep this block above iteration history:

```markdown
## Current Gate Summary

Gate: Pass / Fail / Cap
Latest iteration: N
Last updated: YYYY-MM-DD
Open counts: Blocker X | Major X | Medium X | Minor X
Non-blocking counts: Nit X | Suggestion X
Open findings:
- F-001 Major: One-line finding title.
Next action: Proceed to stage X / Return to stage Y / Operator decision required.
```

The latest full `## Review iteration N` remains the source of truth. Stages 7 and 10 MUST update the latest review iteration and Current Gate Summary atomically when saving review files.

---

## Subagent execution protocol

Each skill can be invoked by an orchestrator as a standalone subagent ([pipeline.spec.md](../pipeline.spec.md) §4). When running as a subagent:

1. **Input brief:** The orchestrator provides: epic ID, stage number, paths to input artefacts, and optional prior-stage summary or review feedback.
2. **Context loading:** Read `ep-context.md` first (when present and current), then YAML front matter of required inputs, then full artefacts only as needed — per the token-optimized context rules above.
3. **Self-contained:** Do not assume access to prior chat history. All context must come from artefacts and the orchestrator brief.
4. **Artefact write:** Write the output artefact per skill rules. In autonomous mode the orchestrator acts as the approver — the "never write until approved" rule is satisfied by the orchestrator's delegation.
5. **Output signal:** On completion, output a structured one-line signal for the orchestrator: `STAGE_<N>_COMPLETE: <artefact_path> [<key_change_summary>]`. For review stages include gate status. For stage 9 tasks: `TASK_COMPLETE: <task_id> [<files_changed>]`.
6. **No ep-context.md writes:** Subagents (except stage 3 initial creation) do not write `ep-context.md` directly. Report material changes in the output signal; the orchestrator applies them.

Each stage skill below includes an **Orchestrator brief** section that specifies the required input, gate preconditions, output signal format, and post-stage validation commands.

---

## Stage → skill (pipeline 1–11)

| Stage | Name | Skill file |
|-------|------|------------|
| 1 | Scope analysis | [01-scope-analysis.skill.md](01-scope-analysis.skill.md) |
| 2 | Strategy analysis | [02-strategy-analysis.skill.md](02-strategy-analysis.skill.md) |
| 3 | Epic planning | [03-epic-planning.skill.md](03-epic-planning.skill.md) |
| 4 | Requirements | [04-requirements.skill.md](04-requirements.skill.md) |
| 5 | Acceptance criteria | [05-acceptance-criteria.skill.md](05-acceptance-criteria.skill.md) |
| 6 | System design | [06-system-design.skill.md](06-system-design.skill.md) |
| 7 | System design review | [07-system-design-review.skill.md](07-system-design-review.skill.md) |
| 8 | Implementation planning | [08-implementation-planning.skill.md](08-implementation-planning.skill.md) |
| 9 | Task execution | [09-task-execution.skill.md](09-task-execution.skill.md) |
| 10 | Code review | [10-code-review.skill.md](10-code-review.skill.md) |
| 11 | Audit | [11-audit.skill.md](11-audit.skill.md) |

**Execution order**

- **Full pipeline (stages 1–11):** 1 → 2 → 3 → … → 11 — see flowchart in [pipeline.spec.md](../pipeline.spec.md) §1.
- **Epic elaboration only:** 3 → 4 → 5 → 6 → 7 → 8 (Epic planning → Requirements → Acceptance criteria → System design → System design review → Implementation planning).
- **After plan approval:** 9 (task execution) → 10 (code review) → 11 (audit).

---

## Cross-cutting skills (not a numbered stage)

| Skill | Use when |
|-------|----------|
| [ep-C4-component.skill.md](ep-C4-component.skill.md) | C4 **C3** Go component diagram for `ep-system-design.md` (optional; complements mandatory C2 container) |
| [project-comparison-report.skill.md](project-comparison-report.skill.md) | Compare an external repo with the consumer product; analytics report under `ai-sdlc-artefacts/analytics/` |
| [user-documentation.skill.md](user-documentation.skill.md) | End-user / operator docs under `docs/` and root `README.md` (installation, config, Docker, operations) |
| [threat-model-report.skill.md](threat-model-report.skill.md) | Code-grounded threat model report (default `docs/threat-model.md` or `ai-sdlc-artefacts/analytics/...`) |

---

## Intent / trigger → skill

When a user request matches an intent below, use the corresponding skill.

| Intent / trigger | Skill |
|------------------|--------|
| Plan or refine **project** scope (`scope.md`) | [01-scope-analysis.skill.md](01-scope-analysis.skill.md) |
| **Epic** planning (`ep-scope.md` per epic) | [03-epic-planning.skill.md](03-epic-planning.skill.md) |
| Define delivery or test **strategy** (`strategy.md`) | [02-strategy-analysis.skill.md](02-strategy-analysis.skill.md) |
| Write or update **requirements** (EARS, `ep-requirements.md`) | [04-requirements.skill.md](04-requirements.skill.md) |
| Write or update **acceptance criteria** | [05-acceptance-criteria.skill.md](05-acceptance-criteria.skill.md) |
| **System design** / architecture / C2 container (`ep-system-design.md`) | [06-system-design.skill.md](06-system-design.skill.md) |
| **C4 C3** component diagram (Go packages) for epic system design | [ep-C4-component.skill.md](ep-C4-component.skill.md) |
| **System design review** / architecture review / requirement traceability | [07-system-design-review.skill.md](07-system-design-review.skill.md) |
| **Implementation plan** (tasks, ordering, verification) | [08-implementation-planning.skill.md](08-implementation-planning.skill.md) |
| **Execute** implementation plan (code tasks, commits) | [09-task-execution.skill.md](09-task-execution.skill.md) |
| **Code review** / PR review | [10-code-review.skill.md](10-code-review.skill.md) |
| **Audit** / quality gate / status report (epic or project) | [11-audit.skill.md](11-audit.skill.md) |
| **Analyse or compare** an external project with PA | [project-comparison-report.skill.md](project-comparison-report.skill.md) |
| **User docs** / operator guide / refresh **README** for deployers | [user-documentation.skill.md](user-documentation.skill.md) |
| **Threat model** / STRIDE / attack surface **from source code** | [threat-model-report.skill.md](threat-model-report.skill.md) |

---

## All skill files in this folder

Numbered pipeline stages: `01-scope-analysis` through `06-system-design`, then `07-system-design-review`, `08-implementation-planning`, `09-task-execution`, `10-code-review`, `11-audit`. Cross-cutting: `ep-C4-component.skill.md`, `project-comparison-report.skill.md`, `user-documentation.skill.md`, `threat-model-report.skill.md`.

**Single source of stage I/O:** [pipeline.spec.md](../pipeline.spec.md) §2 (table of stages, inputs, outputs).
