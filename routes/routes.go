package routes

import (
	"parliament-analytic-be/controllers"
	"parliament-analytic-be/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Auth routes
	authGroup := r.Group("/auth")
	authGroup.POST("/register", controllers.Register)
	authGroup.POST("/login", controllers.Login)

	// Admin routes (with middleware)
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())

	adminGroup.GET("/partai", controllers.GetAllPartai)
	adminGroup.GET("/partai/:id", controllers.GetPartaiByID)
	adminGroup.POST("/partai", controllers.CreatePartai)
	adminGroup.PUT("/partai/:id", controllers.UpdatePartai)
	adminGroup.DELETE("/partai/:id", controllers.DeletePartai)
}
