# Midterm Project Report: Production-Ready Observability

## 1) Custom Architecture and Containerization

Project: **Clothes Store** (real application code, not a Hello World).

- **Frontend** (`frontend`): Nginx as the user entrypoint, reverse proxy to the backend.
- **Backend** (`backend`): Go + Gin, business logic, HTML templates, `/metrics` for Prometheus.
- **Database** (`mongo`): MongoDB.
- **Observability**: Prometheus, Grafana (provisioning + dashboard JSON), Node Exporter.

### Artifacts

- `backend/Dockerfile`, `frontend/Dockerfile`, `frontend/nginx.conf`
- `docker-compose.yml`
- `monitoring/prometheus.yml`, `monitoring/alert_rules.yml`
- `monitoring/grafana/provisioning/`, `monitoring/grafana/dashboards/midterm-overview.json`

### Run

```bash
docker compose up --build -d
```

### URLs (current ports)

- App: http://localhost
- Backend from host: http://localhost:8081 (maps `8081:8080`)
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (`admin` / `admin`)
- Node Exporter: http://localhost:9100/metrics

### Where chart data comes from

1. The app exposes metrics on `/metrics` (`backend/internal/middleware/prometheus.go`, wired in `main.go`).
2. Prometheus **scrapes** `backend:8080` and `node-exporter:9100` (`monitoring/prometheus.yml`).
3. Grafana queries Prometheus at `http://prometheus:9090` and renders panels from `midterm-overview.json`.

### How to add data (non-zero charts)

- Browse the site or call APIs (`curl http://localhost:8081/ping`, etc.).
- Register users, create products, and place orders via UI or REST — data persists in MongoDB.

---

## 2) Reliability Targets (SLI, SLO, Error Budget)

### SLI #1: Availability (successful HTTP request ratio)

\[
\text{Availability} = 1 - \frac{\text{5xx}}{\text{all requests}}
\]

PromQL (5m window):

```promql
1 - (
  sum(rate(clothes_store_http_requests_total{status=~"5.."}[5m]))
  /
  clamp_min(sum(rate(clothes_store_http_requests_total[5m])), 0.001)
)
```

**SLO:** 99.5% monthly availability.

**Error budget (30 days):**  
43,200 minutes × 0.5% = **216 minutes** of allowed “bad” time (per your SLO definition).

### SLI #2: Latency (p95 request duration)

PromQL:

```promql
histogram_quantile(
  0.95,
  sum(rate(clothes_store_http_request_duration_seconds_bucket[5m])) by (le)
)
```

**SLO:** p95 < 1.0s for 99% of 5-minute windows per month (as in the assignment).

**Error budget:** 8,640 windows × 1% ≈ **86 windows** violating the target.

---

## 3) Grafana Dashboard (Golden Signals + SLO)

Dashboard **Midterm - Clothes Store Overview** (`monitoring/grafana/dashboards/midterm-overview.json`):

| Panel | Signal | Source |
|-------|--------|--------|
| Traffic | Golden | `rate(clothes_store_http_requests_total[1m])` |
| Errors (5xx share) | Golden | share of 5xx over all requests |
| Latency p95 | Golden + SLO budget | histogram + 0.5s / 1s thresholds |
| In-flight | Golden (saturation) | `clothes_store_http_requests_in_flight` |
| Availability ≥ 99.5% | SLO compliance | availability formula + thresholds |
| Host CPU % | Saturation (node) | `node_cpu_seconds_total` |

---

## 4) Alerting (Prometheus)

File: `monitoring/alert_rules.yml`

1. **ClothesStoreServiceDown** (critical) — scrape unavailable for 1m.
2. **HighErrorRate** (warning) — 5xx share > 5% for 2m.
3. **HighLatencyP95** (critical) — p95 > 1s for 5m.

### Manual trigger (demo FIRING)

```bash
docker compose stop backend
```

Wait ~1 minute → Prometheus → Alerts.

---

## 5) PDF Evidence Checklist

- [ ] `docker compose ps` — all services up.
- [ ] Screenshot: home page at `http://localhost`.
- [ ] Prometheus `/targets` — UP.
- [ ] Grafana — dashboard with Golden Signals and SLO.
- [ ] One alert in **FIRING** state.
- [ ] SLI/SLO/error-budget text from section 2.
