# WhatsApp Messaging API

Go + Gin + MySQL + RabbitMQ + Worker. **Mock mode only: no real WhatsApp messages are sent.**

## Local setup

Requirements: Go 1.23+, Docker Compose.

```bash
docker compose up -d
cp .env.example .env
set -a; source .env; set +a
go mod tidy
docker compose exec -T mysql mysql -umessaging -pmessaging_local messaging < migrations/001_init.sql
docker compose exec -T mysql mysql -umessaging -pmessaging_local messaging < migrations/002_outbox.sql
```

Run the API and worker in separate terminals (load the environment variables in both):

```bash
go run ./cmd/api
go run ./cmd/worker
```

## Example

```bash
curl -X POST localhost:8080/v1/companies -H 'Content-Type: application/json' -d '{"name":"Empresa Demo"}'
curl -X POST localhost:8080/v1/contacts -H 'Content-Type: application/json' -d '{"company_id":1,"name":"Ana","phone":"+5511999999999","opted_in":true}'
curl -X POST localhost:8080/v1/campaigns -H 'Content-Type: application/json' -d '{"company_id":1,"name":"Teste","template_name":"lembrete_pagamento"}'
curl -X POST localhost:8080/v1/campaigns/1/send
curl localhost:8080/v1/campaigns/1/stats
```

The worker marks messages as SENT **in simulation only**.

## Architecture

The send endpoint inserts messages and outbox records in one MySQL transaction. The outbox publisher publishes durable RabbitMQ messages with publisher confirms. The worker consumes jobs and conditionally updates the database.

## Not production ready

No authentication, tenant isolation, audit-grade opt-in, WhatsApp Cloud API, webhook, rate limiting, opt-out, bounded retry, DLQ, or full integration tests. Publisher concurrency and external idempotency need hardening. Never commit real credentials to the repository.
