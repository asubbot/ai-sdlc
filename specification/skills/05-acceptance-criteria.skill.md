---
name: acceptance-criteria.skill
description: Produce epic acceptance criteria from ep-scope and ep-requirements; output ep-acceptance-criteria.md. Use when defining or refining epic acceptance criteria (stage 5), e.g. "acceptance criteria for this epic", "AC from requirements".
---

# Stage 5: Acceptance criteria

**Pipeline:** [pipeline.spec.md](../pipeline.spec.md)  
**Output:** ai-sdlc-artefacts/epics/<epic-id>/ep-acceptance-criteria.md; update ep-context.md

---

## Orchestrator brief (subagent mode)

When launched as a subagent by the pipeline orchestrator ([pipeline.spec.md](../pipeline.spec.md) §4):

- **Required input:** epic ID, paths to `ep-scope.md` and `ep-requirements.md`
- **Context:** `ep-context.md` (read first for orientation)
- **Gate check before launch:** `ep-scope.md` and `ep-requirements.md` must exist
- **Output signal:** `STAGE_5_COMPLETE: ai-sdlc-artefacts/epics/<epic-id>/ep-acceptance-criteria.md [<N> ACs] [context_delta: <Acceptance Signals + Links summary>]`
- **Validation after:** `./tools/validate/validate req EP-XXX` (REQ↔AC traceability), `./tools/validate/validate structure EP-XXX` (artefact structure, from repo root)
- **ep-context.md:** Subagent does **not** write `ep-context.md`; include compact Acceptance Signals and Links updates in `context_delta` for the orchestrator to apply

---

## Core Principles

Follow these principles for all acceptance criteria work:

1. **HOTL artefact writes** — In pipeline execution, create or overwrite ep-acceptance-criteria.md when ep-scope.md and ep-requirements.md exist, AC traceability is maintained, and no required HITL decision point from [pipeline.spec.md](../pipeline.spec.md) is triggered.
2. **Existing file is baseline** — If ep-acceptance-criteria.md already exists for the epic, treat it as the current baseline; preserve durable acceptance criteria unless a required HITL decision records a change.
3. **Decision points** — Ask the operator only for required HITL decisions (e.g. material acceptance trade-off, missing prerequisite, source-of-truth conflict). For routine AC granularity or format choices, choose the simplest consistent default and record rationale when useful.
4. **References** — Links only to paths under `ai-sdlc-artefacts/`; every linked document must exist. Keep traceability to ep-requirements. Write in English.
5. **Stable IDs only** — Use acceptance criteria IDs in the form **AC-EE.NNN** where EE is the two-digit epic number (e.g. 01 for EP-001, 02 for EP-002) and NNN is the three-digit AC number within that epic (e.g. AC-01.001, AC-02.007). This avoids ID collisions across epics. Do not use internal UUIDs.
6. **Practical and short** — Get to the point. Be practical above all. Be short and specific.
7. **Legacy** — Do not modify content under `legacy` folders; use as reference only.
8. **Token-optimized context** — On save, add YAML front matter to ep-acceptance-criteria.md. **Orchestrated subagent mode:** report Acceptance Signals and Links deltas in `context_delta`; orchestrator updates ep-context.md. **Solo/HOTL without orchestrator:** refresh the Acceptance Signals section in ep-context.md (create from available epic artefacts if missing).

---

## 1. Context and goal

You are an experienced QA / acceptance criteria analyst. Your role is to produce the epic acceptance criteria document (stage 5).

**Goal:** Produce ep-acceptance-criteria.md: testable conditions for the epic in Gherkin (Given/When/Then) or equivalent, with AC ID and traceability to REQ. This output is the input for system design (stage 6), system design review (stage 7), and implementation planning (stage 8).

**Inputs:** ep-scope.md and ep-requirements.md (ai-sdlc-artefacts/epics/<epic-id>/). If either is missing, treat continuing as a required HITL missing-prerequisite decision or run the missing stage when the operator has authorized HOTL execution.

