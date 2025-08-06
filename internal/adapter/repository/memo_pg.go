package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/peconote/peconote/internal/domain"
	domainRepo "github.com/peconote/peconote/internal/domain/repository"
)

type memoRepository struct {
	db *sqlx.DB
}

func NewMemoRepository(db *sqlx.DB) domainRepo.MemoRepository {
	return &memoRepository{db: db}
}

func (r *memoRepository) Create(ctx context.Context, m *domain.Memo) error {
	query := `INSERT INTO memo (id, body, tags, created_at, updated_at) VALUES (:id, :body, :tags, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         m.ID,
		"body":       m.Body,
		"tags":       pq.StringArray(m.Tags),
		"created_at": m.CreatedAt,
		"updated_at": m.UpdatedAt,
	})
	return err
}
