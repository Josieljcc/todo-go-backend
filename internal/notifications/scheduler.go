package notifications

import (
	"log"
	"todo-go-backend/internal/config"

	"github.com/robfig/cron/v3"
)

// StartScheduler starts the notification scheduler
func StartScheduler(cfg *config.Config, notificationService *NotificationService) {
	if !cfg.NotificationsEnabled {
		log.Println("Notifications are disabled")
		return
	}

	c := cron.New()

	// Add notification check job
	_, err := c.AddFunc(cfg.NotificationCheckInterval, func() {
		log.Println("Running notification check...")
		if err := notificationService.CheckAndSendNotifications(); err != nil {
			log.Printf("Error checking notifications: %v", err)
		} else {
			log.Println("Notification check completed")
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule notifications: %v", err)
	}

	log.Printf("Notification scheduler started with interval: %s", cfg.NotificationCheckInterval)
	c.Start()
}

