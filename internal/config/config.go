package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	JWTSecret    string
	DatabasePath string
	// MySQL configuration
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	// CORS configuration
	CORSAllowedOrigins   string // Comma-separated list of allowed origins (e.g., "http://localhost:3000,https://example.com")
	CORSAllowedMethods   string // Comma-separated list of allowed methods (default: "GET,POST,PUT,DELETE,OPTIONS")
	CORSAllowedHeaders   string // Comma-separated list of allowed headers (default: "Content-Type,Authorization")
	CORSExposedHeaders   string // Comma-separated list of exposed headers
	CORSAllowCredentials bool   // Whether to allow credentials (default: true)
	CORSMaxAge           int    // Max age for preflight requests in seconds (default: 3600)
	// Notifications configuration
	NotificationsEnabled      bool   // Enable/disable notifications (default: true)
	NotificationCheckInterval string // Cron expression for notification check (default: "0 * * * *" - every hour)
	// Email SMTP configuration
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	// Telegram Bot configuration
	TelegramBotToken string // Telegram bot token
}

func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		// Log warning but continue (env vars might be set via Docker/OS)
		// This is expected in Docker environments
	}

	// Parse CORS max age
	corsMaxAge := 3600 // Default: 1 hour
	if maxAgeStr := getEnv("CORS_MAX_AGE", ""); maxAgeStr != "" {
		if parsed, err := parseInt(maxAgeStr); err == nil {
			corsMaxAge = parsed
		}
	}

	// Parse CORS allow credentials
	corsAllowCredentials := true // Default: true
	if allowCredsStr := getEnv("CORS_ALLOW_CREDENTIALS", ""); allowCredsStr != "" {
		corsAllowCredentials = allowCredsStr == "true" || allowCredsStr == "1"
	}

	// Parse notifications enabled
	notificationsEnabled := true // Default: enabled
	if enabledStr := getEnv("NOTIFICATIONS_ENABLED", ""); enabledStr != "" {
		notificationsEnabled = enabledStr == "true" || enabledStr == "1"
	}

	config := &Config{
		Port:                      getEnv("PORT", "8080"),
		JWTSecret:                 getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		DatabasePath:              getEnv("DATABASE_PATH", "todo.db"),
		DatabaseHost:              getEnv("DATABASE_HOST", ""),
		DatabasePort:              getEnv("DATABASE_PORT", "3306"),
		DatabaseUser:              getEnv("DATABASE_USER", ""),
		DatabasePassword:          getEnv("DATABASE_PASSWORD", ""),
		DatabaseName:              getEnv("DATABASE_NAME", ""),
		CORSAllowedOrigins:        getEnv("CORS_ALLOWED_ORIGINS", "*"), // Default: allow all origins (including same-origin)
		CORSAllowedMethods:        getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS,PATCH"),
		CORSAllowedHeaders:        getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,Accept,Origin"),
		CORSExposedHeaders:        getEnv("CORS_EXPOSED_HEADERS", ""),
		CORSAllowCredentials:      corsAllowCredentials,
		CORSMaxAge:                corsMaxAge,
		NotificationsEnabled:      notificationsEnabled,
		NotificationCheckInterval: getEnv("NOTIFICATION_CHECK_INTERVAL", "0 * * * *"), // Default: every hour
		SMTPHost:                  getEnv("SMTP_HOST", ""),
		SMTPPort:                  getEnv("SMTP_PORT", "587"),
		SMTPUser:                  getEnv("SMTP_USER", ""),
		SMTPPassword:              getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:                  getEnv("SMTP_FROM", ""),
		TelegramBotToken:          getEnv("TELEGRAM_BOT_TOKEN", ""),
	}

	// Log configuration status (without sensitive data)
	logConfigStatus(config)

	return config, nil
}

// UseMySQL returns true if MySQL configuration is provided
func (c *Config) UseMySQL() bool {
	return c.DatabaseHost != "" && c.DatabaseUser != "" && c.DatabaseName != ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// logConfigStatus logs configuration status without sensitive data
func logConfigStatus(cfg *Config) {
	log.Println("=== Configuration Status ===")
	log.Printf("Port: %s", cfg.Port)
	log.Printf("CORS Allowed Origins: %s", cfg.CORSAllowedOrigins)
	log.Printf("CORS Allow Credentials: %v", cfg.CORSAllowCredentials)
	log.Printf("CORS Allowed Methods: %s", cfg.CORSAllowedMethods)
	log.Printf("CORS Allowed Headers: %s", cfg.CORSAllowedHeaders)
	log.Printf("Notifications Enabled: %v", cfg.NotificationsEnabled)
	log.Printf("Notification Interval: %s", cfg.NotificationCheckInterval)
	log.Printf("SMTP Host: %s", maskIfEmpty(cfg.SMTPHost))
	log.Printf("SMTP Port: %s", cfg.SMTPPort)
	log.Printf("SMTP User: %s", maskIfEmpty(cfg.SMTPUser))
	log.Printf("SMTP Password: %s", maskIfEmpty(cfg.SMTPPassword))
	log.Printf("SMTP From: %s", maskIfEmpty(cfg.SMTPFrom))
	log.Printf("Telegram Bot Token: %s", maskIfEmpty(cfg.TelegramBotToken))
	log.Println("===========================")
}

func maskIfEmpty(s string) string {
	if s == "" {
		return "[NOT CONFIGURED]"
	}
	return "[CONFIGURED]"
}
