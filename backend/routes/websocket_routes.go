package routes

import (
	"gin-project/controllers"
	"gin-project/services"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func SetupWebSocketRoutes(router *gin.RouterGroup, container *utils.Container) {
	// 创建WebSocket服务实例
	wsService := services.NewWebSocketService()
	wsService.Run() // 启动WebSocket服务

	// 创建WebSocket控制器
	wsController := controllers.NewWebSocketController(wsService)

	// WebSocket连接路由 (不需要认证，但可以通过query参数传递用户信息)
	router.GET("", wsController.HandleWebSocket)

	// WebSocket管理API路由 (需要认证)
	wsGroup := router.Group("/api")
	// wsGroup.Use(middleware.AuthMiddleware()) // 需要认证
	{
		// 获取WebSocket状态
		wsGroup.GET("/status", wsController.GetStatus)

		// 广播消息给所有连接的客户端
		wsGroup.POST("/broadcast", wsController.BroadcastMessage)

		// 发送消息给特定用户
		wsGroup.POST("/send/:user_id", wsController.SendToUser)

		// 通知系统事件
		wsGroup.POST("/notify", wsController.NotifyEvent)
	}
}
