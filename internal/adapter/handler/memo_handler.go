package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peconote/peconote/internal/usecase"
)

type MemoHandler struct {
	usecase usecase.MemoUsecase
}

func NewMemoHandler(u usecase.MemoUsecase) *MemoHandler {
	return &MemoHandler{usecase: u}
}

func (h *MemoHandler) CreateMemo(c *gin.Context) {
	var req MemoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.usecase.CreateMemo(c.Request.Context(), req.Body, req.Tags)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidMemo) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.Header("Location", "/api/memos/"+id.String())
	c.JSON(http.StatusCreated, MemoCreateResponse{ID: id.String()})
}
