# Midterm Slides (5–7)

## Slide 1: Project and SRE angle

- Clothes Store: real stack (Go, MongoDB, templates + static assets).
- Midterm goals: containerization, SLOs, monitoring, alerting.

## Slide 2: Docker architecture

- Diagram: Nginx (frontend) → Go (backend) → MongoDB.
- Observability: Prometheus, Grafana, Node Exporter.
- Single stack: `docker compose up`.

## Slide 3: Data and metrics

- Where: Prometheus middleware → `/metrics`; scrape in `prometheus.yml`.
- How to populate charts: HTTP traffic, sign-ups, orders in MongoDB.

## Slide 4: SLI / SLO / Error Budget

- SLIs: availability (non-5xx share), p95 latency.
- SLOs: 99.5% availability; p95 < 1s within the agreed window budget.
- Error budget numbers from the report.

## Slide 5: Grafana

- Golden Signals + SLO panels (dashboard `midterm-overview`).
- Datasource: Prometheus (provisioning).

## Slide 6: Alerts and incident demo

- Warning + critical rules in `alert_rules.yml`.
- Demo: `docker compose stop backend` → FIRING.

## Slide 7 (optional): Swarm / Next steps

- Swarm is optional; bonus: `docker swarm init` + `docker stack deploy`.
- Future work: separate SPA, managed Mongo, secrets management.
