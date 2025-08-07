package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peconote/peconote/internal/adapter/handler/util"
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

func (h *MemoHandler) ListMemos(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}
	sizeStr := c.DefaultQuery("page_size", "20")
	pageSize, err := strconv.Atoi(sizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page_size"})
		return
	}
	var tagPtr *string
	if tag, ok := c.GetQuery("tag"); ok {
		if tag == "" || len(tag) > 30 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag"})
			return
		}
		tagPtr = &tag
	}

	items, pagination, err := h.usecase.ListMemos(c.Request.Context(), page, pageSize, tagPtr)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidMemoQuery) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	resItems := make([]MemoItem, len(items))
	for i, m := range items {
		resItems[i] = MemoItem{
			ID:        m.ID.String(),
			Body:      m.Body,
			Tags:      m.Tags,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	resp := MemoListResponse{Items: resItems, Pagination: *pagination}
	if link := util.BuildLinkHeader("/api/memos", resp.Pagination, tagPtr); link != "" {
		c.Header("Link", link)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MemoHandler) GetMemo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	memo, err := h.usecase.GetMemo(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrMemoNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.JSON(http.StatusOK, MemoItem{
		ID:        memo.ID.String(),
		Body:      memo.Body,
		Tags:      memo.Tags,
		CreatedAt: memo.CreatedAt,
		UpdatedAt: memo.UpdatedAt,
	})
}

func (h *MemoHandler) UpdateMemo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req MemoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.UpdateMemo(c.Request.Context(), id, req.Body, req.Tags); err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidMemo):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrMemoNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *MemoHandler) DeleteMemo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.DeleteMemo(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, usecase.ErrMemoNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
