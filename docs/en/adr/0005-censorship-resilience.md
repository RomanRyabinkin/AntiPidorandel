# ADR-0005: Censorship Resilience

Date: 2025-08-14
Status: Accepted

## Context
Operation is needed under network blocking conditions without mandatory VPN.

## Decision
- Transport ladder: H3 → H2/WS → SSE/long-poll → TURN/WebRTC.
- Multiple domains/ASNs; flocks of mailbox nodes.
- Signed bootstrap manifest from multiple sources.

## Consequences
- Increased resilience; more complex operation (multiple providers).
