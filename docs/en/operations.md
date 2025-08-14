# Operations

- `/healthz` — alive/ready (200 ok).
- Логи — stdout/stderr (Docker best practice).
- Image - distroless:nonroot.
- Monitoring (plan): Prometheus metrics.

## Deploying a "flock" of repeaters
- Set up 3-5 mail nodes in different clouds/ASN.
- Enable `NODE_ENABLED=true` and configure `NODE_PEERS` for replication.
- Set up backup domains/certificates.

## Backups
- Scheduled DB dump; encrypted storage.