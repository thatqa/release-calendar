package handlers

import (
	"net/http"
	"release-calendar/backend/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetSummary struct {
	uc usecase.ReleaseSummarizer
}

func NewGetSummary(db *gorm.DB) *GetSummary {
	return &GetSummary{uc: usecase.NewReleaseSummarizer(db)}
}

func (h *GetSummary) Handle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	resp, err := h.uc.Summarize(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"summary":  resp.Summary,
		"provider": resp.Provider,
	})
}
