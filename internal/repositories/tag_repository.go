package repositories

import (
	"todo-go-backend/internal/database"
	"todo-go-backend/internal/models"
)

// TagRepository defines the interface for tag operations
type TagRepository interface {
	Create(tag *models.Tag) error
	FindByID(id uint) (*models.Tag, error)
	FindByUserID(userID uint) ([]models.Tag, error)
	FindByIDAndUserID(id, userID uint) (*models.Tag, error)
	FindByNameAndUserID(name string, userID uint) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(id uint) error
	FindByIDs(ids []uint, userID uint) ([]models.Tag, error)
	ExistsByNameAndUserID(name string, userID uint) (bool, error)
}

type tagRepository struct{}

// NewTagRepository creates a new instance of TagRepository
func NewTagRepository() TagRepository {
	return &tagRepository{}
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return database.DB.Create(tag).Error
}

func (r *tagRepository) FindByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) FindByUserID(userID uint) ([]models.Tag, error) {
	var tags []models.Tag
	if err := database.DB.Where("user_id = ?", userID).Order("name ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) FindByIDAndUserID(id, userID uint) (*models.Tag, error) {
	var tag models.Tag
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) FindByNameAndUserID(name string, userID uint) (*models.Tag, error) {
	var tag models.Tag
	if err := database.DB.Where("name = ? AND user_id = ?", name, userID).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Update(tag *models.Tag) error {
	return database.DB.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Tag{}, id).Error
}

func (r *tagRepository) FindByIDs(ids []uint, userID uint) ([]models.Tag, error) {
	var tags []models.Tag
	if err := database.DB.Where("id IN ? AND user_id = ?", ids, userID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) ExistsByNameAndUserID(name string, userID uint) (bool, error) {
	var count int64
	if err := database.DB.Model(&models.Tag{}).
		Where("name = ? AND user_id = ?", name, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