**Questions to answer:** When is this epic "done" from a test perspective? What scenarios (Given/When/Then) cover the requirements? How do AC map to REQ?

---

## 2. Acceptance criteria workflow

Follow this order:

1. **Check inputs** — Ensure ep-scope.md and ep-requirements.md exist for the epic. If not, treat continuing as a required HITL missing-prerequisite decision or run the missing stage when the operator has authorized HOTL execution. If ep-context.md exists and is current, read it first for orientation; if stale or missing, fall back to source artefacts.
2. **Check existing ep-acceptance-criteria** — If ep-acceptance-criteria.md exists for the epic, treat it as the baseline; propose changes as edits.
3. **Draft** — Draft acceptance criteria (section by section or by block). In HOTL mode, proceed to write when the draft is internally consistent and no required HITL decision is open.
4. **Resolve decision points** — Use HOTL defaults for routine choices; stop for operator choice only when a required HITL decision point applies.
5. **Write artefact** — Create or update ai-sdlc-artefacts/epics/<epic-id>/ep-acceptance-criteria.md under HOTL when inputs are sufficient and no required HITL decision is open.
6. **Refresh epic context** — **Orchestrated subagent mode:** report Acceptance Signals and Links deltas in the output signal (`context_delta`); do not edit ep-context.md. **Solo/HOTL without orchestrator:** update ep-context.md Acceptance Signals and Links concisely; do not copy the full AC list.
7. **Legacy** — Do not modify content under `legacy` folders; use as reference only.

---

## 3. Output structure (ep-acceptance-criteria.md)

Use these section headings (or user-agreed equivalents).

Start the document with YAML front matter:

```yaml
---
artefact: ep-acceptance-criteria
epic_id: EP-XXX
status: draft
source_of_truth: true
updated_at: YYYY-MM-DD
---
```

- **Introduction** — Brief summary of the epic and purpose of this document.
- **Acceptance criteria index** — Table: AC ID (link to AC in document), REQ (link to ep-requirements.md section), Summary (one-line). Required.
- **Acceptance criteria** — List of AC: ID (AC-EE.NNN format, e.g. AC-01.001), formulation in Gherkin (Given/When/Then) or equivalent, traceability to REQ (links to ep-requirements.md).

**Gherkin:** Prefer Given/When/Then. Scenario order: happy path first, then negative path, alternative flows, edge cases.

**Traceability:** Every AC must trace to one or more REQ from ep-requirements.md. In the index, the REQ column must use links to the requirement section (e.g. [REQ-01.001](ep-requirements.md#interface-and-deployment)).

**Quality:** Each AC must be testable; one clear scenario per AC; short and specific; no vague wording.

**Example** — One AC in Gherkin with REQ link:

```markdown
**AC-01.001** (Trace: REQ-01.001)
Given the user is logged in
When the user requests the dashboard
Then the system SHALL display the summary widget within 2 seconds
```

---

## 4. Done when

Verify all before considering the stage complete:

- [ ] ep-acceptance-criteria.md exists at ai-sdlc-artefacts/epics/<epic-id>/ep-acceptance-criteria.md
- [ ] YAML front matter is present and consistent with the saved acceptance criteria content
- [ ] ep-context.md was refreshed with compact Acceptance Signals and Links, **or** (orchestrated subagent mode) `context_delta` was reported for the orchestrator to apply
- [ ] Document contains **Introduction** (epic summary and document purpose), **Acceptance criteria index** (AC ID | REQ with links | Summary), and **Acceptance criteria** (AC-EE.NNN with Gherkin or equivalent and traceability to REQ)
- [ ] Every link in the document points to an existing path under `ai-sdlc-artefacts/` (no broken links)
- [ ] Every AC traces to at least one REQ from ep-requirements.md (REQ blocks use `### REQ-EE.NNN —` in ep-requirements.md)
- [ ] `./bin/validate req EP-XXX` exits 0
- [ ] Content was written under HOTL, or any required HITL decision was recorded before writing
