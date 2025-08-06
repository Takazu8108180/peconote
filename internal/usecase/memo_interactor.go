package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/peconote/peconote/internal/domain"
	"github.com/peconote/peconote/internal/domain/repository"
)

var ErrInvalidMemo = errors.New("invalid memo")

type MemoUsecase interface {
	CreateMemo(ctx context.Context, body string, tags []string) (uuid.UUID, error)
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
