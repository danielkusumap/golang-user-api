package routes

import (
	"golang-api/handlers"
	"golang-api/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigureUserRoutes(router *gin.Engine, client *mongo.Client) {
	userRouter := router.Group("/users")

	userHandlers := handlers.Handler{Client: client}

	userRouter.POST("/register", userHandlers.RegisterUser)
	userRouter.POST("/login", userHandlers.LoginUser)
	userRouter.GET("/all", middleware.AuthMiddleware(), userHandlers.GetAllUsers)
	userRouter.POST("/logout", middleware.AuthMiddleware(), userHandlers.LogoutUser)
}
