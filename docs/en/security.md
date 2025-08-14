# Security

## Possible Threats
- Passive network observer (DPI, ISP).
- Partially malicious relays.
- Loss of mailbox node.
- Replays/spam/queue overflow.

## Measures
- E2EE: servers see only the padded envelope (AEAD).
- Sealed sender (planned): hide the sender from relays.
- Signed ACK (Ed25519): only the inbox owner can delete.
- Idempotency by `message_id` (UUIDv7/ULID).
- PoW/quotas/ban by key on `PUT`.
- Logs without payload/PII; encrypted DB backups.
