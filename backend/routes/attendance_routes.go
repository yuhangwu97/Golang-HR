package routes

import (
	"gin-project/controllers"
	"gin-project/middleware"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupAttendanceRoutes(router *gin.RouterGroup, container *utils.Container) {
	attendance := router.Group("/attendance")
	attendance.Use(middleware.JWTAuth())
	{
		// 考勤打卡
		attendance.POST("/checkin", 
			middleware.RequireAnyRole("admin", "hr", "user"), 
			middleware.ValidateAndBindJSON(), 
			utils.CreateHandlerFunc[controllers.AttendanceController](container, "CheckIn"))
		
		attendance.POST("/checkout", 
			middleware.RequireAnyRole("admin", "hr", "user"), 
			middleware.ValidateAndBindJSON(), 
			utils.CreateHandlerFunc[controllers.AttendanceController](container, "CheckOut"))

		// 考勤记录查询
		attendance.GET("/records", 
			middleware.RequireAnyRole("admin", "hr", "user"), 
			utils.CreateHandlerFunc[controllers.AttendanceController](container, "GetAttendanceRecords"))
		
		attendance.GET("/today", 
			middleware.RequireAnyRole("admin", "hr", "user"), 
			utils.CreateHandlerFunc[controllers.AttendanceController](container, "GetTodayAttendance"))

		// 考勤统计
		attendance.GET("/statistics", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.AttendanceController](container, "GetAttendanceStatistics"))
	}

	// 请假管理
	leave := router.Group("/leave")
	leave.Use(middleware.JWTAuth())
	{
		// 请假申请
		leave.POST("", 
			middleware.RequireAnyRole("admin", "hr", "user"), 
			middleware.ValidateAndBindJSON(), 
			utils.CreateHandlerFunc[controllers.AttendanceController](container, "CreateLeave"))
		
		// 请假记录查询
		leave.GET("/records", 
			middleware.RequireAnyRole("admin", "hr", "user"), 
			utils.CreateHandlerFunc[controllers.AttendanceController](container, "GetLeaveRecords"))
		
		// 请假审批
		leave.PUT("/:id/approve", 
			middleware.ValidateNumericID(), 
			middleware.RequireAnyRole("admin", "hr"), 
			middleware.ValidateAndBindJSON(), 
			utils.CreateHandlerFunc[controllers.AttendanceController](container, "ApproveLeave"))
		
		// 待审批请假
		leave.GET("/pending", 
			middleware.RequireAnyRole("admin", "hr"), 
			utils.CreateHandlerFunc[controllers.AttendanceController](container, "GetPendingLeaves"))
	}
}