package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peconote/peconote/internal/domain"
	"github.com/peconote/peconote/internal/domain/model"
	"github.com/peconote/peconote/internal/usecase"
)

type stubMemoUsecase struct {
	id         uuid.UUID
	err        error
	items      []*domain.Memo
	pagination *model.Pagination
	memo       *domain.Memo
}

func (s *stubMemoUsecase) CreateMemo(ctx context.Context, body string, tags []string) (uuid.UUID, error) {
	return s.id, s.err
}

func (s *stubMemoUsecase) ListMemos(ctx context.Context, page, pageSize int, tag *string) ([]*domain.Memo, *model.Pagination, error) {
	return s.items, s.pagination, s.err
}

func (s *stubMemoUsecase) GetMemo(ctx context.Context, id uuid.UUID) (*domain.Memo, error) {
	return s.memo, s.err
}

func (s *stubMemoUsecase) UpdateMemo(ctx context.Context, id uuid.UUID, body string, tags []string) error {
	return s.err
}

func (s *stubMemoUsecase) DeleteMemo(ctx context.Context, id uuid.UUID) error {
	return s.err
}

func TestCreateMemoHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	h := NewMemoHandler(&stubMemoUsecase{id: id})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/memos", bytes.NewBufferString(`{"body":"hi","tags":["t"]}`))
	h.CreateMemo(c)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/api/memos/"+id.String() {
		t.Fatalf("expected location header, got %s", loc)
	}
}

func TestCreateMemoHandler_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMemoHandler(&stubMemoUsecase{err: usecase.ErrInvalidMemo})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/memos", bytes.NewBufferString(`{"body":"","tags":[]}`))
	h.CreateMemo(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestListMemosHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	items := []*domain.Memo{
		{ID: uuid.New(), Body: "a", Tags: []string{"t"}, CreatedAt: now, UpdatedAt: now},
	}
	stub := &stubMemoUsecase{items: items, pagination: &model.Pagination{Page: 1, PageSize: 1, TotalPages: 2, TotalCount: 2}}
	h := NewMemoHandler(stub)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/memos?page=1&page_size=1", nil)
	h.ListMemos(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if link := w.Header().Get("Link"); link != "</api/memos?page=2&page_size=1>; rel=\"next\"" {
		t.Fatalf("unexpected Link header: %s", link)
	}
	var resp MemoListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item")
	}
}

func TestListMemosHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewMemoHandler(&stubMemoUsecase{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/memos?page=0", nil)
	h.ListMemos(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestGetMemoHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	now := time.Now()
	memo := &domain.Memo{ID: id, Body: "hi", Tags: []string{"t"}, CreatedAt: now, UpdatedAt: now}
	h := NewMemoHandler(&stubMemoUsecase{memo: memo})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/memos/"+id.String(), nil)
	h.GetMemo(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	var resp MemoItem
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.ID != id.String() {
		t.Fatalf("unexpected id")
	}
}

func TestGetMemoHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	h := NewMemoHandler(&stubMemoUsecase{err: usecase.ErrMemoNotFound})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/memos/"+id.String(), nil)
	h.GetMemo(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", w.Code)
	}
}

func TestUpdateMemoHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	h := NewMemoHandler(&stubMemoUsecase{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
	c.Request = httptest.NewRequest(http.MethodPut, "/api/memos/"+id.String(), bytes.NewBufferString(`{"body":"hi","tags":["t"]}`))
	h.UpdateMemo(c)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", w.Code)
	}
}

func TestUpdateMemoHandler_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	h := NewMemoHandler(&stubMemoUsecase{err: usecase.ErrInvalidMemo})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
	c.Request = httptest.NewRequest(http.MethodPut, "/api/memos/"+id.String(), bytes.NewBufferString(`{"body":"","tags":[]}`))
	h.UpdateMemo(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestDeleteMemoHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	h := NewMemoHandler(&stubMemoUsecase{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/memos/"+id.String(), nil)
	h.DeleteMemo(c)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", w.Code)
	}
}

func TestDeleteMemoHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.New()
	h := NewMemoHandler(&stubMemoUsecase{err: usecase.ErrMemoNotFound})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: id.String()}}
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/memos/"+id.String(), nil)
	h.DeleteMemo(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", w.Code)
	}
}
