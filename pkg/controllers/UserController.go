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

func (c *UserController) GetUserById(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, err := c.UserService.GetUserById(ctx, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user"})
		return
	}

	ctx.JSON(http.StatusCreated, user)

}

func (c *UserController) UpdateUser(ctx *gin.Context) {
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

	if err := c.UserService.UpdateUser(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	ctx.JSON(http.StatusCreated, user)

}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("id") // Retrieve user ID from the URL parameter

	// Check if user exists
	if !c.UserService.IsUserAvailable(userId) {
		ctx.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Attempt to delete the user
	if err := c.UserService.DeleteUser(ctx, userId); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	// Return success response
	ctx.JSON(200, gin.H{"message": "User deleted successfully"})
}
