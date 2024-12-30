package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type AuthResponse struct {
	Valid bool   `json:"valid"`
	User  string `json:"user"`
}

type MiddlewareService struct {
	redisClient *redis.Client
}

func NewMiddlewareService(redisClient *redis.Client) *MiddlewareService {
	return &MiddlewareService{
		redisClient: redisClient,
	}
}

func (m *MiddlewareService) ValidateTokenWithAuthService(ctx context.Context, token string) (*AuthResponse, error) {
	// Check if the token is already cached
	cachedData, err := m.redisClient.Get(ctx, "auth-token:"+token).Result()
	if err == nil {
		// Token found in cache, unmarshal the cached data
		var cachedResponse AuthResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
			return &cachedResponse, nil
		}
		log.Printf("Failed to unmarshal cached data: %v", err)
	}

	// Prepare the request body with the token
	reqBody, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token: %v", err)
	}

	// Send the request to the authentication service for validation
	resp, err := http.Post("http://localhost:8080/validate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to authentication service: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// If the response is OK, parse the JSON response
	if resp.StatusCode == http.StatusOK {
		var authResponse AuthResponse
		if err := json.Unmarshal(responseBody, &authResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		// Cache the token and user as a JSON string
		cachedData, err := json.Marshal(authResponse)
		if err == nil {
			if cacheErr := m.redisClient.Set(ctx, "auth-token:"+token, cachedData, 10*time.Minute).Err(); cacheErr != nil {
				log.Printf("Failed to cache token: %v", cacheErr)
			}
		} else {
			log.Printf("Failed to marshal auth response for caching: %v", err)
		}

		return &authResponse, nil
	}

	return nil, fmt.Errorf("authentication failed with status code: %d", resp.StatusCode)
}
