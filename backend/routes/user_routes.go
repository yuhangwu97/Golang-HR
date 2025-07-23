package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, container *utils.Container) {
	users := router.Group("/users")
	users.Use(middleware.JWTAuth())
	{
		users.POST("", 
			middleware.RequireRole("admin"), 
			middleware.ValidateAndBindJSON(), 
			utils.CreateHandlerFunc[controllers.UserController](container, "CreateUser"))
		
		users.GET("", 
			middleware.RequireAnyRole("admin", "user"), 
			utils.CreateHandlerFunc[controllers.UserController](container, "GetUsers"))
		
		users.GET("/:id", 
			middleware.ValidateObjectID(), 
			middleware.RequireAnyRole("admin", "user"), 
			utils.CreateHandlerFunc[controllers.UserController](container, "GetUser"))
		
		users.PUT("/:id", 
			middleware.ValidateObjectID(), 
			middleware.RequireAnyRole("admin", "user"), 
			middleware.ValidateAndBindJSON(), 
			utils.CreateHandlerFunc[controllers.UserController](container, "UpdateUser"))
		
		users.DELETE("/:id", 
			middleware.ValidateObjectID(), 
			middleware.RequireRole("admin"), 
			utils.CreateHandlerFunc[controllers.UserController](container, "DeleteUser"))
	}
}