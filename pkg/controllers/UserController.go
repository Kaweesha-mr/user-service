package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/pkg/model"
	"user-service/pkg/service"
	"user-service/pkg/utils"
)

type UserController struct {
	UserService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.UserService.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Use the reusable validation module
	validationErrors, valid := utils.ValidateStruct(user)
	if !valid {
		// If validation fails, return validation errors
		ctx.JSON(http.StatusBadRequest, gin.H{"validation_errors": validationErrors})
		return
	}

	if err := c.UserService.CreateUser(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}
