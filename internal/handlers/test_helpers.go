package handlers

import (
	"todo-go-backend/internal/database"
	"todo-go-backend/internal/middleware"
	"todo-go-backend/internal/models"
	"todo-go-backend/internal/repositories"
	"todo-go-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB cria um banco de dados em memória para testes
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Tag{}, &models.Comment{}, &models.Notification{})
	database.DB = db
	return db
}

// setupTestRouter cria um router de teste com handlers configurados
func setupTestRouter(jwtSecret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	taskRepo := repositories.NewTaskRepository()

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtSecret)
	tagRepo := repositories.NewTagRepository()
	taskService := services.NewTaskService(taskRepo, userRepo, tagRepo)

	// Initialize handlers
	authHandler := NewAuthHandler(authService)
	taskHandler := NewTaskHandler(taskService)

	// Public routes
	api := router.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.GET("/tasks", taskHandler.GetTasks)
		protected.GET("/tasks/:id", taskHandler.GetTask)
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)
	}

	return router
}
