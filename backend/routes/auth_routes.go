package routes

import (
	"gin-project/controllers"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, container *utils.Container) {
	auth := router.Group("/auth")
	{
		auth.POST("/register",
			//middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.AuthController](container, "Register"))

		auth.POST("/login",
			//middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.AuthController](container, "Login"))

		// Token validation endpoint
		auth.POST("/validate",
			utils.CreateHandlerFunc[controllers.AuthController](container, "ValidateToken"))

		// Token refresh endpoint
		auth.POST("/refresh",
			utils.CreateHandlerFunc[controllers.AuthController](container, "RefreshToken"))

		// Logout endpoint
		auth.POST("/logout",
			utils.CreateHandlerFunc[controllers.AuthController](container, "Logout"))
	}
}
