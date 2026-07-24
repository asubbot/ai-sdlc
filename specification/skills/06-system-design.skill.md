---
name: system-design.skill
description: Produce epic system design from ep-requirements and ep-acceptance-criteria; output ep-system-design.md. Use when defining or refining system design (stage 6), e.g. "system design for this epic", "architecture", "components and interfaces".
---

# Stage 6: System design

**Pipeline:** [pipeline.spec.md](../pipeline.spec.md)  
**Output:** ai-sdlc-artefacts/epics/<epic-id>/ep-system-design.md; update ep-context.md

---

## Orchestrator brief (subagent mode)

When launched as a subagent by the pipeline orchestrator ([pipeline.spec.md](../pipeline.spec.md) §4):

- **Required input:** epic ID, paths to `ep-requirements.md`, `ep-acceptance-criteria.md`; optional: latest `ep-system-design-review.md` (when iterating per §2.1)
- **Context:** `ep-context.md` (read first for orientation)
- **Gate check before launch:** `ep-requirements.md` and `ep-acceptance-criteria.md` must exist
- **Output signal:** `STAGE_6_COMPLETE: ai-sdlc-artefacts/epics/<epic-id>/ep-system-design.md [<key design decisions>] [context_delta: <Design Decisions + Interfaces + Open Questions + Links summary>]`
- **Validation after:** `./tools/validate/validate structure EP-XXX` (artefact structure, from repo root)
- **ep-context.md:** Subagent does **not** write `ep-context.md`; include compact context updates in `context_delta` for the orchestrator to apply

---

## Core Principles

Follow these principles for all system design work:

1. **HOTL artefact writes** — In pipeline execution, create or overwrite ep-system-design.md when requirements and acceptance criteria exist, traceability is maintained, and no required HITL decision point from [pipeline.spec.md](../pipeline.spec.md) is triggered.
2. **Existing file is baseline** — If ep-system-design.md already exists for the epic, treat it as the current baseline; preserve durable architecture/contract decisions unless a required HITL decision records a change.
3. **Decision points** — Ask the operator only for required HITL decisions (e.g. material architecture/security/reliability trade-off, missing prerequisite, source-of-truth conflict). For routine structure depth, module grouping, or wording choices, choose the simplest consistent default and record rationale when useful.
4. **References** — Relative links only to paths under `ai-sdlc-artefacts/`; every linked document must exist. External **`https://`** links to upstream pattern/vendor documentation are allowed. Keep traceability to ep-requirements. Write in English.
5. **Full REQ traceability** — Every requirement (REQ-EE.NNN format) from ep-requirements.md must be mentioned at least once in ep-system-design.md; the document must not omit any REQ.
6. **Practical and short** — Get to the point. Be practical above all. Be short and specific.
7. **Legacy** — Do not modify content under `legacy` folders; use as reference only (e.g. technical discovery or research).
8. **Token-optimized context** — Read current ep-context.md first when present, but open full requirements and acceptance criteria for traceability. On save, add YAML front matter to ep-system-design.md. **Orchestrated subagent mode:** report Design Decisions, Interfaces / Contracts, Open Questions, and Links deltas in `context_delta`; orchestrator updates ep-context.md. **Solo/HOTL without orchestrator:** refresh those sections in ep-context.md.

---

## 1. Context and goal

You are the Tech Lead for this epic. Your role is to produce the epic system design document (stage 6).

**Goal:** Produce ep-system-design.md: components, interfaces, data models, and key design decisions. This output is the input for system design review (stage 7), then implementation planning (stage 8); keep it traceable to ep-requirements and, where useful, to ep-acceptance-criteria. Optionally include technical discovery (options, comparison, recommendation, risks) and references to research under `legacy/` (read-only).

**Design–review iteration ([pipeline.spec.md](../pipeline.spec.md) §2.1):** When **stage 7** reports **Blocker**, **Major**, **Medium**, or **Minor** findings in the latest `## Review iteration N` section of `ep-system-design-review.md` (or in chat before save), run **stage 6** again: revise `ep-system-design.md` to resolve those findings, then the orchestrator schedules another **stage 7** (delegated, fresh session). Repeat until **zero** Blocker/Major/Medium/Minor or until the **operator decides** after **five** stage 7 iterations.

