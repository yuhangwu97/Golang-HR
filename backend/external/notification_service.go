package external

import (
	"log"
)

type NotificationService struct {
	apiKey  string
	baseURL string
	client  *HTTPClient
}

func NewNotificationService(apiKey, baseURL string) *NotificationService {
	return &NotificationService{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  NewHTTPClient(baseURL, 30000),
	}
}

func (n *NotificationService) SendNotification(userID, message string) error {
	headers := map[string]string{
		"Authorization": "Bearer " + n.apiKey,
	}

	data := map[string]interface{}{
		"user_id": userID,
		"message": message,
		"type":    "push",
	}

	_, err := n.client.Post("/notify", data, headers)
	if err != nil {
		log.Printf("Failed to send notification to user %s: %v", userID, err)
		return err
	}

	log.Printf("Notification sent successfully to user %s", userID)
	return nil
}

func (n *NotificationService) SendBulkNotification(userIDs []string, message string) error {
	headers := map[string]string{
		"Authorization": "Bearer " + n.apiKey,
	}

	data := map[string]interface{}{
		"user_ids": userIDs,
		"message":  message,
		"type":     "bulk_push",
	}

	_, err := n.client.Post("/bulk-notify", data, headers)
	if err != nil {
		log.Printf("Failed to send bulk notification: %v", err)
		return err
	}

	log.Printf("Bulk notification sent successfully to %d users", len(userIDs))
	return nil
}