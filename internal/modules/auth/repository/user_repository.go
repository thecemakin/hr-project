package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/auth/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
}

type sqlUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &sqlUserRepository{db: db}
}

func (r *sqlUserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqlUserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqlUserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *sqlUserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
