package repository

import (
	"gorm.io/gorm"
	"user-service/pkg/model"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *UserRepository {

	return &UserRepository{DB: db}

}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUpdateUser saves a new user into the database
func (r *UserRepository) CreateUpdateUser(user *model.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserById(id string) (model.User, error) {
	var user model.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepository) UserAvailable(id string) bool {
	var user model.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return true
	}
	return false
}
func (r *UserRepository) DeleteUser(id string) error {
	// Delete the user from the database
	if err := r.DB.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}
