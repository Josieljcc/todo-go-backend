package notifications

import (
	"fmt"
	"net/smtp"
	"todo-go-backend/internal/models"
)

// EmailService handles email notifications
type EmailService struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

// NewEmailService creates a new email service
func NewEmailService(host, port, user, password, from string) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
	}
}

// SendNotification sends a notification email
func (s *EmailService) SendNotification(user *models.User, task *models.Task, notificationType models.NotificationType) error {
	if s.host == "" || s.user == "" || s.password == "" {
		return fmt.Errorf("email service not configured")
	}

	subject, body := s.buildEmailContent(task, notificationType)

	// Setup authentication
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	// Email message
	msg := []byte(fmt.Sprintf("To: %s\r\n", user.Email) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	err := smtp.SendMail(addr, auth, s.from, []string{user.Email}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// buildEmailContent builds email subject and body based on notification type
func (s *EmailService) buildEmailContent(task *models.Task, notificationType models.NotificationType) (string, string) {
	var subject string
	var body string

	switch notificationType {
	case models.NotificationTypeDueSoon:
		subject = fmt.Sprintf("⏰ Tarefa vence amanhã: %s", task.Title)
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Tarefa vence amanhã!</h2>
				<p><strong>%s</strong></p>
				<p>%s</p>
				<p><strong>Prioridade:</strong> %s</p>
				<p><strong>Data de vencimento:</strong> %s</p>
			</body>
			</html>
		`, task.Title, task.Description, task.Priority, task.DueDate.Format("02/01/2006"))
	case models.NotificationTypeDueToday:
		subject = fmt.Sprintf("📅 Tarefa vence hoje: %s", task.Title)
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Tarefa vence hoje!</h2>
				<p><strong>%s</strong></p>
				<p>%s</p>
				<p><strong>Prioridade:</strong> %s</p>
				<p><strong>Data de vencimento:</strong> %s</p>
			</body>
			</html>
		`, task.Title, task.Description, task.Priority, task.DueDate.Format("02/01/2006"))
	case models.NotificationTypeOverdue:
		subject = fmt.Sprintf("⚠️ Tarefa atrasada: %s", task.Title)
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Tarefa atrasada!</h2>
				<p><strong>%s</strong></p>
				<p>%s</p>
				<p><strong>Prioridade:</strong> %s</p>
				<p><strong>Data de vencimento:</strong> %s</p>
			</body>
			</html>
		`, task.Title, task.Description, task.Priority, task.DueDate.Format("02/01/2006"))
	}

	return subject, body
}

