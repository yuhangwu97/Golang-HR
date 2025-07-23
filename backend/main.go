package main

import (
	"log"
	"net/http"

	"gin-project/config"
	"gin-project/middleware"
	"gin-project/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	config.ConnectMySQL()
	config.ConnectRedis()
	//config.AutoMigrate()
	defer config.CloseConnections()

	// 初始化依赖注入容器
	config.InitContainer()

	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "HR Management System API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// API路由组
	api := r.Group("/api/v1")

	// 使用依赖注入容器设置所有路由
	routes.SetupAuthRoutes(api, config.Container)
	routes.SetupUserRoutes(api, config.Container)
	routes.SetupEmployeeRoutes(api, config.Container)
	routes.SetupDepartmentRoutes(api, config.Container)
	routes.SetupPositionRoutes(api, config.Container)
	routes.SetupJobLevelRoutes(api, config.Container)
	routes.SetupSalaryRoutes(api, config.Container)
	routes.SetupAttendanceRoutes(api, config.Container)
	routes.SetupOrganizationRoutes(api, config.Container)

	ws := r.Group("/ws")
	routes.SetupWebSocketRoutes(ws, config.Container)

	log.Println("🚀 HR Management System API is running on :8090")
	r.Run(":8090")
}
