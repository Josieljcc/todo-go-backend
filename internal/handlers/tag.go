package handlers

import (
	"net/http"
	"strconv"
	"todo-go-backend/internal/errors"
	"todo-go-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// TagHandler manages tag handlers
type TagHandler struct {
	tagService services.TagService
}

// NewTagHandler creates a new instance of TagHandler
func NewTagHandler(tagService services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// CreateTagRequest represents a tag creation request
type CreateTagRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=50" example:"Important"`
	Color string `json:"color" example:"#FF5733"` // Optional: hex color code
}

// UpdateTagRequest represents a tag update request
type UpdateTagRequest struct {
	Name  *string `json:"name" example:"Updated Tag"`
	Color *string `json:"color" example:"#33FF57"`
}

// CreateTag creates a new tag
// @Summary      Create a new tag
// @Description  Creates a new custom tag for the authenticated user
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateTagRequest  true  "Tag creation data"
// @Success      201      {object}  models.Tag
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	userID := c.GetUint("user_id")

	createReq := &services.CreateTagRequest{
		Name:  req.Name,
		Color: req.Color,
	}

	tag, err := h.tagService.Create(userID, createReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetTags lists user tags
// @Summary      List user tags
// @Description  Retrieves all tags for the authenticated user
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Tag
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	userID := c.GetUint("user_id")

	tags, err := h.tagService.GetByUserID(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tags)
}

// GetTag retrieves a specific tag
// @Summary      Get a tag by ID
// @Description  Retrieves a specific tag by its ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Tag ID"
// @Success      200  {object}  models.Tag
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	userID := c.GetUint("user_id")
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.NewInvalidInputError("Invalid tag ID"))
		return
	}

	tag, err := h.tagService.GetByID(userID, uint(tagID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tag)
}

// UpdateTag updates a tag
// @Summary      Update a tag
// @Description  Updates an existing tag
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int              true  "Tag ID"
// @Param        request  body      UpdateTagRequest true  "Tag update data"
// @Success      200      {object}  models.Tag
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	userID := c.GetUint("user_id")
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.NewInvalidInputError("Invalid tag ID"))
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	updateReq := &services.UpdateTagRequest{
		Name:  req.Name,
		Color: req.Color,
	}

	tag, err := h.tagService.Update(userID, uint(tagID), updateReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteTag deletes a tag
// @Summary      Delete a tag
// @Description  Deletes a tag by its ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Tag ID"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	userID := c.GetUint("user_id")
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.NewInvalidInputError("Invalid tag ID"))
		return
	}

	if err := h.tagService.Delete(userID, uint(tagID)); err != nil {
		handleError(c, err)
		return
	}

	handleSuccess(c, http.StatusOK, "Tag deleted successfully", nil)
}