**Inputs:** ep-context.md when present and current, ep-requirements.md and ep-acceptance-criteria.md (ai-sdlc-artefacts/epics/<epic-id>/). If either source artefact is missing, treat continuing as a required HITL missing-prerequisite decision or run the missing stage when the operator has authorized HOTL execution. When addressing review feedback, read the latest **Current Gate Summary** and latest relevant **`ep-system-design-review.md`** iteration section (and prior sections if needed for context). Platform constraints and research or technical discovery (e.g. under epic `legacy/`) may be used as reference.

**Questions to answer:** How is the system structured? What are the main components and interfaces? What are the key design decisions and risks?

---

## 2. System design workflow

Follow this order:

1. **Check inputs** — Ensure ep-requirements.md and ep-acceptance-criteria.md exist for the epic. If not, treat continuing as a required HITL missing-prerequisite decision or run the missing stage when the operator has authorized HOTL execution. If ep-context.md exists and is current, read it first for orientation; if stale or missing, fall back to source artefacts.
2. **Check existing ep-system-design** — If ep-system-design.md exists for the epic, treat it as the baseline; propose changes as edits.
3. **Consult architecture patterns catalog** — Identify architecturally significant decisions (module boundaries, sync vs async, resilience, consistency, auth boundaries). For each, open `reference/architecture-patterns/index.md` (at the ai-sdlc checkout root), pick 1–3 relevant cards, and fetch their `sources` for trade-offs. Record in the **Design decisions** section of ep-system-design.md (see § Output structure; not only in ep-context.md): **chosen / rejected / why**, the marker `architecture-pattern: <pattern-id>` (pattern id as plain text, no links into `ai-sdlc/`), and upstream https links; respect each card's `kiss_default` and `when_not`. When no card applies to a decision, record `architecture-pattern: n/a — <one-line reason>`. Skip this step entirely for trivial wording/layout edits with no architecturally significant decisions. If the catalog directory is absent from the checkout, record `architecture-pattern: n/a — catalog unavailable in checkout` and continue.
4. **Draft** — Draft the design (section by section or by block). In HOTL mode, proceed to write when the draft is internally consistent and no required HITL decision is open.
5. **Resolve decision points** — Use HOTL defaults for routine choices; stop for operator choice only when a required HITL decision point applies.
6. **Write artefact** — Create or update ai-sdlc-artefacts/epics/<epic-id>/ep-system-design.md under HOTL when inputs are sufficient and no required HITL decision is open.
7. **Render C4 C2 PNG** — After saving `diagrams/c4-container.puml`, follow [ep-plantuml-export.skill.md](ep-plantuml-export.skill.md): run `plantuml -tpng diagrams/c4-container.puml` from the epic directory, verify `diagrams/c4-container.png` exists, and embed the centered `<img>` in ep-system-design.md before the Source line. If C3 was added, render `c4-component-go.puml` the same way.
8. **Refresh epic context** — **Orchestrated subagent mode:** report Design Decisions, Interfaces / Contracts, Open Questions, and Links deltas in the output signal (`context_delta`); do not edit ep-context.md. **Solo/HOTL without orchestrator:** update ep-context.md concisely; do not copy full design sections.
9. **Legacy** — Do not modify content under `legacy` folders; use as reference only.

---

## 3. Output structure (ep-system-design.md)

Use these section headings (or user-agreed equivalents).

Start the document with YAML front matter:

```yaml
---
artefact: ep-system-design
epic_id: EP-XXX
status: draft
source_of_truth: true
updated_at: YYYY-MM-DD
---
```

