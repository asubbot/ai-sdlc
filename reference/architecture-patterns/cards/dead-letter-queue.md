---
# OKF v0.1
type: Pattern
title: Dead Letter Queue
description: Route messages that repeatedly fail processing to a side channel
  for inspection, replay, or discard instead of blocking the main queue.
timestamp: 2026-07-27
tags: [messaging, reliability, operations]

# ai-sdlc extensions (advisory selection policy)
sources:
  - url: https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dead-letter-queues
    note: Primary upstream (Azure Service Bus dead-letter queues)
  - url: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html
    note: Amazon SQS dead-letter queues
forces: [poison-message isolation, operability of failed work, retention and PII risk]
when: async consumers can fail on bad payloads or permanent errors and must not
  stall the healthy queue forever
when_not: there is no queue/async path, or failures are always transient and
  already bounded by retry without needing a side channel
kiss_default: bounded retries + surface the error; add a DLQ when poison messages
  are observed or when operators need a replay path
quality: [reliability, operability]
related: [retry-and-timeouts, idempotency, publisher-subscriber, sync-vs-async]
---

# Dead Letter Queue

**Problem.** A poison message that always fails can block a consumer, burn
retries, and stop healthy work behind it.

**Options.**
- Drop/fail after N retries with alerting — simplest for low volume.
- Dead-letter queue: move exhausted messages aside with reason metadata.
- Quarantine + automated classification — for high-volume platforms.

**Failure modes.** DLQ without owners or replay runbook (silent pile-up);
replaying without fixing the root cause; retaining sensitive payloads too long.
