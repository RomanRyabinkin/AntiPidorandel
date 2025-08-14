# Censorship resistance


## Layered Approach
1. **Transports**: H3 (QUIC) → H2/WS(443) → SSE/long‑poll → TURN/WebRTC DataChannel.
2. **Infrastructure**: multiple domains/ASN/clouds; some traffic via repeaters.
3. **P2P storage**: mailer node swarms allow offline message delivery without a central domain.
4. **Bootstrap**: signed list of endpoints in several showcases (GitHub Pages, mirrors).
5. **Availability telemetry** (no identifiers): prompts clients for the order of transports.

## What to consider
- TLS fingerprint: as "browser-like" as possible. Headers - no exotics.
- Padding and jitter: fixed envelope size, random delivery delays.
- It’s better not to count on “domain fronting”, since it is closed almost everywhere.
- UDP can be blocked - a reliable path via TCP/H2/443 is required.