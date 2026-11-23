# Go Task API

A clean, idiomatic Go REST API for managing tasks.  
Built to demonstrate solid Go engineering practices suitable for client review.

## Highlights

- **Idiomatic Go** — clear package structure, interfaces for storage, proper error handling
- **Concurrent-safe** in-memory store using `sync.RWMutex`
- **Context-aware** — all store methods respect `context.Context`
- **Structured logging** with `log/slog`
- **Middleware** — request logging + panic recovery
- **Graceful shutdown** on SIGINT/SIGTERM
- **Unit tests** for the store layer
- **Modern Go** — uses Go 1.22+ routing patterns (`GET /tasks/{id}`)

## Project Structure

```
go-task-api/
├── cmd/server/          # Application entrypoint
├── internal/
│   ├── handler/         # HTTP handlers + middleware
│   ├── model/           # Domain types + request DTOs
│   └── store/           # Storage interface + in-memory implementation
├── go.mod
└── README.md
```

## Quick Start

```bash
# Run the server
go run ./cmd/server

# Server listens on :8080
```

## API Endpoints

| Method | Path              | Description                  |
|--------|-------------------|------------------------------|
| GET    | `/health`         | Health check                 |
| POST   | `/tasks`          | Create a new task            |
| GET    | `/tasks`          | List tasks (optional filter) |
| GET    | `/tasks/{id}`     | Get a single task            |
| PUT    | `/tasks/{id}`     | Update a task                |
| DELETE | `/tasks/{id}`     | Delete a task                |

### Examples

**Create a task**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go concurrency", "description": "Channels and goroutines"}'
```

**List all tasks**
```bash
curl http://localhost:8080/tasks
```

**List only completed tasks**
```bash
curl "http://localhost:8080/tasks?completed=true"
```

**Update a task**
```bash
curl -X PUT http://localhost:8080/tasks/<id> \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

**Delete a task**
```bash
curl -X DELETE http://localhost:8080/tasks/<id>
```

## Design Decisions

1. **Interface-based store**  
   The `Store` interface allows easy swapping of the persistence layer (e.g. PostgreSQL, Redis) without changing handlers.

2. **In-memory implementation**  
   Chosen for simplicity and zero external dependencies. The concurrent design (`RWMutex`) shows proper care for shared state.

3. **Context propagation**  
   Every store method accepts `context.Context` so cancellation and timeouts work correctly end-to-end.

4. **Thin handlers**  
   Handlers focus on HTTP concerns (validation, status codes, JSON). Business logic lives in the store.

5. **Graceful shutdown**  
   The server handles `SIGINT`/`SIGTERM` cleanly and gives in-flight requests time to finish.

## Running Tests

```bash
go test ./...
```

## Extending This Project

Possible next steps a client might ask for:
- Add PostgreSQL / SQLite persistence
- JWT authentication
- Rate limiting middleware
- OpenAPI / Swagger documentation
- Docker + docker-compose setup
- Observability (metrics + tracing)

---

Built with care using **only the Go standard library**.
