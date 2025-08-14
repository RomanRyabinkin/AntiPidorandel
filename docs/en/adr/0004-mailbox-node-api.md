# ADR-0004: Mailbox Node API

Date: 2025-08-14
Status: Accepted

## Context
A store-and-forward layer without a single point of failure is needed.

## Decision
- HTTP API `/node/put`, `/node/get`, `/node/ack`.
- Signed ACK (Ed25519) for the string `ack:<uuid>` (optional).
- Replication to `NODE_PEERS` (best-effort).

## Consequences
- Delivery works even if the main domain is unavailable.
