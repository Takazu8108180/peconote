package domain

import (
	"time"

	"github.com/google/uuid"
)

type Memo struct {
	ID        uuid.UUID
	Body      string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
