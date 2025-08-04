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

