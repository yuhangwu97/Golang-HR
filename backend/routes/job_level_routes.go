package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupJobLevelRoutes(router *gin.RouterGroup, container *utils.Container) {
	jobLevels := router.Group("/job_levels")
	jobLevels.Use(middleware.JWTAuth())
	{
		// 职级基本CRUD
		jobLevels.GET("", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.JobLevelController](container, "GetJobLevels"))
		
		jobLevels.POST("", 
			middleware.RequireAnyRole("admin", "hr"), 
			middleware.ValidateAndBindJSON(), 
			utils.CreateHandlerFunc[controllers.JobLevelController](container, "CreateJobLevel"))
		
		jobLevels.GET("/:id", 
			middleware.ValidateNumericID(), 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.JobLevelController](container, "GetJobLevel"))
		
		jobLevels.PUT("/:id", 
			middleware.ValidateNumericID(), 
			middleware.RequireAnyRole("admin", "hr"), 
			middleware.ValidateAndBindJSON(), 
			utils.CreateHandlerFunc[controllers.JobLevelController](container, "UpdateJobLevel"))
		
		jobLevels.DELETE("/:id", 
			middleware.ValidateNumericID(), 
			middleware.RequireRole("admin"), 
			utils.CreateHandlerFunc[controllers.JobLevelController](container, "DeleteJobLevel"))

		// 职级查询和统计
		jobLevels.GET("/by-level", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.JobLevelController](container, "GetJobLevelsByLevel"))
		
		jobLevels.GET("/statistics", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.JobLevelController](container, "GetJobLevelStatistics"))

		// 薪资范围相关
		jobLevels.POST("/validate-salary", 
			middleware.RequireAnyRole("admin", "hr"), 
			middleware.ValidateAndBindJSON(), 
			utils.CreateHandlerFunc[controllers.JobLevelController](container, "ValidateSalaryRange"))
		
		jobLevels.GET("/level/:level/salary-range", 
			middleware.ValidateNumericParam("level"),
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.JobLevelController](container, "GetSalaryRangeByLevel"))
	}
}