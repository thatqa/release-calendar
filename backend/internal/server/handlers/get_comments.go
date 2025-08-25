package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetComments struct {
	useCase usecase.CommentGetter
}

func NewGetComments(
	db *gorm.DB,
) *GetComments {
	return &GetComments{
		useCase: usecase.NewCommentGetter(db),
	}
}

func (h *GetComments) Handle(c *gin.Context) {
	releaseID, _ := strconv.Atoi(c.Param("id"))
	comments, err := h.useCase.GetComments(uint(releaseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}
