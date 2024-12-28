package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/pkg/model"
	"user-service/pkg/service"
)

type UserController struct {
	UserService *service.UserService
}

// NewUserController creates a new UserController instance
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{UserService: userService}
}

// GetUsers handles GET requests to fetch all users
func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.UserService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// CreateUser handles POST requests to create a new user
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := c.UserService.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}
