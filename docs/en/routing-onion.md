# P2P and Onion Routing

## Idea
A message is encrypted in layers ("onion") for a chain of hops: each hop knows only the **next** one.
Nodes act as relays and/or temporary store-and-forward holders.

## Layers
- Outer layer: address of the next hop (public key), MAC, TTL.
- Inner layers: repeated for each hop.
- The innermost layer is an encrypted envelope for the recipient (AEAD).

## Route Selection
- 3 hops by default (configurable).
- Selected from a pool of nodes: randomly, considering ASN/provider diversity.
- For offline delivery, the last hop is one of the recipient's **flock mailbox nodes**.

```mermaid
flowchart LR
	A[Sender] --> H1[Hop 1]
	H1 --> H2[Hop 2]
	H2 --> H3[Hop 3 / Mailbox]
	H3 --> B[Recipient]
```
