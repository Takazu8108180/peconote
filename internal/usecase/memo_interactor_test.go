package usecase

import (
	"context"
	"database/sql"
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

func (m *mockMemoRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Memo, error) {
	return m.memo, m.err
}

func (m *mockMemoRepository) Update(ctx context.Context, memo *domain.Memo) error {
	m.memo = memo
	return m.err
}

func (m *mockMemoRepository) Delete(ctx context.Context, id uuid.UUID) error {
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

func TestGetMemo_Success(t *testing.T) {
	now := time.Now()
	memo := &domain.Memo{ID: uuid.New(), Body: "b", CreatedAt: now, UpdatedAt: now}
	repo := &mockMemoRepository{memo: memo}
	u := NewMemoUsecase(repo)
	got, err := u.GetMemo(context.Background(), memo.ID)
	if err != nil || got.ID != memo.ID {
		t.Fatalf("unexpected result")
	}
}

func TestGetMemo_NotFound(t *testing.T) {
	repo := &mockMemoRepository{err: sql.ErrNoRows}
	u := NewMemoUsecase(repo)
	if _, err := u.GetMemo(context.Background(), uuid.New()); !errors.Is(err, ErrMemoNotFound) {
		t.Fatalf("expected not found")
	}
}

func TestUpdateMemo_Validation(t *testing.T) {
	repo := &mockMemoRepository{}
	u := NewMemoUsecase(repo)
	if err := u.UpdateMemo(context.Background(), uuid.New(), "", nil); !errors.Is(err, ErrInvalidMemo) {
		t.Fatalf("expected validation error")
	}
}

func TestDeleteMemo_NotFound(t *testing.T) {
	repo := &mockMemoRepository{err: sql.ErrNoRows}
	u := NewMemoUsecase(repo)
	if err := u.DeleteMemo(context.Background(), uuid.New()); !errors.Is(err, ErrMemoNotFound) {
		t.Fatalf("expected not found")
	}
}
