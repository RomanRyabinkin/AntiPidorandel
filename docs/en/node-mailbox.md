# Mailbox Node

The Mailbox node stores encrypted incoming envelopes using the recipient's key and delivers them upon request.

## Purpose of a swarm
- `swarm = k` nodes based on consistent hash `hash(PubKey)`.
- The sender places an envelope on all `k` nodes (or `m` of `n` in erasure coding).
- The receiver periodically `GET`s to its relay flock and removes envelopes via `ACK` (signed).

## HTTP API

### `POST /node/put`
Puts down the envelope.
```json
{
  "to": "bob",
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "header_b64": "...",
  "nonce_b64": "...",
  "cipher_b64": "...",
  "ttl_seconds": 1209600
}
```

### `GET /node/get?to=<user>`
Returns a list of envelopes (up to `PENDING_BATCH`).

### `POST /node/ack`
Deletes by `message_id`. Recommended to **sign**:
```json
{
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "pubkey_b64": "<ed25519-pub>",
  "sig_b64": "<ed25519-sign('ack:'+message_id)>"
}
```

## Replication and neighbor nodes
- After `put`, the node can **replicate** to `NODE_PEERS` (best‑effort).
- ACK removal occurs on all replicas.