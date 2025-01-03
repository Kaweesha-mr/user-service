package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"user-service/config"
	"user-service/pkg/controllers"
	"user-service/pkg/middleware"
	"user-service/pkg/repository"
	"user-service/pkg/service"
)

func SetUpRouter() *gin.Engine {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to the databse: %v", err)
	}

	redisClient, err := config.ConnectRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	userRepo := repository.NewRepository(db)
	userService := service.NewUserService(userRepo, redisClient)
	middlewareService := service.NewMiddlewareService(redisClient)
	userController := controllers.NewUserController(userService)

	r := gin.Default()
	r.Use(middleware.JWTMiddleware(middlewareService))

	// Health Check route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Service is up and running",
		})
	})

	v1 := r.Group("v1")
	{
		v1.GET("/users", userController.GetUsers)
		v1.GET("/user/:id", userController.GetUserById)
		v1.POST("/users", userController.CreateUser)
		v1.PUT("/update", userController.UpdateUser)
		v1.DELETE("/delete/:id", userController.DeleteUser)

	}

	return r
}
