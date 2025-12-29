// @title           Todo API
// @version         1.0
// @description     A RESTful API for managing tasks with JWT authentication and user task assignment
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"log"
	_ "todo-go-backend/docs" // Swagger documentation
	"todo-go-backend/internal/config"
	"todo-go-backend/internal/database"
	"todo-go-backend/internal/handlers"
	"todo-go-backend/internal/middleware"
	"todo-go-backend/internal/notifications"
	"todo-go-backend/internal/repositories"
	"todo-go-backend/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	taskRepo := repositories.NewTaskRepository()
	tagRepo := repositories.NewTagRepository()
	commentRepo := repositories.NewCommentRepository()

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	taskService := services.NewTaskService(taskRepo, userRepo, tagRepo)
	tagService := services.NewTagService(tagRepo)
	commentService := services.NewCommentService(commentRepo, taskRepo)

	// Initialize notification services
	emailService := notifications.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)
	telegramService := notifications.NewTelegramService(cfg.TelegramBotToken)
	notificationRepo := repositories.NewNotificationRepository()
	notificationService := notifications.NewNotificationService(
		emailService,
		telegramService,
		notificationRepo,
		taskRepo,
		userRepo,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)
	tagHandler := handlers.NewTagHandler(tagService)
	commentHandler := handlers.NewCommentHandler(commentService)
	userHandler := handlers.NewUserHandler(notificationService)

	// Start notification scheduler
	go notifications.StartScheduler(cfg, notificationService)

	// Setup router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware(cfg))

	// Health check endpoint
	// @Summary     Health check endpoint
	// @Description Returns the health status of the API
	// @Tags        health
	// @Accept      json
	// @Produce     json
	// @Success     200  {object}  map[string]string
	// @Router      /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	api := router.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Tasks routes
		protected.GET("/tasks", taskHandler.GetTasks)
		protected.POST("/tasks", taskHandler.CreateTask)

		// Comments routes for tasks (must be before /tasks/:id to avoid route conflict)
		// Using /tasks/:id/comments with same parameter name to avoid Gin route conflict
		protected.GET("/tasks/:id/comments", commentHandler.GetComments)

		// Tasks routes with ID (must be after /tasks/:id/comments)
		protected.GET("/tasks/:id", taskHandler.GetTask)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)

		// Tags routes
		protected.GET("/tags", tagHandler.GetTags)
		protected.GET("/tags/:id", tagHandler.GetTag)
		protected.POST("/tags", tagHandler.CreateTag)
		protected.PUT("/tags/:id", tagHandler.UpdateTag)
		protected.DELETE("/tags/:id", tagHandler.DeleteTag)

		// Comments routes
		protected.POST("/comments", commentHandler.CreateComment)
		protected.GET("/comments/:id", commentHandler.GetComment)
		protected.PUT("/comments/:id", commentHandler.UpdateComment)
		protected.DELETE("/comments/:id", commentHandler.DeleteComment)

		// User settings routes
		protected.PUT("/users/telegram-chat-id", userHandler.UpdateTelegramChatID)
		protected.PUT("/users/notifications-enabled", userHandler.UpdateNotificationsEnabled)

		// Notification test routes (for testing)
		protected.POST("/notifications/test", userHandler.TestNotifications)
		protected.GET("/notifications/debug", userHandler.GetNotificationDebugInfo)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
