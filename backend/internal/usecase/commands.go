package usecase

import (
	"release-calendar/backend/internal/models"
	"time"
)

type ReleaseCMD struct {
	Title     string               `json:"title" binding:"required"`
	Date      time.Time            `json:"date" binding:"required"`
	Status    models.ReleaseStatus `json:"status" binding:"required"`
	Notes     string               `json:"notes"`
	DutyUsers []string             `json:"dutyUsers"`
	Links     []LinkCMD            `json:"links"`
}

type LinkCMD struct {
	ID   uint   `json:"id"`
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

type CommentCMD struct {
	Author  string `json:"author" binding:"required"`
	Message string `json:"message" binding:"required"`
}
