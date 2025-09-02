package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeleteRelease struct {
	useCase usecase.ReleaseDeleter
}

func NewDeleteRelease(
	db *gorm.DB,
) *DeleteRelease {
	return &DeleteRelease{
		useCase: usecase.NewReleaseDeleter(db),
	}
}

func (h *DeleteRelease) Handle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.useCase.DeleteRelease(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "ok")
}
