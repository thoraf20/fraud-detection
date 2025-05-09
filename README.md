# Real-Time Fraud Detection System (Go + Redis Streams)

A distributed, event-driven fraud detection system built with Go and Redis Streams. Processes transactions in real-time, applies ML-based rules, and alerts on suspicious activity with exactly-once processing guarantees.

## Features
- **Event-Driven Architecture**: Redis Streams for pub/sub with consumer groups.
- **Production-Ready**: Metrics, structured logging, and dead-letter queues.
- **Scalable**: Horizontal scaling of Go consumers via Redis Consumer Groups.
- **Real-Time ML**: Pre-trained model scoring (ONNX) for anomaly detection.
- **Idempotency**: Deduplication using Redis SETNX.
- **Monitoring**: Prometheus/Grafana dashboards for throughput, fraud rate, and latency.

---

## Architecture

flowchart LR
    A[Payment Gateway] -->|HTTP| B[Producer (Go)]
    B -->|Redis Stream| C[transactions:raw]
    C --> D[Consumer Group 1 (Fraud Check)]
    C --> E[Consumer Group 2 (ML Scoring)]
    D -->|Approved| F[PostgreSQL]
    D -->|Fraud| G[transactions:fraud]
    G --> H[Alert Service (Slack/Email)]
    G --> I[PostgreSQL (Fraud Logs)]
    H --> J[Prometheus]



## Tech Stack

Component	                      Technology
Language	                      Go 1.21+
Event Broker	                  Redis Streams (v7+)
Database	                      PostgreSQL (TimescaleDB)
Caching/State	                  Redis (v7+)
ML Runtime	                    ONNX (Go bindings)
Monitoring	                    Prometheus + Grafana
CI/CD	                          GitHub Actions


## Prerequisites

Go 1.21+

Redis 7.0+ (with stream support)

PostgreSQL 15+

## Quick Start (Local)

# Clone the repo
git clone https://github.com/thoraf20/fraud-detection.git
cd fraud-detection

# Build and run producer + consumers
go run cmd/producer/main.go --stream transactions:raw
go run cmd/consumer/main.go --group fraud-group --stream transactions:raw

## Configuration

REDIS_URL=redis://localhost:6379
POSTGRES_URL=postgres://user:pass@localhost:5432/fraud?sslmode=disable
ML_MODEL_PATH=./models/fraud.onnx
CONSUMER_GROUP=fraud-group
CONSUMER_ID=host-1 # Unique per instance