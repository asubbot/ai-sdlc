# Architecture patterns index

MOC for stage 6 consultation. Pick 1–3 cards per architecturally significant
decision; read `when_not` and `kiss_default` before adopting a pattern.

| Pattern id | Core force (one line) | Card |
|------------|----------------------|------|
| `module-boundaries` | Keep modules cohesive and dependencies one-way as the codebase grows | [cards/module-boundaries.md](cards/module-boundaries.md) |
| `sync-vs-async` | Decouple callers from slow or unreliable work without losing responses | [cards/sync-vs-async.md](cards/sync-vs-async.md) |
| `retry-and-timeouts` | Survive transient faults without amplifying load or hanging forever | [cards/retry-and-timeouts.md](cards/retry-and-timeouts.md) |
| `circuit-breaker` | Stop hammering a dependency that is already failing | [cards/circuit-breaker.md](cards/circuit-breaker.md) |
| `bulkhead` | Isolate capacity so one dependency cannot exhaust the whole process | [cards/bulkhead.md](cards/bulkhead.md) |
| `rate-limiting` | Cap burst traffic to protect capacity, cost, and partner quotas | [cards/rate-limiting.md](cards/rate-limiting.md) |
| `idempotency` | Make retries and duplicate deliveries safe | [cards/idempotency.md](cards/idempotency.md) |
| `transactional-outbox` | Deliver an event if and only if the local transaction committed | [cards/transactional-outbox.md](cards/transactional-outbox.md) |
| `dead-letter-queue` | Park poison messages so they do not stall the healthy queue | [cards/dead-letter-queue.md](cards/dead-letter-queue.md) |
| `publisher-subscriber` | Fan out events to independent consumers via a broker topic | [cards/publisher-subscriber.md](cards/publisher-subscriber.md) |
| `saga-or-compensating` | Coordinate multi-step work across commit boundaries with compensations | [cards/saga-or-compensating.md](cards/saga-or-compensating.md) |
| `caching` | Cut latency/load on hot reads while controlling staleness | [cards/caching.md](cards/caching.md) |
| `authn-boundary` | Concentrate authentication/authorization at an explicit trust boundary | [cards/authn-boundary.md](cards/authn-boundary.md) |
| `strangler-fig` | Replace legacy incrementally by routing slices to the new system | [cards/strangler-fig.md](cards/strangler-fig.md) |
| `health-liveness-readiness` | Give orchestrators a clear alive vs ready-to-serve contract | [cards/health-liveness-readiness.md](cards/health-liveness-readiness.md) |
