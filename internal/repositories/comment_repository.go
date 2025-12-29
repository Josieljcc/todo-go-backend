package repositories

import (
	"todo-go-backend/internal/database"
	"todo-go-backend/internal/models"
)

// CommentRepository defines the interface for comment operations
type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByID(id uint) (*models.Comment, error)
	FindByTaskID(taskID uint) ([]models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type commentRepository struct{}

// NewCommentRepository creates a new instance of CommentRepository
func NewCommentRepository() CommentRepository {
	return &commentRepository{}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return database.DB.Create(comment).Error
}

func (r *commentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := database.DB.
		Preload("User").
		Preload("Task").
		First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByTaskID(taskID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := database.DB.
		Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return database.DB.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Comment{}, id).Error
}

func (r *commentRepository) Exists(id uint) (bool, error) {
	var count int64
	if err := database.DB.Model(&models.Comment{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

