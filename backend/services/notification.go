package services

import (
	"fmt"
	"log"
	"time"
)

type NotificationService struct {
	wsService *WebSocketService
}

func NewNotificationService(wsService *WebSocketService) *NotificationService {
	return &NotificationService{
		wsService: wsService,
	}
}

// 通知类型
type NotificationType string

const (
	NotificationInfo    NotificationType = "info"
	NotificationWarning NotificationType = "warning"
	NotificationError   NotificationType = "error"
	NotificationSuccess NotificationType = "success"
)

// 通知数据结构
type Notification struct {
	ID       string           `json:"id"`
	Type     NotificationType `json:"type"`
	Title    string           `json:"title"`
	Message  string           `json:"message"`
	UserID   string           `json:"user_id,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
	CreateAt time.Time        `json:"create_at"`
}

// 发送通知给特定用户
func (ns *NotificationService) SendToUser(userID string, notificationType NotificationType, title, message string, data map[string]interface{}) error {
	notification := Notification{
		ID:       fmt.Sprintf("notif_%d", time.Now().UnixNano()),
		Type:     notificationType,
		Title:    title,
		Message:  message,
		UserID:   userID,
		Data:     data,
		CreateAt: time.Now(),
	}

	err := ns.wsService.SendToUser(userID, MessageTypeNotification, notification)
	if err != nil {
		log.Printf("Failed to send notification to user %s: %v", userID, err)
		return err
	}

	log.Printf("Notification sent to user %s: %s", userID, title)
	return nil
}

// 广播通知给所有用户
func (ns *NotificationService) Broadcast(notificationType NotificationType, title, message string, data map[string]interface{}) {
	notification := Notification{
		ID:       fmt.Sprintf("broadcast_%d", time.Now().UnixNano()),
		Type:     notificationType,
		Title:    title,
		Message:  message,
		Data:     data,
		CreateAt: time.Now(),
	}

	ns.wsService.Broadcast(MessageTypeNotification, notification)
	log.Printf("Notification broadcasted: %s", title)
}

// HR系统特定的通知方法

// 员工入职通知
func (ns *NotificationService) NotifyEmployeeOnboard(employeeName, department string, managerUserID string) {
	ns.SendToUser(managerUserID, NotificationInfo, 
		"新员工入职", 
		fmt.Sprintf("员工 %s 已加入 %s 部门", employeeName, department),
		map[string]interface{}{
			"employee_name": employeeName,
			"department":    department,
			"event_type":    "employee_onboard",
		})
}

// 薪资发放通知
func (ns *NotificationService) NotifySalaryPayment(userID string, amount float64, month string) {
	ns.SendToUser(userID, NotificationSuccess,
		"薪资发放通知",
		fmt.Sprintf("您的 %s 月薪资 ¥%.2f 已发放", month, amount),
		map[string]interface{}{
			"amount":     amount,
			"month":      month,
			"event_type": "salary_payment",
		})
}

// 考勤异常通知
func (ns *NotificationService) NotifyAttendanceAnomaly(userID string, anomalyType, date string) {
	ns.SendToUser(userID, NotificationWarning,
		"考勤异常提醒",
		fmt.Sprintf("您在 %s 的考勤记录存在异常: %s", date, anomalyType),
		map[string]interface{}{
			"anomaly_type": anomalyType,
			"date":         date,
			"event_type":   "attendance_anomaly",
		})
}

// 系统维护通知
func (ns *NotificationService) NotifySystemMaintenance(startTime, endTime, reason string) {
	ns.Broadcast(NotificationWarning,
		"系统维护通知",
		fmt.Sprintf("系统将于 %s 至 %s 进行维护，原因: %s", startTime, endTime, reason),
		map[string]interface{}{
			"start_time":   startTime,
			"end_time":     endTime,
			"reason":       reason,
			"event_type":   "system_maintenance",
		})
}

// 任务提醒通知
func (ns *NotificationService) NotifyTaskReminder(userID string, taskTitle, dueDate string) {
	ns.SendToUser(userID, NotificationInfo,
		"任务提醒",
		fmt.Sprintf("任务 '%s' 将于 %s 到期，请及时处理", taskTitle, dueDate),
		map[string]interface{}{
			"task_title": taskTitle,
			"due_date":   dueDate,
			"event_type": "task_reminder",
		})
}