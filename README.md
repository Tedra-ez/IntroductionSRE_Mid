# Clothes Store

An e-commerce web application built with Go (Gin) and MongoDB: storefront, cart, checkout, admin tools, and analytics. The UI uses HTML templates and static assets. For the midterm, an observability stack (Prometheus, Grafana, Node Exporter) is included and runs via Docker Compose.

## Team

- Yskak Zhanibek
- Nauanov Alikhan
- Zhumagali Beibarys

## Features

### Customer Experience

- **Shop**: filters (category, gender, color, size) and sorting.
- **Product**: product page, sizes, stock.
- **Cart & Wishlist**: client-side persistence.
- **Checkout**: order placement flow.
- **Accounts**: sign up, sign in, order history.

### Admin Dashboard

- **Analytics**: KPIs, revenue, top products.
- **Products**: CRUD and image uploads.
- **Orders**: order status updates.
- **Users**: registered user list.

## Tech Stack

- **Backend**: Go (Gin)
- **Database**: MongoDB
- **Auth**: JWT + `auth_token` cookie
- **Frontend**: HTML templates, CSS, vanilla JavaScript
- **Observability**: Prometheus, Grafana, Node Exporter (midterm)

## Project Structure

```text
├── backend/
│   ├── cmd/server/          # Application entrypoint
│   ├── internal/            # API, handlers, services, middleware (incl. Prometheus)
│   └── Dockerfile
├── frontend/
│   ├── static/              # CSS, JS, assets (uploads: static/assets/products/)
│   ├── templates/           # HTML (Gin templates)
│   ├── nginx.conf           # Reverse proxy to backend (in Docker)
│   └── Dockerfile
├── monitoring/
│   ├── prometheus.yml       # Scrape backend + node-exporter
│   ├── alert_rules.yml      # Alerting rules
│   └── grafana/
│       ├── provisioning/    # Datasource + dashboard provisioning
│       └── dashboards/      # Dashboard JSON (Midterm overview)
├── docker-compose.yml
├── MIDTERM_REPORT.md        # Report draft (SLI/SLO, screenshot checklist)
├── PRESENTATION.md          # Slide-by-slide defense deck (copy to Slides)
└── DEFENSE_SIMPLE_RU.md     # Simple Russian “explain like I’m 5” for oral defense
```

## Where metrics and charts come from

| What | Source |
|------|--------|
| **Application HTTP metrics** | Middleware `backend/internal/middleware/prometheus.go`: `clothes_store_http_requests_total`, `clothes_store_http_request_duration_seconds_*`, `clothes_store_http_requests_in_flight`. `/metrics` is registered in `backend/cmd/server/main.go`. |
| **Scraping** | Prometheus pulls `http://backend:8080/metrics` inside the Docker network. Config: `monitoring/prometheus.yml`, job `clothes-store-app`. |
| **Host metrics (CPU, etc.)** | Node Exporter (`node-exporter:9100`), job `node-exporter` in the same `prometheus.yml`. |
| **Grafana charts** | Prometheus datasource (`monitoring/grafana/provisioning/datasources/prometheus.yml`, `uid: prometheus`). Dashboard `monitoring/grafana/dashboards/midterm-overview.json` (Golden Signals, SLO, CPU). |

## How to feed data (so charts are not flat)

- **HTTP traffic**: open the site in a browser or call the API, e.g.:
  - `curl http://localhost:8081/ping` (host port **8081** maps to container **8080**)
  - or open `http://localhost` (Nginx → backend).
- **Business data**: use the UI or REST (`/auth/register`, `/api/product`, `/orders`, …) — data is stored in **MongoDB** (service `mongo` in Compose, URI `mongodb://mongo:27017/clothes_store`).
- **Product images**: saved under `frontend/static/assets/products/` (in containers, `FRONTEND_ROOT` points at the bundled frontend tree).

## Local run without Docker

Requirements: Go 1.25+, MongoDB.

Create a `.env` in the repo root (or export variables):

```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017/clothes_store
JWT_SECRET=your_secret
```

Run:

```bash
cd backend && go run ./cmd/server
```

Server: `http://localhost:8080` (or the port from `PORT`).

## Full stack (midterm): Docker Compose

```bash
docker compose up --build -d
```

If Grafana still shows an old datasource and panels look empty:

```bash
docker compose down -v
docker compose up --build -d
```

### Ports (current)

| Service | URL |
|---------|-----|
| Site (Nginx → backend) | http://localhost |
| Backend (host) | http://localhost:8081 |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 (default login: `admin` / `admin`) |
| Node Exporter | http://localhost:9100/metrics |
| MongoDB | localhost:27017 |

Check Prometheus targets: http://localhost:9090/targets — `clothes-store-app` and `node-exporter` should be **UP**.

In Grafana: **Dashboards → Midterm folder → Midterm - Clothes Store Overview** (Golden Signals, SLO availability, p95 latency, CPU).

### Quick traffic for charts

```bash
for i in {1..100}; do curl -s http://localhost:8081/ping >/dev/null; done
```

## API (short)

When calling from the host, use port **8081** for direct backend access.

- `POST /auth/register`, `POST /auth/login` (JSON or form)
- `GET /api/product`, `GET /api/product/:id`
- `POST /orders`, `GET /orders/:id`, …
- Admin: `POST /api/product` (multipart), analytics under `/api/analytics/...`

See handlers and `internal/api/router.go` for full routes.

## Midterm: alerts

Rules: `monitoring/alert_rules.yml`. Example manual trigger (service unavailable for scrape):

```bash
docker compose stop backend
```

After ~1 minute, Prometheus → **Alerts** should show a **FIRING** alert (e.g. `ClothesStoreServiceDown`).

## Code quality

- Layering: handlers → services → repository.
- MongoDB indexes for orders and line items.
- Prometheus metrics via middleware.
