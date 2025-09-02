package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddComment struct {
	useCase usecase.CommentCreator
}

func NewAddComment(
	db *gorm.DB,
) *AddComment {
	return &AddComment{
		useCase: usecase.NewCommentCreator(db),
	}
}

func (h *AddComment) Handle(c *gin.Context) {
	releaseID, _ := strconv.Atoi(c.Param("id"))
	var request usecase.CommentCMD
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment, err := h.useCase.AddComment(uint(releaseID), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}
