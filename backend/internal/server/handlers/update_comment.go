package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateComment struct {
	useCase usecase.CommentUpdater
}

func NewUpdateComment(
	db *gorm.DB,
) *UpdateComment {
	return &UpdateComment{
		useCase: usecase.NewCommentUpdater(db),
	}
}

func (h *UpdateComment) Handle(c *gin.Context) {
	commentID, _ := strconv.Atoi(c.Param("commentId"))
	var request usecase.CommentCMD
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment, err := h.useCase.UpdateComment(uint(commentID), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}
