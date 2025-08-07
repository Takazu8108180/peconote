package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/peconote/peconote/internal/domain"
	"github.com/peconote/peconote/internal/domain/model"
	"github.com/peconote/peconote/internal/domain/repository"
)

var ErrInvalidMemo = errors.New("invalid memo")
var ErrInvalidMemoQuery = errors.New("invalid memo query")

type MemoUsecase interface {
	CreateMemo(ctx context.Context, body string, tags []string) (uuid.UUID, error)
	ListMemos(ctx context.Context, page, pageSize int, tag *string) ([]*domain.Memo, *model.Pagination, error)
}

type memoUsecase struct {
	repo repository.MemoRepository
}

func NewMemoUsecase(r repository.MemoRepository) MemoUsecase {
	return &memoUsecase{repo: r}
}

func (u *memoUsecase) CreateMemo(ctx context.Context, body string, tags []string) (uuid.UUID, error) {
	if strings.TrimSpace(body) == "" || len(body) > 2000 {
		return uuid.Nil, ErrInvalidMemo
	}
	if len(tags) > 10 {
		return uuid.Nil, ErrInvalidMemo
	}
	for _, t := range tags {
		if l := len(t); l < 1 || l > 30 {
			return uuid.Nil, ErrInvalidMemo
		}
	}
	id := uuid.New()
	now := time.Now().UTC()
	memo := &domain.Memo{
		ID:        id,
		Body:      body,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := u.repo.Create(ctx, memo); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (u *memoUsecase) ListMemos(ctx context.Context, page, pageSize int, tag *string) ([]*domain.Memo, *model.Pagination, error) {
	if pageSize < 1 || pageSize > 100 {
		return nil, nil, ErrInvalidMemoQuery
	}
	if tag != nil {
		t := strings.TrimSpace(*tag)
		if t == "" || len(t) > 30 {
			return nil, nil, ErrInvalidMemoQuery
		}
		*tag = t
	}
	offset := (page - 1) * pageSize
	items, total, err := u.repo.List(ctx, tag, pageSize, offset)
	if err != nil {
		return nil, nil, err
	}
	totalPages := 0
	if pageSize > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	p := &model.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: total,
	}
	return items, p, nil
}
