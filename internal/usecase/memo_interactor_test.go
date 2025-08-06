package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/peconote/peconote/internal/domain"
)

type mockMemoRepository struct {
	memo *domain.Memo
	err  error
}

func (m *mockMemoRepository) Create(ctx context.Context, mem *domain.Memo) error {
	m.memo = mem
	return m.err
}

func TestCreateMemo_Success(t *testing.T) {
	repo := &mockMemoRepository{}
	u := NewMemoUsecase(repo)

	id, err := u.CreateMemo(context.Background(), "hello", []string{"tag"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == uuid.Nil {
		t.Fatalf("expected valid id")
	}
	if repo.memo == nil || repo.memo.Body != "hello" {
		t.Fatalf("memo not saved")
	}
}

func TestCreateMemo_Validation(t *testing.T) {
	repo := &mockMemoRepository{}
	u := NewMemoUsecase(repo)

	_, err := u.CreateMemo(context.Background(), "", nil)
	if !errors.Is(err, ErrInvalidMemo) {
		t.Fatalf("expected validation error")
	}
}
