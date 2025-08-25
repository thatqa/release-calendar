package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListRelease struct {
	useCase usecase.ReleaseGetter
}

func NewListRelease(
	db *gorm.DB,
) *ListRelease {
	return &ListRelease{
		useCase: usecase.NewReleaseGetter(db),
	}
}

type ListReleaseParams struct {
	Date   string `form:"date"`
	Status string `form:"status"`
	Duty   string `form:"duty"`
}

func (h *ListRelease) Handle(c *gin.Context) {
	var params ListReleaseParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	release, err := h.useCase.ListRelease(params.Date, params.Status, params.Duty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, release)
}
