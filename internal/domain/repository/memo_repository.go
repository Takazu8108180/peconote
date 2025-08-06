package repository

import (
	"context"

	"github.com/peconote/peconote/internal/domain"
)

type MemoRepository interface {
	Create(ctx context.Context, m *domain.Memo) error
}
