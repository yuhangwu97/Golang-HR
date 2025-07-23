package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupPositionRoutes(router *gin.RouterGroup, container *utils.Container) {
	positions := router.Group("/positions")
	positions.Use(middleware.JWTAuth())
	{
		// 职位基本CRUD
		positions.GET("",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.PositionController](container, "GetPositions"))

		positions.POST("",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.PositionController](container, "CreatePosition"))

		positions.GET("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.PositionController](container, "GetPosition"))

		positions.PUT("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.PositionController](container, "UpdatePosition"))

		positions.DELETE("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			utils.CreateHandlerFunc[controllers.PositionController](container, "DeletePosition"))

		// 职位查询和统计
		positions.GET("/search",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.PositionController](container, "SearchPositions"))

		positions.GET("/statistics",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.PositionController](container, "GetPositionStatistics"))

		// 职位树形结构
		positions.GET("/tree",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.PositionController](container, "GetPositionTree"))

		positions.GET("/all",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.PositionController](container, "GetAllPositions"))

		// 部门相关
		positions.GET("/department/:departmentId",
			middleware.ValidateNumericParam("departmentId"),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.PositionController](container, "GetPositionsByDepartment"))

		// 批量操作
		positions.POST("/bulk",
			middleware.RequireRole("admin"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.PositionController](container, "BulkCreatePositions"))
	}
}
