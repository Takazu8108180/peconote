# peconote

This repository contains a minimal backend API skeleton built with Go. It adopts the Gin web framework, Gorm ORM, and follows a Clean Architecture with Domain-Driven Design approach.

## Run

```bash
go run ./cmd/api
```

## Structure

- `cmd/api` - Application entry point
- `internal/domain` - Entity and repository interfaces
- `internal/usecase` - Business logic
- `internal/interfaces` - HTTP controllers
- `internal/infrastructure` - Database, router, and repository implementations
## API

### Create Memo

`POST /api/memos`

Request body:

```json
{
  "body": "memo text",
  "tags": ["tag1", "tag2"]
}
```

Example:

```bash
curl -X POST http://localhost:8080/api/memos \
  -H "Content-Type: application/json" \
  -d '{"body":"hello","tags":["sample"]}'
```

Response:

`201 Created`

```json
{"id":"<uuid>"}
```

The `Location` header contains `/api/memos/{id}`.
