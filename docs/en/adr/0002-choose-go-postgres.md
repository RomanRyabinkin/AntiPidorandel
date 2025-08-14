# ADR-0002: Choosing Go + PostgreSQL

Date: 2025-08-14
Status: Accepted

## Context
A fast and easy-to-deploy stack with minimal dependencies is needed.

## Decision
- Backend: Go 1.23+ (`gorilla/websocket`, `pgx`).
- Storage: PostgreSQL 16+ (indexes, reliability, simple TTL policies).

## Alternatives
- Rust/Node — postponed until post-MVP for time-to-market reasons.

## Consequences
- High performance and ease of operation.
