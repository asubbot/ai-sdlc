---
name: implementation-planning.skill
description: Produce one implementation plan per epic (stage 8); output ep-implementation-plan.md. Use when planning tasks for an epic, e.g. "implementation plan for this epic", "break down epic into tasks".
---

# Stage 8: Implementation planning

**Pipeline:** [pipeline.spec.md](../pipeline.spec.md)  
**Output:** ai-sdlc-artefacts/epics/<epic-id>/ep-implementation-plan.md; may update ep-context.md

---

## Core Principles

Follow these principles for all implementation planning work:

1. **Never write until approved** — Do not create or overwrite ep-implementation-plan.md until the user explicitly approves the draft (e.g. "lgtm", "save", "approve"). All edits go into the draft in chat; do not write to file until approval.
2. **Existing file is baseline** — If ep-implementation-plan.md already exists for the epic, treat it as the current baseline; propose changes as edits and overwrite only after user approval.
3. **Options when in doubt** — When multiple valid choices exist (e.g. task granularity, ordering, checkpoint placement), present options (e.g. A/B) and ask the user to choose before proceeding.
4. **References** — Links only to paths under `ai-sdlc-artefacts/`; every linked document must exist. Write in English.
5. **Practical and short** — Get to the point. Be practical above all. Be short and specific. Only include tasks that can be performed by a coding agent. Each task must have a verification criterion.
6. **Legacy** — Do not modify content under `legacy` folders; use as reference only.
7. **Token-optimized context** — Read YAML front matter, Current Gate Summary, and current ep-context.md first. Open full source artefacts only for traceability, stale context, missing details, or gate disputes.

---

## 1. Context and goal

You are the Tech Lead for this epic. Your role is to produce the implementation plan per epic (stage 8).

**Goal:** Produce ep-implementation-plan.md: ordered tasks for the epic with dependencies, verification per task, traceability to REQ and AC (from ep-requirements and ep-acceptance-criteria), and checkpoints. Optionally group tasks by theme or label with user-story-like IDs (e.g. US-01) if useful; no separate story-level artefact is required.

**Inputs:** ep-context.md when present and current; ep-scope.md, ep-requirements.md, ep-acceptance-criteria.md, ep-system-design.md (all under ai-sdlc-artefacts/epics/<epic-id>/); and test strategy (e.g. strategy.md under ai-sdlc-artefacts/). **Recommended:** ep-system-design-review.md after stage 7 (system design review), starting with YAML front matter and Current Gate Summary; read full `## Review iteration N` sections only when the summary reports open findings, repeated issues, or a traceability dispute. **Prerequisite:** Stages **6↔7** must meet **exit** in [pipeline.spec.md](../pipeline.spec.md) **§2.1** (zero Blocker/Major/Medium/Minor) or have an explicit **operator decision** after the iteration cap—do not start stage 8 before that. If any required epic input is missing, ask the user to run the corresponding prior stage first.

**Questions to answer:** What are the discrete coding steps? In what order do we execute them? Where do we place checkpoints and how do we verify each step?

---

## 2. Implementation planning workflow

Follow this order:

1. **Check inputs** — Ensure ep-scope.md, ep-requirements.md, ep-acceptance-criteria.md, and ep-system-design.md exist for the epic. If not, ask the user to run the missing stage(s) first. If ep-context.md exists and is current, read it first for orientation. If it is missing or stale, open the changed source artefacts and refresh context after approval. If ep-system-design-review.md is missing, recommend completing stage 7 first; proceed only if the user explicitly chooses to skip the review. Optionally use test strategy (e.g. strategy.md) as reference.
2. **Check existing ep-implementation-plan** — If ep-implementation-plan.md exists for the epic, treat it as the baseline; propose changes as edits.
3. **Draft in chat** — Draft the implementation plan in chat (task list, verification, checkpoints). Show it to the user and ask for clarification or changes. Do not write to file yet.
4. **Resolve choices** — When multiple valid options exist (e.g. task breakdown, ordering), present options (e.g. A/B) and ask the user to choose.
5. **Write after approval** — Create or update ai-sdlc-artefacts/epics/<epic-id>/ep-implementation-plan.md only when the user explicitly approves (e.g. "lgtm", "save", "approve").
6. **Refresh epic context** — Optionally update ep-context.md with plan-level execution notes and Links when useful; keep it compact.
7. **Legacy** — Do not modify content under `legacy` folders; use as reference only.

---

## 3. Output structure (ep-implementation-plan.md)

Use these elements (or user-agreed equivalents):

- **YAML front matter** — Start with artefact metadata:
  ```yaml
  ---
  artefact: ep-implementation-plan
  epic_id: EP-XXX
  status: draft
  source_of_truth: true
  updated_at: YYYY-MM-DD
  ---
  ```

- **Document header** — Purpose, pipeline link, Previous/Related links to ep-acceptance-criteria.md, ep-requirements.md, ep-system-design.md, ep-system-design-review.md when present (and optionally strategy.md).
- **Task list** — Numbered tasks with dependencies (e.g. "Task 2 depends on Task 1"). Clear objective per task; sub-bullets for details. Group into sections by theme if helpful.
- **Verification per task** — For each task, state how to confirm the step is done (e.g. test passes, build succeeds, review done).
- **Checkpoints** — Explicit checkpoints (e.g. "Ensure all tests pass before proceeding"; "Ask the user if questions arise").
- **Traceability** — Each task that implements scope MUST reference REQ and AC where applicable (links to ep-requirements.md and ep-acceptance-criteria.md). Use "—" when a task has no direct link (e.g. checkpoints, tooling). Optional: add _User Stories:_ or similar labels per task if the epic uses story IDs for grouping.

**Task format:**

- One clear objective per task.
- Sub-bullets for technical or procedural details.
- Traceability: **_Requirements:_** REQ-EE.NNN (link to ep-requirements); **_Acceptance Criteria:_** AC-EE.NNN (link to ep-acceptance-criteria; anchor form e.g. #ac-01-001). Use "—" when the task is a checkpoint or has no direct REQ/AC.

**Quality:** Tasks must be actionable by a coding agent. Each task has a verification criterion. Keep the plan short and specific.

**Example** (one task):

```markdown
- [ ] 1.1 Implement config load and validation
  - Define config struct; load JSON from path; validate required fields.
  - _Requirements:_ [REQ-01.003](ep-requirements.md#nodes-and-ssh), [REQ-01.004](ep-requirements.md#nodes-and-ssh)
  - _Acceptance Criteria:_ [AC-01.005](ep-acceptance-criteria.md#ac-01-005)
  - **Verification:** `go build ./...` passes; invalid config returns error.
```

---

## 4. Done when

Verify all before considering the stage complete:

- [ ] ep-implementation-plan.md exists at ai-sdlc-artefacts/epics/<epic-id>/ep-implementation-plan.md
- [ ] YAML front matter is present and consistent with the saved implementation plan
- [ ] Document contains header with links to epic artefacts, task list (numbered, with dependencies), verification per task, and checkpoints
- [ ] Each task that implements scope has traceability to REQ/AC (or "—" for checkpoints)
- [ ] Every link in the document points to an existing path under `ai-sdlc-artefacts/` (no broken links)
- [ ] The user has explicitly approved the plan
