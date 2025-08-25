package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateRelease struct {
	useCase usecase.ReleaseUpdater
}

func NewUpdateRelease(
	db *gorm.DB,
) *UpdateRelease {
	return &UpdateRelease{
		useCase: usecase.NewReleaseUpdater(db),
	}
}

func (h *UpdateRelease) Handle(c *gin.Context) {
	var request usecase.ReleaseCMD
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	release, err := h.useCase.UpdateRelease(uint(id), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, release)
}
