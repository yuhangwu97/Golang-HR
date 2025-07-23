package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupEmployeeRoutes(router *gin.RouterGroup, container *utils.Container) {
	employees := router.Group("/employees")
	employees.Use(middleware.JWTAuth())
	{
		// 员工基本CRUD
		employees.GET("", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "GetEmployees"))
		
		employees.POST("", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "CreateEmployee"))
		
		employees.GET("/:id", 
			middleware.ValidateNumericID(), 
			middleware.RequireAnyRole("admin", "hr", "user"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "GetEmployee"))
		
		employees.PUT("/:id", 
			middleware.ValidateNumericID(), 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "UpdateEmployee"))
		
		employees.DELETE("/:id", 
			middleware.ValidateNumericID(), 
			middleware.RequireRole("admin"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "DeleteEmployee"))

		// 员工搜索和统计
		employees.GET("/search", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "SearchEmployees"))
		
		employees.GET("/statistics", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "GetEmployeeStatistics"))

		// 部门相关
		employees.GET("/department/:departmentId", 
			middleware.ValidateNumericParam("departmentId"),
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "GetEmployeesByDepartment"))

		// 导入导出
		employees.GET("/export", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "ExportEmployees"))
		
		employees.POST("/import", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "ImportEmployees"))

		// 批量操作
		employees.PATCH("/bulk", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "BulkUpdateEmployees"))
		
		employees.DELETE("/bulk", 
			middleware.RequireRole("admin"), 
			utils.CreateHandlerFunc[controllers.EmployeeController](container, "BulkDeleteEmployees"))
	}
}
