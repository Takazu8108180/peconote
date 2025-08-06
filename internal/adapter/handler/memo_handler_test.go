package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peconote/peconote/internal/usecase"
)

type stubMemoUsecase struct {
	id  uuid.UUID
	err error
}

func (s *stubMemoUsecase) CreateMemo(ctx context.Context, body string, tags []string) (uuid.UUID, error) {
	return s.id, s.err
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
