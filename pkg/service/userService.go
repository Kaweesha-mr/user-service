package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
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
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {

	err := s.UserRepo.CreateUpdateUser(user)
	if err != nil {
		return err
	}

	err = s.redisClient.Del(ctx, "users").Err()
	if err != nil {
		log.Printf("Error clearing cache: %v", err)
	}

	return nil
}
func (s *UserService) GetUserById(ctx context.Context, id string) (model.User, error) {
	// Check cache first
	var user model.User

	// Try fetching from cache
	cachedData, err := s.redisClient.Get(ctx, "user:"+id).Result()
	if err == nil {
		if jsonErr := json.Unmarshal([]byte(cachedData), &user); jsonErr == nil {
			log.Println("Cache hit for user:", id)
			return user, nil
		} else {
			log.Printf("Failed to unmarshal cached data for user %v: %v", id, jsonErr)
		}
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("Error while checking cache for user %v: %v", id, err)
	}

	// Fetch from DB if not in cache or unmarshalling failed
	log.Printf("Cache miss for user %v, fetching from DB", id)
	user, err = s.UserRepo.GetUserById(id)
	if err != nil {
		return model.User{}, err
	}

	// Serialize and cache the result
	data, jsonErr := json.Marshal(user)
	if jsonErr != nil {
		log.Printf("Failed to marshal user %v for caching: %v", id, jsonErr)
	} else {
		// Cache the JSON result
		cacheErr := s.redisClient.Set(ctx, "user:"+id, data, 10*time.Minute).Err()
		if cacheErr != nil {
			log.Printf("Failed to cache user %v: %v", id, cacheErr)
		}
	}

	return user, nil
}
func (s *UserService) UpdateUser(ctx *gin.Context, user *model.User) interface{} {

	err := s.UserRepo.CreateUpdateUser(user)
	if err != nil {
		return err
	}

	err = s.redisClient.Del(ctx, "user:"+user.ID.String()).Err()
	if err != nil {
		log.Printf("Error clearing cache: %v", err)
	}

	return nil
}
func (s *UserService) IsUserAvailable(id string) bool {
	return s.UserRepo.UserAvailable(id)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {

	if err := s.UserRepo.DeleteUser(id); err != nil {
		return err
	}

	cacheKey := "user" + id
	if err := s.redisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Error clearing cache for user ID %s: %v", id, err)
		return err
	}

	return nil
}
