package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
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

func (r *memoRepository) List(ctx context.Context, tag *string, limit, offset int) ([]*domain.Memo, int, error) {
	type memoRow struct {
		ID        uuid.UUID      `db:"id"`
		Body      string         `db:"body"`
		Tags      pq.StringArray `db:"tags"`
		CreatedAt time.Time      `db:"created_at"`
		UpdatedAt time.Time      `db:"updated_at"`
	}

	var rows []memoRow
	query := `SELECT id, body, tags, created_at, updated_at
FROM memo
WHERE ($1::text IS NULL OR $1 = ANY(tags))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`
	if err := r.db.SelectContext(ctx, &rows, query, tag, limit, offset); err != nil {
		return nil, 0, err
	}
	memos := make([]*domain.Memo, len(rows))
	for i, row := range rows {
		memos[i] = &domain.Memo{
			ID:        row.ID,
			Body:      row.Body,
			Tags:      []string(row.Tags),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
	}
	var total int
	countQuery := `SELECT COUNT(*) FROM memo WHERE ($1::text IS NULL OR $1 = ANY(tags))`
	if err := r.db.GetContext(ctx, &total, countQuery, tag); err != nil {
		return nil, 0, err
	}
	return memos, total, nil
}

func (r *memoRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Memo, error) {
	type memoRow struct {
		ID        uuid.UUID      `db:"id"`
		Body      string         `db:"body"`
		Tags      pq.StringArray `db:"tags"`
		CreatedAt time.Time      `db:"created_at"`
		UpdatedAt time.Time      `db:"updated_at"`
	}
	var row memoRow
	query := `SELECT id, body, tags, created_at, updated_at FROM memo WHERE id = $1`
	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		return nil, err
	}
	return &domain.Memo{
		ID:        row.ID,
		Body:      row.Body,
		Tags:      []string(row.Tags),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *memoRepository) Update(ctx context.Context, m *domain.Memo) error {
	query := `UPDATE memo SET body = :body, tags = :tags, updated_at = :updated_at WHERE id = :id`
	res, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         m.ID,
		"body":       m.Body,
		"tags":       pq.StringArray(m.Tags),
		"updated_at": m.UpdatedAt,
	})
	if err != nil {
		return err
	}
	if cnt, err := res.RowsAffected(); err == nil && cnt == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *memoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM memo WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cnt, err := res.RowsAffected(); err == nil && cnt == 0 {
		return sql.ErrNoRows
	}
	return nil
}
