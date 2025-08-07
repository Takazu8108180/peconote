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

### List Memos

`GET /api/memos`

Query parameters:

- `page` (default `1`)
- `page_size` (default `20`, max `100`)
- `tag` (optional)

Example:

```bash
curl "http://localhost:8080/api/memos?page=2&page_size=10"
```

Response:

```json
{
  "items": [{"id": "<uuid>", "body": "memo", "tags": ["tag"], "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}],
"pagination": {"page": 2, "page_size": 10, "total_pages": 5, "total_count": 50}
}
```

### Get Memo

`GET /api/memos/{id}`

Example:

```bash
curl http://localhost:8080/api/memos/<id>
```

Response:

```json
{"id": "<uuid>", "body": "memo", "tags": ["tag"], "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}
```

### Update Memo

`PUT /api/memos/{id}`

Request body:

```json
{
  "body": "updated text",
  "tags": ["tag1", "tag2"]
}
```

Example:

```bash
curl -X PUT http://localhost:8080/api/memos/<id> \
  -H "Content-Type: application/json" \
  -d '{"body":"updated","tags":["sample"]}'
```

Response: `204 No Content`

### Delete Memo

`DELETE /api/memos/{id}`

Example:

```bash
curl -X DELETE http://localhost:8080/api/memos/<id>
```

Response: `204 No Content`
