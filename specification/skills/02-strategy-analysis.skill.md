---
name: strategy-analysis.skill
description: Define delivery strategy and test strategy from scope; produce strategy.md as input for epic planning. Use when the user wants to define or refine delivery strategy, test strategy, or when moving from scope (stage 1) to analysis (stage 2).
---

# Stage 2: Strategy analysis

**Pipeline:** [pipeline.spec.md](../pipeline.spec.md)
**Output:** ai-sdlc-artefacts/strategy.md

## Core Principles

Follow these principles for all strategy analysis work:

1. **HOTL artefact writes** — In pipeline execution, create or overwrite strategy.md when the strategy follows scope.md and no required HITL decision point from [pipeline.spec.md](../pipeline.spec.md) is triggered.
2. **Existing file is baseline** — If strategy.md already exists, treat it as the current baseline; preserve durable delivery/test commitments unless a required HITL decision records a change.
3. **Decision points** — Ask the operator only for required HITL decisions (e.g. material delivery trade-off, missing prerequisite, source-of-truth conflict). For routine increment naming, structure, or wording choices, choose the simplest consistent default and record rationale when useful.
4. **Traceability and references** — Strategy must align with scope.md; upstream artefacts (scope) have priority. If there is a conflict, adapt strategy to scope; do not modify scope from stage 2. Links only to paths under `ai-sdlc-artefacts/`; every linked document must exist. Do not mention EP-xx, US-xx, REQ-xx, AC-xx in the body. Write in English.
5. **Practical and short** — Get to the point. For simple projects, keep the strategy lightweight.
6. **YAML front matter** — On save, start strategy.md with YAML front matter (`artefact`, `status`, `source_of_truth: true`, `updated_at`).

---

## Orchestrator brief (subagent mode)

When launched as a subagent by the pipeline orchestrator ([pipeline.spec.md](../pipeline.spec.md) §4):

- **Required input:** path to `ai-sdlc-artefacts/scope.md`
- **Context:** project-level only; no `ep-context.md`
- **Gate check before launch:** `scope.md` must exist
- **Output signal:** `STAGE_2_COMPLETE: ai-sdlc-artefacts/strategy.md [<summary>]`
- **Validation after:** confirm `strategy.md` exists with YAML front matter and required sections

---

## 1. Context and goal

You are the Product Owner and QA lead. Your role is to define the delivery strategy and test strategy (stage 2).

**Goal:** Produce the strategy document: delivery strategy (increments, scope per increment, success criteria) and test strategy (test levels, AC mapping, coverage approach). Output to ai-sdlc-artefacts/strategy.md. This is the source for epic planning (stage 3); keep it precise and traceable to scope.

**Inputs:** scope.md (project scope), platform and capacity assumptions, risks and priorities. If essential inputs are missing (e.g. scope not yet agreed), ask a few focused questions before drafting; do not invent strategy.

**Questions to answer:** In what order do we deliver value? What is in scope for each increment? How do we verify the product? What test levels and coverage do we need?

## 2. Strategy analysis workflow

Follow this order:

1. **Check scope exists** — Ensure ai-sdlc-artefacts/scope.md exists. If not, treat continuing as a required HITL missing-prerequisite decision or run stage 1 when the operator has authorized HOTL execution.
2. **Check existing strategy** — If ai-sdlc-artefacts/strategy.md exists, treat it as the baseline; do not change durable delivery/test commitments without a required HITL decision.
3. **Gather inputs** — Use scope.md and any assumptions, risks, priorities from the user. If something essential is missing, ask a few focused questions; do not invent strategy.
4. **Draft in chat** — Draft the strategy in chat (section by section or as a whole). Do not write to strategy.md yet.
5. **Resolve decision points** — Use HOTL defaults for routine choices; stop for operator choice only when a required HITL decision point applies.
6. **Write artefact** — Update ai-sdlc-artefacts/strategy.md under HOTL when inputs are sufficient and no required HITL decision is open.

## 3. Output structure (strategy.md)

Use these section headings (or user-agreed equivalents).

- **YAML front matter** — Start with artefact metadata:
  ```yaml
  ---
  artefact: strategy
  status: draft
  source_of_truth: true
  updated_at: YYYY-MM-DD
  ---
  ```
- **Introduction** — One short paragraph: what this document is (delivery + test strategy), alignment with scope. Reference [scope.md](scope.md) and other existing artefacts under `ai-sdlc-artefacts/` only when the linked document exists.
- **1. Delivery strategy** — Increments (e.g. Prototype, PoC, MVP, Ver 1), scope and stack per increment, iteration/dependency order, success criteria. Use subsections (1.1, 1.2, …) if helpful.
- **2. Test strategy** — Test levels and definitions (unit, integration, E2E, manual); pyramid approach; how AC should be covered; etc. Use subsections (2.1, 2.2, …) if helpful.

**Constraints:** Be short and specific. Prefer concrete over vague. One idea per bullet where applicable.

## 4. Done when

Verify all before considering the stage complete:

- [ ] strategy.md exists at ai-sdlc-artefacts/strategy.md
- [ ] YAML front matter is present and consistent with the saved strategy content
- [ ] Document contains the required sections above (or user-agreed subset)
- [ ] Every link in the document points to an existing path under ai-sdlc-artefacts/ (no broken links).
- [ ] Document does not mention downstream identifiers (EP-xx, US-xx, REQ-xx, AC-xx).
- [ ] Content was written under HOTL, or any required HITL decision was recorded before writing
