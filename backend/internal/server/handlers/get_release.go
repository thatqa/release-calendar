package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetRelease struct {
	useCase usecase.ReleaseGetter
}

func NewGetRelease(
	db *gorm.DB,
) *GetRelease {
	return &GetRelease{
		useCase: usecase.NewReleaseGetter(db),
	}
}

type GetReleaseParams struct {
	ID int `form:"id" binding:"required"`
}

func (h *GetRelease) Handle(c *gin.Context) {
	var params GetReleaseParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter id is required"})
		return
	}
	release, err := h.useCase.GetRelease(uint(params.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, release)
}
