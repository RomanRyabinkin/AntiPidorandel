# ADR-0003: WebSocket Frame Contract

Date: 2025-08-14
Status: Accepted

## Context
A simple and compatible delivery and acknowledgment contract is needed.

## Decision
- Types: `send`, `deliver`, `ack`, `ping`, `pong`.
- Binary fields — base64.
- Idempotency by `message_id`.
- Optionally — `ack_pubkey_b64` and `ack_sig_b64` (Ed25519).

## Consequences
- Minimal client changes when adding ACK signature support.
