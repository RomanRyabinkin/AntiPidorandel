# Discovery and Bootstrap

## Bootstrap Manifesto
- A signed file with a list of domains, IPs, ports and transport priorities.
- Splits into several "showcases" (GitHub Pages, mirrors, main site).
- The client caches and periodically updates; verifies the signature with the project key.

## DHT (later)
- Kademlia: `hash(PubKey) -> swarm узлов`.
- Record validation: signatures + rate‑limit.

## Happy‑eyeballs for transports
- The client tries H3→H2/WS→SSE/long‑poll→TURN/WebRTC in parallel and chooses the first successful path.