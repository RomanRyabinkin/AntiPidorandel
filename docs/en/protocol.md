# Protocol

## WebSocket frames (MVP)
```json
// client -> server
{ "type":"send", "to":"<user>", "message_id":"<uuid>", "header_b64":"...", "nonce_b64":"...", "cipher_b64":"..." }
{ "type":"ack",  "message_id":"<uuid>", "ack_pubkey_b64":"...", "ack_sig_b64":"..." }
{ "type":"ping" }

// server -> client
{ "type":"deliver", "to":"<user>", "message_id":"<uuid>", "header_b64":"...", "nonce_b64":"...", "cipher_b64":"..." }
{ "type":"pong" }
```
> `ack_pubkey_b64/ack_sig_b64` are optional fields. If `REQUIRE_ACK_SIGNATURE=true` is enabled, the server checks `Ed25519.Verify(pub, "ack:"+uuid, sig)`.

## HTTP Mailbox
Endpoints: `/node/put`, `/node/get`, `/node/ack`. See the Mailer Nodes section for details.