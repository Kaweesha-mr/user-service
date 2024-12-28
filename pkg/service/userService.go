package service

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
	"user-service/pkg/model"
	"user-service/pkg/repository"
)

type UserService struct {
	UserRepo    *repository.UserRepository
	redisClient *redis.Client
}

func NewUserService(userRepo *repository.UserRepository, redisClient *redis.Client) *UserService {
	return &UserService{
		UserRepo:    userRepo,
		redisClient: redisClient,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	// Check cache first
	users := []model.User{}

	// Try fetching from cache
	cachedData, err := s.redisClient.Get(ctx, "users").Result()
	if err == nil {

		if jsonErr := json.Unmarshal([]byte(cachedData), &users); jsonErr == nil {
			return users, nil
		} else {
			log.Printf("Failed to unmarshal cached data: %v", jsonErr)
		}
	} else if err != redis.Nil {
		log.Printf("Error while checking cache: %v", err)
	}

	// Fetch from DB if not in cache or if unmarshalling failed
	log.Println("Cache miss for users, fetching from DB")
	users, err = s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// Serialize and cache the result
	data, jsonErr := json.Marshal(users)
	if jsonErr != nil {
		log.Printf("Failed to marshal users for caching: %v", jsonErr)
	} else {
		// Cache the JSON result
		if cacheErr := s.redisClient.Set(ctx, "users", data, 10*time.Minute).Err(); cacheErr != nil {
			log.Printf("Failed to cache users: %v", cacheErr)
		}
	}
	return users, nil
}

func (s *UserService) CreateUser(user *model.User) error {
	return s.UserRepo.CreateUser(user)
}
