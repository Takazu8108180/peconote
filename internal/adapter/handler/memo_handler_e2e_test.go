package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peconote/peconote/internal/domain"
	"github.com/peconote/peconote/internal/usecase"
)

type memoryMemoRepo struct {
	memos []*domain.Memo
}

func (m *memoryMemoRepo) Create(ctx context.Context, memo *domain.Memo) error {
	m.memos = append(m.memos, memo)
	return nil
}

func (m *memoryMemoRepo) List(ctx context.Context, tag *string, limit, offset int) ([]*domain.Memo, int, error) {
	filtered := make([]*domain.Memo, 0, len(m.memos))
	for _, me := range m.memos {
		if tag != nil {
			ok := false
			for _, t := range me.Tags {
				if t == *tag {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		filtered = append(filtered, me)
	}
	total := len(filtered)
	end := offset + limit
	if end > total {
		end = total
	}
	if offset > total {
		return []*domain.Memo{}, total, nil
	}
	return filtered[offset:end], total, nil
}

func TestListMemos_E2E(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &memoryMemoRepo{}
	now := time.Now()
	for i := 0; i < 30; i++ {
		repo.memos = append(repo.memos, &domain.Memo{
			ID:        uuid.New(),
			Body:      fmt.Sprintf("memo %d", i),
			Tags:      []string{"t"},
			CreatedAt: now.Add(-time.Duration(i) * time.Minute),
			UpdatedAt: now.Add(-time.Duration(i) * time.Minute),
		})
	}
	u := usecase.NewMemoUsecase(repo)
	h := NewMemoHandler(u)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/memos?page=2&page_size=10", nil)
	h.ListMemos(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if link := w.Header().Get("Link"); link == "" {
		t.Fatalf("expected Link header")
	}
	var resp MemoListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(resp.Items) != 10 {
		t.Fatalf("expected 10 items")
	}
	if !resp.Items[0].CreatedAt.After(resp.Items[9].CreatedAt) {
		t.Fatalf("items not sorted")
	}
}
