package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupSalaryRoutes(router *gin.RouterGroup, container *utils.Container) {
	salaries := router.Group("/salaries")
	salaries.Use(middleware.JWTAuth())
	{
		// 薪资记录基本CRUD
		salaries.GET("",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetSalaryRecords"))

		salaries.GET("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetSalaryDetail"))

		salaries.POST("",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "CreateSalary"))

		salaries.PUT("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "UpdateSalary"))

		salaries.DELETE("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "DeleteSalary"))

		// 薪资计算
		salaries.POST("/calculate",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "CalculateSalary"))

		salaries.POST("/batch-calculate",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "BatchCalculateSalary"))

		// 薪资审批
		salaries.PUT("/:id/approve",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "ApproveSalary"))

		// 统计和报表
		salaries.GET("/statistics",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetSalaryStatistics"))

		salaries.GET("/export",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "ExportSalaryReport"))

		// 个人薪资查询
		salaries.GET("/my-salary",
			middleware.RequireAnyRole("admin", "hr", "user"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetMySalary"))

		// 个人增强版薪资查询
		salaries.GET("/my-enhanced-salary",
			middleware.RequireAnyRole("admin", "hr", "user"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetMyEnhancedSalary"))

		// 个人薪资历史
		salaries.GET("/my-salary-history",
			middleware.RequireAnyRole("admin", "hr", "user"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetMySalaryHistory"))

		// 个人发放记录
		salaries.GET("/my-payroll-records",
			middleware.RequireAnyRole("admin", "hr", "user"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetMyPayrollRecords"))

		// 个人薪资仪表板
		salaries.GET("/my-dashboard",
			middleware.RequireAnyRole("admin", "hr", "user"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetMySalaryDashboard"))

		// 薪资发放
		salaries.POST("/:id/process-payroll",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "ProcessPayroll"))

		//salaries.GET("/:salary_id/payroll-records",
		//	middleware.ValidateNumericParam("salary_id"),
		//	middleware.RequireAnyRole("admin", "hr"),
		//	utils.CreateHandlerFunc[controllers.SalaryController](container, "GetPayrollRecords"))
	}

	// ========================= Salary Component Management =========================
	components := router.Group("/salary/components")
	components.Use(middleware.JWTAuth())
	{
		components.POST("",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "CreateSalaryComponent"))

		components.GET("",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetSalaryComponents"))

		components.GET("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetSalaryComponent"))

		components.PUT("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "UpdateSalaryComponent"))

		components.DELETE("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "DeleteSalaryComponent"))
	}

	// ========================= Salary Grade Management =========================
	grades := router.Group("/salary/grades")
	grades.Use(middleware.JWTAuth())
	{
		grades.POST("",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "CreateSalaryGrade"))

		grades.GET("",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetSalaryGrades"))
	}

	// ========================= Salary Structure Management =========================
	structures := router.Group("/salary/structures")
	structures.Use(middleware.JWTAuth())
	{
		structures.POST("",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "CreateSalaryStructure"))

		structures.GET("",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetSalaryStructures"))
	}

	// ========================= Payroll Period Management =========================
	periods := router.Group("/payroll/periods")
	periods.Use(middleware.JWTAuth())
	{
		periods.POST("",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "CreatePayrollPeriod"))

		periods.GET("",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetPayrollPeriods"))

		periods.PUT("/:id/lock",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "LockPayrollPeriod"))

		periods.PUT("/:id/unlock",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "UnlockPayrollPeriod"))
	}

	// ========================= Enhanced Salary Management =========================
	enhanced := router.Group("/salary/enhanced")
	enhanced.Use(middleware.JWTAuth())
	{
		// 薪资计算
		enhanced.POST("/calculate",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "CalculateEmployeeSalary"))

		enhanced.POST("/batch-calculate",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "BatchCalculateSalaries"))

		// 薪资记录管理
		enhanced.GET("",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetEnhancedSalaries"))

		enhanced.GET("/:id",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetEnhancedSalary"))

		enhanced.PUT("/:id/details",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "UpdateSalaryDetails"))

		// 审批流程
		enhanced.PUT("/:id/review",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "ReviewSalary"))

		enhanced.PUT("/:id/approve",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "ApproveSalary"))

		enhanced.POST("/bulk-approve",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "BulkApproveSalaries"))
	}

	// ========================= Payment Processing =========================
	payments := router.Group("/payroll/payments")
	payments.Use(middleware.JWTAuth())
	{
		payments.POST("/batches",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "CreatePaymentBatch"))

		payments.GET("/batches",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetPaymentBatches"))

		payments.PUT("/batches/:id/process",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "ProcessPaymentBatch"))
	}

	// ========================= Analytics and Reporting =========================
	analytics := router.Group("/salary/analytics")
	analytics.Use(middleware.JWTAuth())
	{
		analytics.GET("",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetSalaryAnalytics"))

		analytics.GET("/department-report",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "GetDepartmentSalaryReport"))

		analytics.GET("/export",
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "ExportSalaryReport"))
	}

	// ========================= Utilities =========================
	utils_group := router.Group("/salary/utils")
	utils_group.Use(middleware.JWTAuth())
	{
		utils_group.POST("/validate-formula",
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "ValidateFormula"))
	}

	// 薪资发放管理
	payroll := router.Group("/payroll")
	payroll.Use(middleware.JWTAuth())
	{
		payroll.PUT("/:id/status",
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.SalaryController](container, "UpdatePayrollStatus"))
	}
}
