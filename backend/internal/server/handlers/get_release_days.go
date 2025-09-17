package handlers

import (
	"net/http"
	"time"

	"release-calendar/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetReleaseDays struct{ uc usecase.ReleaseDaysGetter }

func NewGetReleaseDays(db *gorm.DB) *GetReleaseDays {
	return &GetReleaseDays{uc: usecase.NewReleaseDaysGetter(db)}
}

func (h *GetReleaseDays) Handle(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to are required (YYYY-MM-DD)"})
		return
	}
	from, err1 := time.Parse("2006-01-02", fromStr)
	to, err2 := time.Parse("2006-01-02", toStr)
	if err1 != nil || err2 != nil || !to.After(from) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
		return
	}

	m, err := h.uc.Get(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}
