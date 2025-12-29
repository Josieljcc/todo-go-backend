package handlers

import (
	"net/http"
	"strconv"
	"todo-go-backend/internal/errors"
	"todo-go-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// CommentHandler manages comment handlers
type CommentHandler struct {
	commentService services.CommentService
}

// NewCommentHandler creates a new instance of CommentHandler
func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateCommentRequest represents a comment creation request
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000" example:"This is a comment on the task"`
	TaskID  uint   `json:"task_id" binding:"required" example:"1"`
}

// UpdateCommentRequest represents a comment update request
type UpdateCommentRequest struct {
	Content *string `json:"content" binding:"omitempty,min=1,max=5000" example:"Updated comment text"`
}

// CreateComment creates a new comment on a task
// @Summary      Create a comment on a task
// @Description  Creates a new comment on a task. User must own the task or have assigned it.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateCommentRequest  true  "Comment creation data"
// @Success      201      {object}  models.Comment
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      403      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	userID := c.GetUint("user_id")

	createReq := &services.CreateCommentRequest{
		Content: req.Content,
		TaskID:  req.TaskID,
	}

	comment, err := h.commentService.Create(userID, createReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComments retrieves all comments for a task
// @Summary      Get comments for a task
// @Description  Retrieves all comments for a specific task. User must own the task or have assigned it.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int  true  "Task ID"
// @Success      200      {array}   models.Comment
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      403      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /tasks/{id}/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.NewInvalidInputError("Invalid task ID"))
		return
	}

	comments, err := h.commentService.GetByTaskID(userID, uint(taskID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

// GetComment retrieves a specific comment
// @Summary      Get a comment by ID
// @Description  Retrieves a specific comment by its ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Comment ID"
// @Success      200  {object}  models.Comment
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.NewInvalidInputError("Invalid comment ID"))
		return
	}

	comment, err := h.commentService.GetByID(userID, uint(commentID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, comment)
}

// UpdateComment updates a comment
// @Summary      Update a comment
// @Description  Updates an existing comment. Only the comment author can update it.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "Comment ID"
// @Param        request  body      UpdateCommentRequest true  "Comment update data"
// @Success      200      {object}  models.Comment
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      403      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.NewInvalidInputError("Invalid comment ID"))
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	updateReq := &services.UpdateCommentRequest{
		Content: req.Content,
	}

	comment, err := h.commentService.Update(userID, uint(commentID), updateReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment deletes a comment
// @Summary      Delete a comment
// @Description  Deletes a comment by its ID. Only the comment author can delete it.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Comment ID"
// @Success      200  {object}  SuccessResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.NewInvalidInputError("Invalid comment ID"))
		return
	}

	if err := h.commentService.Delete(userID, uint(commentID)); err != nil {
		handleError(c, err)
		return
	}

	handleSuccess(c, http.StatusOK, "Comment deleted successfully", nil)
}

