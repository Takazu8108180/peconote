package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/peconote/peconote/internal/domain"
)

type mockMemoRepository struct {
	memo      *domain.Memo
	err       error
	listItems []*domain.Memo
	total     int
}

func (m *mockMemoRepository) Create(ctx context.Context, mem *domain.Memo) error {
	m.memo = mem
	return m.err
}

func (m *mockMemoRepository) List(ctx context.Context, tag *string, limit, offset int) ([]*domain.Memo, int, error) {
	return m.listItems, m.total, m.err
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

func TestListMemos_Success(t *testing.T) {
	now := time.Now()
	repo := &mockMemoRepository{listItems: []*domain.Memo{{ID: uuid.New(), Body: "b", CreatedAt: now, UpdatedAt: now}}, total: 1}
	u := NewMemoUsecase(repo)
	items, p, err := u.ListMemos(context.Background(), 1, 20, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || p.TotalCount != 1 {
		t.Fatalf("unexpected result")
	}
}

func TestListMemos_Validation(t *testing.T) {
	repo := &mockMemoRepository{}
	u := NewMemoUsecase(repo)
	if _, _, err := u.ListMemos(context.Background(), 1, 101, nil); !errors.Is(err, ErrInvalidMemoQuery) {
		t.Fatalf("expected validation error")
	}
	tag := ""
	if _, _, err := u.ListMemos(context.Background(), 1, 10, &tag); !errors.Is(err, ErrInvalidMemoQuery) {
		t.Fatalf("expected validation error")
	}
	longTag := "1234567890123456789012345678901"
	if _, _, err := u.ListMemos(context.Background(), 1, 10, &longTag); !errors.Is(err, ErrInvalidMemoQuery) {
		t.Fatalf("expected validation error")
	}
}
