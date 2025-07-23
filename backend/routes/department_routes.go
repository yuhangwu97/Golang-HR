package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupDepartmentRoutes(router *gin.RouterGroup, container *utils.Container) {
	departments := router.Group("/departments")
	departments.Use(middleware.JWTAuth())
	{
		// 部门基本CRUD
		departments.GET("",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "GetDepartments"))

		departments.POST("",
			middleware.RequireRole("admin"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "CreateDepartment"))

		departments.GET("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "GetDepartment"))

		departments.PUT("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "UpdateDepartment"))

		departments.DELETE("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "DeleteDepartment"))

		// 部门树形结构
		departments.GET("/tree",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "GetDepartmentTree"))

		// 部门移动
		departments.PUT("/:id/move",
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "MoveDepartment"))

		// 统计信息
		departments.GET("/statistics",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "GetDepartmentStatistics"))

		// 层级管理
		departments.GET("/:id/hierarchy",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "GetDepartmentHierarchy"))

		departments.GET("/:id/subdepartments",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "GetAllSubDepartments"))

		departments.GET("/:id/path",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "GetDepartmentPath"))

		departments.POST("/bulk-sort",
			middleware.RequireRole("admin"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "BulkUpdateDepartmentSort"))

		// 组织架构图
		departments.GET("/chart",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.DepartmentController](container, "GetDepartmentChart"))
	}
}
