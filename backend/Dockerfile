FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /src/backend
RUN go mod download

COPY backend/ /src/backend/
COPY frontend/ /src/frontend/

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.20

RUN adduser -D -H appuser
WORKDIR /app

COPY --from=builder /out/server /app/server
COPY --from=builder /src/frontend /app/frontend

ENV PORT=8080
ENV FRONTEND_ROOT=/app/frontend

EXPOSE 8080

USER appuser

CMD ["/app/server"]
