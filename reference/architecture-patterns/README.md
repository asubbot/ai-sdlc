# Architecture patterns reference

Advisory catalog of thin architecture-pattern cards for **stage 6 (system design)**.
It is **not** part of the normative pipeline: artefact paths, gates, and validate
rules do not depend on this directory.

## Purpose

- Give the stage 6 agent a fast, curated shortlist of patterns to consult when a
  design contains **architecturally significant decisions** (ASD): module
  boundaries, sync vs async, resilience, consistency, auth boundaries.
- Keep cards **thin**: forces, when / when-not, KISS default, and links. The
  **source of truth for the pattern body is the upstream URL** in `sources`.

## Anti-goals

- No copies of upstream articles; card bodies stay ≤ ~30 lines.
- No product-specific knowledge base (that lives outside ai-sdlc).
- No new validate rules or severity tiers; enforcement happens through the
  stage 6 Done-when and the stage 7 review checklist.

## How agents use this catalog

1. During stage 6, identify architecturally significant decisions.
2. Open [index.md](index.md), pick 1–3 relevant cards, read `when_not` and
   `kiss_default` first, fetch `sources` for trade-offs.
3. Record in Design Decisions: **chosen / rejected / why**, the marker
   `architecture-pattern: <pattern-id>`, and upstream https links.
4. When no card applies to an ASD, record
   `architecture-pattern: n/a — <one-line reason>`.
5. If this directory is absent from the checkout (older pin), record
   `architecture-pattern: n/a — catalog unavailable in checkout` and continue.

## Citation rules (in SDLC artefacts)

- Reference cards by **pattern id as plain text** (the filename stem, e.g.
  `transactional-outbox`) plus the **upstream https URL**.
- Do **not** add markdown links from artefacts into `ai-sdlc/` paths; relative
  artefact links must stay under `ai-sdlc-artefacts/`. External `https://`
  links to upstream pattern documentation are allowed.

## Card format

See [SCHEMA.md](SCHEMA.md). Front matter is aligned with
[Open Knowledge Format (OKF) v0.1](https://github.com/GoogleCloudPlatform/knowledge-catalog/blob/main/okf/SPEC.md):
required `type: Pattern`, `title`, `description`, `timestamp`, `tags`, plus
ai-sdlc extension keys. Card identity is the filename.
