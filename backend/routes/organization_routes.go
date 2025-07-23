package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupOrganizationRoutes(router *gin.RouterGroup, container *utils.Container) {
	org := router.Group("/organization")
	org.Use(middleware.JWTAuth())
	{
		// 组织单元基础 API
		setupUnitRoutes(org, container)
		
		// 组织架构结构 API
		setupStructureRoutes(org, container)
		
		// 员工分配管理 API
		setupAssignmentRoutes(org, container)
		
		// 组织变更管理 API
		setupChangeRoutes(org, container)
		
		// 历史记录 API
		setupHistoryRoutes(org, container)
		
		// 统计分析 API
		setupAnalyticsRoutes(org, container)
		
		// 工具类 API
		setupUtilityRoutes(org, container)
	}
}

// 组织单元基础路由
func setupUnitRoutes(router *gin.RouterGroup, container *utils.Container) {
	units := router.Group("/units")
	{
		units.GET("", 
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetOrganizationUnits"))
		
		units.POST("", 
			middleware.RequireRole("admin"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "CreateOrganizationUnit"))
		
		units.GET("/:id", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetOrganizationUnit"))
		
		units.PUT("/:id", 
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "UpdateOrganizationUnit"))
		
		units.DELETE("/:id", 
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "DeleteOrganizationUnit"))
		
		// 单元相关操作
		units.PUT("/:id/move", 
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "MoveUnit"))
		
		units.GET("/:id/path", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetHierarchyPath"))
		
		units.GET("/:id/children", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetSubunits"))
		
		units.GET("/:id/assignments", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetUnitAssignments"))
		
		units.GET("/:id/statistics", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetUnitStatistics"))
		
		units.GET("/:id/history", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetUnitHistory"))
		
		units.GET("/:id/history/compare", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "CompareHistoryVersions"))
		
		units.PUT("/:id/rollback", 
			middleware.ValidateNumericID(),
			middleware.RequireRole("admin"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "RollbackToHistory"))
		
		units.GET("/:id/evolution", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetOrganizationEvolution"))
	}
}

// 组织架构结构路由
func setupStructureRoutes(router *gin.RouterGroup, container *utils.Container) {
	router.GET("/tree", 
		middleware.RequireAnyRole("admin", "hr"),
		utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetOrganizationTree"))
	
	router.GET("/search", 
		middleware.RequireAnyRole("admin", "hr"),
		utils.CreateHandlerFunc[controllers.OrganizationController](container, "SearchUnits"))
	
	router.GET("/timeline", 
		middleware.RequireAnyRole("admin", "hr"),
		utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetOrganizationTimeline"))
}

// 员工分配管理路由
func setupAssignmentRoutes(router *gin.RouterGroup, container *utils.Container) {
	assignments := router.Group("/assignments")
	{
		assignments.POST("", 
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "AssignEmployee"))
		
		assignments.PUT("/:id", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "UpdateAssignment"))
		
		assignments.DELETE("/:id", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "RemoveAssignment"))
	}
}

// 组织变更管理路由
func setupChangeRoutes(router *gin.RouterGroup, container *utils.Container) {
	changes := router.Group("/changes")
	{
		changes.GET("", 
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetOrganizationChanges"))
		
		changes.POST("", 
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "CreateOrganizationChange"))
		
		changes.PUT("/:id/approve", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "ApproveChange"))
		
		changes.PUT("/:id/reject", 
			middleware.ValidateNumericID(),
			middleware.RequireAnyRole("admin", "hr"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "RejectChange"))
	}
}

// 历史记录路由
func setupHistoryRoutes(router *gin.RouterGroup, container *utils.Container) {
	snapshots := router.Group("/snapshots")
	{
		snapshots.GET("", 
			middleware.RequireAnyRole("admin", "hr"),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetOrganizationSnapshots"))
		
		snapshots.POST("", 
			middleware.RequireRole("admin"),
			middleware.ValidateAndBindJSON(),
			utils.CreateHandlerFunc[controllers.OrganizationController](container, "CreateOrganizationSnapshot"))
	}
}

// 统计分析路由
func setupAnalyticsRoutes(router *gin.RouterGroup, container *utils.Container) {
	// 变更统计 - 匹配前端期望的路径
	router.GET("/statistics/changes", 
		middleware.RequireAnyRole("admin", "hr"),
		utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetChangeStatistics"))
}

// 工具类路由
func setupUtilityRoutes(router *gin.RouterGroup, container *utils.Container) {
	// 系统配置
	router.GET("/types", 
		middleware.RequireAnyRole("admin", "hr"),
		utils.CreateHandlerFunc[controllers.OrganizationController](container, "GetUnitTypes"))
	
	// 验证工具
	router.GET("/validate-code", 
		middleware.RequireAnyRole("admin", "hr"),
		utils.CreateHandlerFunc[controllers.OrganizationController](container, "ValidateUnitCode"))
	
	// 批量操作
	router.PUT("/batch", 
		middleware.RequireRole("admin"),
		middleware.ValidateAndBindJSON(),
		utils.CreateHandlerFunc[controllers.OrganizationController](container, "BatchUpdateUnits"))
	
	// 导出功能
	router.GET("/export", 
		middleware.RequireAnyRole("admin", "hr"),
		utils.CreateHandlerFunc[controllers.OrganizationController](container, "ExportOrganization"))
}