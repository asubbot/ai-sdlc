---
name: scope-analysis.skill
description: Capture project scope from user request or conversation; produce scope.md as single source for strategy and epic planning. Use when the user wants to define or refine project scope (e.g. "define scope", "what are we building?", "capture the scope", "scope this feature").
---

# Stage 1: Scope analysis

**Pipeline:** [pipeline.spec.md](../pipeline.spec.md)
**Output:** ai-sdlc-artefacts/scope.md

## Core Principles

Follow these principles for all scope analysis work:

1. **HOTL artefact writes** — In pipeline execution, create or overwrite scope.md when the scope draft is internally consistent and no required HITL decision point from [pipeline.spec.md](../pipeline.spec.md) is triggered.
2. **Existing file is baseline** — If scope.md already exists, treat it as the current baseline; preserve durable scope unless a required HITL decision records a scope change.
3. **Decision points** — Ask the operator only for required HITL decisions (e.g. missing essential scope input, material scope trade-off, source-of-truth conflict). For routine structure or wording choices, choose the simplest consistent default and record rationale when useful.
4. **References** — Links only to paths under `ai-sdlc-artefacts/`; every linked document must exist.
5. **YAML front matter** — On save, start scope.md with YAML front matter (`artefact`, `status`, `source_of_truth: true`, `updated_at`).

---

## Orchestrator brief (subagent mode)

When launched as a subagent by the pipeline orchestrator ([pipeline.spec.md](../pipeline.spec.md) §4):

- **Required input:** user request / chat context; optional existing `ai-sdlc-artefacts/scope.md`
- **Context:** project-level only; no `ep-context.md`
- **Gate check before launch:** none (first pipeline stage)
- **Output signal:** `STAGE_1_COMPLETE: ai-sdlc-artefacts/scope.md [<summary>]`
- **Validation after:** confirm `scope.md` exists with YAML front matter and required sections (no epic `validate structure` for project-level artefacts)

---

## 1. Context and goal

You are an expert requirements analyst. Your role is to capture project scope from the user's request or conversation.

**Goal:** Produce a project scope document that answers: What are we building? What terms do we use? What is in scope, out of scope, or deferred? The output is the source for strategy (stage 2) and epic planning (stage 3); keep it precise and traceable.

**Inputs:** Chat request, stakeholder vision, problem statement, success criteria, constraints (platform, audience), and any references the user provides. If essential inputs are missing, ask a few focused questions before drafting; do not invent scope.

## 2. Scope analysis workflow

Follow this order:

1. **Check existing scope** — If ai-sdlc-artefacts/scope.md exists, treat it as the baseline; do not change durable scope without a required HITL decision.
2. **Gather inputs** — Use chat and any references. If something essential is missing, ask a few focused questions; do not invent scope.
3. **Draft in chat** — Draft the scope in chat (section by section or as a whole). Do not write to scope.md yet.
4. **Resolve decision points** — Use HOTL defaults for routine choices; stop for operator choice only when a required HITL decision point applies.
5. **Write artefact** — Update ai-sdlc-artefacts/scope.md under HOTL when inputs are sufficient and no required HITL decision is open.

## 3. Output structure (scope.md)

Use these section headings (or user-agreed equivalents).

- **YAML front matter** — Start with artefact metadata:
  ```yaml
  ---
  artefact: scope
  status: draft
  source_of_truth: true
  updated_at: YYYY-MM-DD
  ---
  ```
- **Introduction** — Short summary of the project or feature and what this document covers (2–4 sentences).
- **Glossary** — Key system names and technical terms the team will use. One row per term: term and short definition. Only terms that affect scope or downstream stages. Example: "Personal Assistant (PA)" — "Agent-driven app that executes user requests in the repo and environment."
- **In scope** — What is included: capabilities, features, or themes. Use concrete, testable phrasing (e.g. "Telegram bot for conversation" not "we will have a bot"). Bullet list.
- **Out of scope / deferred** — What is explicitly excluded or postponed, with brief reason if helpful. Bullet list.

**Constraints:** Be short and specific. Prefer concrete over vague. One idea per bullet.

## 4. Done when

Verify all before considering the stage complete:

- [ ] scope.md exists at ai-sdlc-artefacts/scope.md
- [ ] YAML front matter is present and consistent with the saved scope content
- [ ] Document contains the required sections above (or user-agreed subset)
- [ ] Every link in the document points to an existing path under `ai-sdlc-artefacts/` (no broken links).
- [ ] Document does not mention downstream identifiers (EP-xx, US-xx, AC-xx).
- [ ] Content was written under HOTL, or any required HITL decision was recorded before writing
