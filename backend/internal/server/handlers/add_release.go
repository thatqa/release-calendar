package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddRelease struct {
	useCase usecase.ReleaseCreator
}

func NewAddRelease(
	db *gorm.DB,
) *AddRelease {
	return &AddRelease{
		useCase: usecase.NewReleaseCreator(db),
	}
}

func (h *AddRelease) Handle(c *gin.Context) {
	var request usecase.ReleaseCMD
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	release, err := h.useCase.AddRelease(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, release)
}
