# Event Ticket Booking

Auth REST API server written in Go.

## 1. Tech Stack

- Language: Go (1.25+)
- Web framework: Gin
- Database: MySQL (GORM)
- Cache: Redis (access-token blacklist on logout)
- Message Queue: Kafka
- Auth: JWT (access token)

## 2. Getting Started

### Step 1 — Prerequisites

Make sure the following are installed and running before you continue:

- Go >= 1.25
- MySQL on port 3306
- Redis on port 6379
- Kafka on port 9092

### Step 2 — Create the database & tables

Run the schema file (creates the `event_ticket_booking` database and its tables):

```bash
mysql -u root -p < migrations/schema.sql
```

> All SQL statements live in [migrations/schema.sql](migrations/schema.sql).
> The app does not auto-migrate, so this step is required before the first run.

### Step 3 — Configure

**`.env`** — environment variables:

```env
ENVIRONMENT=dev
```

### Step 4 — Run the server

```bash
go mod tidy
go mod vendor
make run
# or
go run cmd/main.go
```

### Step 5 — Run tests

```bash
make test
# or
go test ./...
```

The server runs at `http://localhost:8080`.

### Step 6 — Verify

```bash
curl http://localhost:8080/ping        # {"data":"pong"}
```

## 3. API Testing

A Postman collection is included at [`EventTicketBooking.postman_collection.json`](EventTicketBooking.postman_collection.json).

Import it into Postman via **File → Import** to get all available endpoints with sample requests ready to use.

Test account:

```json
{
  "email": "user@mail.com",
  "password": "12345678"
}
```
