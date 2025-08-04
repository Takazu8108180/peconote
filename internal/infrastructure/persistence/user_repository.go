package persistence

import (
	"github.com/peconote/peconote/internal/domain/model"
	domainRepository "github.com/peconote/peconote/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainRepository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}
