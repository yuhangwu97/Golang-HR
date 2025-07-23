package controllers

import (
	"gin-project/services"
	"gin-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebSocketController struct {
	wsService *services.WebSocketService
}

func NewWebSocketController(wsService *services.WebSocketService) *WebSocketController {
	return &WebSocketController{
		wsService: wsService,
	}
}

// WebSocket连接端点
func (wsc *WebSocketController) HandleWebSocket(c *gin.Context) {
	wsc.wsService.HandleConnection(c)
}

// 获取WebSocket状态信息
func (wsc *WebSocketController) GetStatus(c *gin.Context) {
	clients := wsc.wsService.GetClients()
	onlineCount := wsc.wsService.GetOnlineUsers()

	utils.SuccessResponse(c, http.StatusOK, "WebSocket status retrieved successfully", gin.H{
		"online_users": onlineCount,
		"clients":      clients,
		"status":       "running",
	})
}

// 广播消息API端点
func (wsc *WebSocketController) BroadcastMessage(c *gin.Context) {
	var request struct {
		Type    string      `json:"type" binding:"required"`
		Message interface{} `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	messageType := services.MessageType(request.Type)
	wsc.wsService.Broadcast(messageType, request.Message)

	utils.SuccessResponse(c, http.StatusOK, "Message broadcasted successfully", gin.H{
		"message": "Message broadcasted successfully",
	})
}

// 发送消息给特定用户
func (wsc *WebSocketController) SendToUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	var request struct {
		Type    string      `json:"type" binding:"required"`
		Message interface{} `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	messageType := services.MessageType(request.Type)
	err := wsc.wsService.SendToUser(userID, messageType, request.Message)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not connected")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Message sent successfully", gin.H{
		"user_id": userID,
	})
}

// 通知系统事件
func (wsc *WebSocketController) NotifyEvent(c *gin.Context) {
	var request struct {
		Event string                 `json:"event" binding:"required"`
		Data  map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	wsc.wsService.NotifySystemEvent(request.Event, request.Data)

	utils.SuccessResponse(c, http.StatusOK, "Event notification sent", gin.H{
		"event": request.Event,
	})
}
