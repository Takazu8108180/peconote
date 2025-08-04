package usecase

import (
	"github.com/peconote/peconote/internal/domain/model"
	"github.com/peconote/peconote/internal/domain/repository"
)

type UserUsecase interface {
	GetUsers() ([]model.User, error)
	CreateUser(user *model.User) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{repo: r}
}

func (u *userUsecase) GetUsers() ([]model.User, error) {
	return u.repo.FindAll()
}

func (u *userUsecase) CreateUser(user *model.User) error {
	return u.repo.Create(user)
}
