package repository

import (
	"gorm.io/gorm"
	"user-service/pkg/model"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserepository(db *gorm.DB) *UserRepository {

	return &UserRepository{DB: db}

}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser saves a new user into the database
func (r *UserRepository) CreateUser(user *model.User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}
