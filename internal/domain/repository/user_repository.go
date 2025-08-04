package repository

import "github.com/peconote/peconote/internal/domain/model"

type UserRepository interface {
	FindAll() ([]model.User, error)
	Create(user *model.User) error
}
