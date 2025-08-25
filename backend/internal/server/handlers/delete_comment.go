package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeleteComment struct {
	useCase usecase.CommentDeleter
}

func NewDeleteComment(
	db *gorm.DB,
) *DeleteComment {
	return &DeleteComment{
		useCase: usecase.NewCommentDeleter(db),
	}
}

func (h *DeleteComment) Handle(c *gin.Context) {
	commentID, _ := strconv.Atoi(c.Param("commentId"))
	err := h.useCase.DeleteComment(uint(commentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "ok")
}
