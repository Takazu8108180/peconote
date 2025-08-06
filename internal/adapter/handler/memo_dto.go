package handler

type MemoCreateRequest struct {
	Body string   `json:"body" binding:"required,max=2000"`
	Tags []string `json:"tags" binding:"max=10"`
}

type MemoCreateResponse struct {
	ID string `json:"id"`
}
