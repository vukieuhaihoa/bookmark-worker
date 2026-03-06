# bookmark-worker

A background worker service that processes bookmark import jobs from a Redis queue, persists them to PostgreSQL, and invalidates the Redis cache so the API service always serves fresh data. Built with Go using Domain-Driven Design (DDD) principles.

## Architecture

```
Redis Queue
    │
    ▼
Engine (main loop)
    │  pops messages
    ▼
Worker Pool (5 concurrent workers)
    │  dispatches messages
    ▼
Handler → Service ──► 1. Invalidate Redis cache (list_bookmarks_{user_id})
                  └──► 2. Insert bookmarks → PostgreSQL
```

- **Engine**: polls Redis queue in a loop, dispatches messages to the worker pool
- **Worker Pool**: fixed pool of goroutines that process messages concurrently
- **Handler**: deserializes the JSON message and calls the bookmark service
- **Service**: invalidates the user's bookmark list cache, then persists bookmarks
- **Cache Repository**: deletes the `list_bookmarks_{user_id}` key from Redis
- **Bookmark Repository**: inserts bookmarks into PostgreSQL via GORM

## Message Format

Workers consume JSON messages from Redis with the following structure:

```json
{
  "user_id": "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
  "bookmarks": [
    {
      "url": "https://example.com/bookmark1",
      "description": "My bookmark"
    }
  ]
}
```

## Tech Stack

| Component  | Technology               |
|------------|--------------------------|
| Language   | Go 1.25                  |
| Queue      | Redis (go-redis v9)      |
| Cache      | Redis (go-redis v9)      |
| Database   | PostgreSQL (GORM)        |
| Logging    | zerolog                  |
| Container  | Docker                   |

## Configuration

All configuration is via environment variables:

| Variable       | Default                   | Description                          |
|----------------|---------------------------|--------------------------------------|
| `QUEUE_NAME`   | `bookmark_import_queue`   | Redis queue key to consume from      |
| `SERVICE_NAME` | `bookmark-worker`         | Service identifier for logging       |
| `INSTANCE_ID`  | _(auto-generated UUID)_   | Unique instance identifier           |
| `LOG_LEVEL`    | `debug`                   | Log level (`debug`, `info`, `error`) |
| `DB_*`         | _(see bookmark-libs)_     | PostgreSQL connection settings       |
| `REDIS_*`      | _(see bookmark-libs)_     | Redis connection settings            |

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Redis
- PostgreSQL

### Run locally

Start dependencies:

```bash
make dev-up
```

Run the worker:

```bash
make dev-run
```

### Run with Docker

```bash
make docker-build
make docker-up
```

Stop:

```bash
make docker-down
```

## Development

### Generate mocks

```bash
make mock-gen
```

### Run tests

```bash
make test
```

Tests require a minimum of **80% code coverage**. Coverage report is generated at `./coverage/coverage.html`.

### Run tests in Docker

```bash
make docker-test
```

### Database migrations

Create a new migration:

```bash
make new-schema name=add_column_to_bookmarks
```

Run migrations:

```bash
make migrate
```

### Redis utilities

```bash
make redis-run      # start a local Redis container
make redis-cli      # open redis-cli
make redis-monitor  # monitor Redis commands in real time
```

## Project Structure

```
bookmark-worker/
├── cmd/
│   └── worker/             # entrypoint
├── internal/
│   ├── app/
│   │   ├── handler/        # message deserialization
│   │   ├── model/          # domain models
│   │   ├── repository/
│   │   │   ├── bookmark/   # PostgreSQL bookmark repository
│   │   │   ├── cache/      # Redis cache repository
│   │   │   └── queue/      # Redis queue repository
│   │   └── service/        # business logic (cache invalidation + DB write)
│   ├── infrastructure/     # wiring (DB, Redis, engine setup)
│   └── worker/             # engine, worker pool, config
├── migrations/             # SQL migration files
└── test/
    └── integration/        # integration tests
```
