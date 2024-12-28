package service

import (
	"user-service/pkg/model"
	"user-service/pkg/repository"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.UserRepo.GetAllUsers()
}

func (s *UserService) CreateUser(user *model.User) error {
	return s.UserRepo.CreateUser(user)
}
