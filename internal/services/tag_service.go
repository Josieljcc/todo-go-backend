package services

import (
	"todo-go-backend/internal/errors"
	"todo-go-backend/internal/models"
	"todo-go-backend/internal/repositories"
)

// TagService defines the interface for tag operations
type TagService interface {
	Create(userID uint, req *CreateTagRequest) (*models.Tag, error)
	GetByID(userID, tagID uint) (*models.Tag, error)
	GetByUserID(userID uint) ([]models.Tag, error)
	Update(userID, tagID uint, req *UpdateTagRequest) (*models.Tag, error)
	Delete(userID, tagID uint) error
}

// CreateTagRequest represents a tag creation request
type CreateTagRequest struct {
	Name  string
	Color string // Hex color code (e.g., #FF5733)
}

// UpdateTagRequest represents a tag update request
type UpdateTagRequest struct {
	Name  *string
	Color *string
}

type tagService struct {
	tagRepo repositories.TagRepository
}

// NewTagService creates a new instance of TagService
func NewTagService(tagRepo repositories.TagRepository) TagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

func (s *tagService) Create(userID uint, req *CreateTagRequest) (*models.Tag, error) {
	// Check if tag with same name already exists for this user
	exists, err := s.tagRepo.ExistsByNameAndUserID(req.Name, userID)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	if exists {
		return nil, errors.NewInvalidInputError("A tag with this name already exists")
	}

	// Validate color format if provided
	if req.Color != "" && !isValidHexColor(req.Color) {
		return nil, errors.NewInvalidInputError("Invalid color format. Use hex color code (e.g., #FF5733)")
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#808080" // Default gray
	}

	tag := &models.Tag{
		Name:   req.Name,
		Color:  color,
		UserID: userID,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	return tag, nil
}

func (s *tagService) GetByID(userID, tagID uint) (*models.Tag, error) {
	tag, err := s.tagRepo.FindByIDAndUserID(tagID, userID)
	if err != nil {
		return nil, errors.NewTaskNotFoundError() // Reuse error type
	}
	return tag, nil
}

func (s *tagService) GetByUserID(userID uint) ([]models.Tag, error) {
	tags, err := s.tagRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return tags, nil
}

func (s *tagService) Update(userID, tagID uint, req *UpdateTagRequest) (*models.Tag, error) {
	tag, err := s.tagRepo.FindByIDAndUserID(tagID, userID)
	if err != nil {
		return nil, errors.NewTaskNotFoundError()
	}

	if req.Name != nil {
		// Check if another tag with the same name already exists for this user
		existingTag, err := s.tagRepo.FindByNameAndUserID(*req.Name, userID)
		if err == nil && existingTag.ID != tagID {
			return nil, errors.NewInvalidInputError("A tag with this name already exists")
		}
		tag.Name = *req.Name
	}
	if req.Color != nil {
		if !isValidHexColor(*req.Color) {
			return nil, errors.NewInvalidInputError("Invalid color format. Use hex color code (e.g., #FF5733)")
		}
		tag.Color = *req.Color
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	return tag, nil
}

func (s *tagService) Delete(userID, tagID uint) error {
	tag, err := s.tagRepo.FindByIDAndUserID(tagID, userID)
	if err != nil {
		return errors.NewTaskNotFoundError()
	}

	if err := s.tagRepo.Delete(tag.ID); err != nil {
		return errors.NewInternalServerError(err)
	}

	return nil
}

// isValidHexColor validates hex color format
func isValidHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
