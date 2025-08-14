# Spirit — Overview

**Spirit** — is a privacy-first messenger that uses client-to-client store-and-forward on top of a relay network. Messages are end-to-end encrypted and routed through multiple hops using onion-style layered encryption. The server can act as a blind relay and/or a mailbox node.

- 🧭 **Goals**: P2P/store-and-forward, minimal metadata, resilience without a mandatory VPN.
- 🔐 **Principles**: E2EE, sealed sender (planned), fixed-size padded envelopes, replay protection.
- 🧱 **Censorship resistance**: multiple transports and domains, a fallback ladder, mailbox swarms, and multi-source bootstrap.

> **Rename**: the project was previously called AntiPidorandel (working title). The current name is**Spirit**.