# Event Ticket Booking

Auth REST API server written in Go.

## 1. Tech Stack

- **Language:** Go (1.25+)
- **Web framework:** Gin
- **Database:** MySQL (GORM)
- **Cache:** Redis (access-token blacklist on logout)
- **Auth:** JWT (access token)

## 2. Prerequisites

- Go >= 1.25
- MySQL >= 8.0 (running)
- Redis (running)

## 3. Getting Started

### Step 1 — Create the database & tables

Run the schema file (creates the `event_ticket_booking` database and its tables):

```bash
mysql -u root -p < migrations/schema.sql
```

> All SQL statements live in [migrations/schema.sql](migrations/schema.sql).
> The app does **not** auto-migrate, so this step is required before the first run.

### Step 2 — Configure

**`.env`** — environment variables:

```env
ENVIRONMENT=dev
TRUSTED_PROXIES=127.0.0.1
ACCESS_SECRET=<jwt-signing-secret>   # change this before production
```

**`config/config-dev.json`** — match your local MySQL/Redis:

```jsonc
{
  "Server": { "Port": "8080" },
  "Db": {
    "Username": "root",
    "Password": "root",
    "Host": "127.0.0.1",
    "Port": "3306",
    "DbName": "event_ticket_booking"   // must match the DB name from Step 1
  },
  "Authentication": {
    "AccessTokenExpirationMinutes": 5
  },
  "Redis": { "Host": "127.0.0.1", "Port": "6379", "Username": "", "Password": "", "Db": 0 }
}
```

### Step 3 — Run the server

```bash
make run
# or
go run cmd/main.go
```

The server runs at `http://localhost:8080`. Dependencies are vendored (`vendor/`), so no network access is needed to build.

### Step 4 — Verify

```bash
curl http://localhost:8080/ping        # {"data":"pong"}
```

## 4. API

| Method | Endpoint           | Description                       |
|--------|--------------------|-----------------------------------|
| GET    | `/ping`            | Health check                      |
| POST   | `/signup`          | Register a new account            |
| POST   | `/login`           | Log in                            |
| POST   | `/logout`          | Log out (blacklist the token)     |

Examples:

```bash
# Sign up
curl -X POST http://localhost:8080/signup \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"Password@123"}'

# Log in
curl -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"Password@123"}'
```

Every response follows the shape: `{ "data": ..., "error": { "message": ... } }`.