- **Overview** — Brief summary of the system and design scope, with traceability to key requirements. One short paragraph or bullet list.
- **Architecture** — High-level structure and boundaries. **Must include C4 Level 2 (Containers):** C4-PlantUML diagram: source in `diagrams/c4-container.puml`, PNG in `diagrams/c4-container.png`. In ep-system-design: centered image, then "Source:" with link to .puml and regeneration command (`plantuml -tpng diagrams/c4-container.puml` from epic directory). System context (C1) is in ep-requirements; C2 is drawn here. Optional subsection **Module boundaries**: layers/modules, dependency rules, and wiring responsibilities; optional diagram and verification step (script or checklist).
- **Components and interfaces** — Table or list: Component | Responsibility | Key interface/contract. Each component should trace to requirements where applicable.
- **Data models** — Main entities, schemas, and state transitions relevant to the design. Reference upstream artefacts under ai-sdlc-artefacts only.
- **Design decisions** — Required when the design contains architecturally significant decisions (workflow step 3); omit only when there are none. One entry per decision: **chosen / rejected / why**, the marker `architecture-pattern: <pattern-id>` or `architecture-pattern: n/a — <one-line reason>`, and upstream https links. These records live **in this section of ep-system-design.md** (stage 7 reviews this document, not ep-context.md); ep-context.md carries only a compact summary. Example entry:

```markdown
- DD-1: Deliver scheduler notifications via an outbox table.
  architecture-pattern: transactional-outbox
  Chosen: atomic write of task + event in one transaction.
  Rejected: direct send inside the transaction (dual-write, loss on crash).
  Upstream: https://microservices.io/patterns/data/transactional-outbox.html
- DD-2: Package layout for the notify module.
  architecture-pattern: n/a — local code grouping without cross-component trade-offs
```

- **Error handling** — How validation and runtime failure modes are handled; trace to requirements where applicable.
- **Testing strategy** (optional) — Short note on unit, integration, e2e, and deploy verification; can reference strategy.md under ai-sdlc-artefacts.
- **Optional:** Technical discovery, options comparison, recommendation, risks; references to research (e.g. legacy/technical-discovery) only for reading; no edits to legacy.

**Diagrams:** Use the epic's `diagrams/` folder (same as in ep-requirements); store `c4-container.puml` there. **Automatically render** PNG via [ep-plantuml-export.skill.md](ep-plantuml-export.skill.md) (`diagrams/c4-container.png`); do not finish stage 6 with `.puml` only.

**TOC:** Include a first-level table of contents (links to section anchors).

**Traceability:** In the body, link to ep-requirements.md sections (e.g. [REQ-01.001](ep-requirements.md#interface-and-deployment)). **Every REQ from ep-requirements.md must be referenced at least once** in the document (full requirement traceability). Requirement IDs use format REQ-EE.NNN (epic number EE, requirement NNN). Every relative linked path must be under `ai-sdlc-artefacts/` and exist; external `https://` links to upstream pattern/vendor documentation are allowed.

**Quality:** Be specific and testable where possible; avoid vague wording. Keep the document maintainable for system design review (stage 7) and implementation planning (stage 8).

---

## 4. Done when

Verify all before considering the stage complete:

- [ ] ep-system-design.md exists at ai-sdlc-artefacts/epics/<epic-id>/ep-system-design.md
- [ ] YAML front matter is present and consistent with the saved system design content
- [ ] ep-context.md was refreshed with compact Design Decisions, Interfaces / Contracts, Open Questions, and Links, **or** (orchestrated subagent mode) `context_delta` was reported for the orchestrator to apply
- [ ] Document contains **Overview** (system summary and traceability to REQ), **Architecture** including **C4 C2** (source in `diagrams/c4-container.puml`, PNG in `diagrams/c4-container.png` embedded centered; Source line with regeneration command), **Components and interfaces**, and **Data models** (if applicable); **Error handling** and **Testing strategy** recommended where relevant
- [ ] `diagrams/c4-container.png` exists and was rendered per [ep-plantuml-export.skill.md](ep-plantuml-export.skill.md) (not `.puml`-only)
- [ ] Every relative link in the document points to an existing path under `ai-sdlc-artefacts/` (no broken links); external links, where present, use `https://` and point to upstream documentation
- [ ] Every architecturally significant decision has an entry in the **Design decisions** section of ep-system-design.md with an **`architecture-pattern:`** field: either `<pattern-id>` (plus chosen / rejected / why and upstream https link) or `n/a — <one-line reason>` (catalog: `reference/architecture-patterns/` at the ai-sdlc checkout root)
- [ ] Traceability to ep-requirements is maintained: **every REQ from ep-requirements.md is referenced at least once** in the document
- [ ] Content was written under HOTL, or any required HITL decision was recorded before writing
