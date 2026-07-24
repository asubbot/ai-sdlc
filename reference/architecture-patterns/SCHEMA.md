# Card schema

Normative for this package only (not for pipeline.spec). Every file under
`cards/` MUST follow this schema.

## Identity

The card identity is the **filename**: `cards/<pattern-id>.md`. There is no
`id` key in front matter. SDLC artefacts reference cards via
`architecture-pattern: <pattern-id>` (filename stem, plain text).

## Front matter

Aligned with [OKF v0.1](https://github.com/GoogleCloudPlatform/knowledge-catalog/blob/main/okf/SPEC.md):
markdown file with parseable YAML front matter; unknown keys are preserved by
consumers.

### OKF fields (required)

| Key | Value |
|-----|-------|
| `type` | Always `Pattern` |
| `title` | Human-readable pattern name |
| `description` | 1–2 sentence summary of the pattern |
| `timestamp` | Date the card was last edited (`YYYY-MM-DD`) |
| `tags` | Short lowercase topic tags |

### ai-sdlc extension fields (required unless noted)

| Key | Value |
|-----|-------|
| `sources` | List of `{url, note}`; first entry is the primary upstream (SOT for the pattern body). Prefer official vendor docs |
| `forces` | Short list of forces/constraints the pattern balances |
| `when` | One sentence: the situation where the pattern applies |
| `when_not` | One sentence: where the pattern is overkill or harmful |
| `kiss_default` | The simplest alternative to try first, and when to escalate to the pattern |
| `quality` | Short quality-attribute tags (e.g. `reliability`, `consistency`) |
| `related` | Optional list of related pattern ids (filename stems) |

## Body

English, ≤ ~30 lines, three sections:

1. **Problem.** Sketch of the problem, 1–3 lines.
2. **Options.** Bullet list of realistic alternatives including the pattern.
3. **Failure modes.** What goes wrong when the pattern is applied.

Do not paste upstream article content; the body is an index entry, not a copy.
