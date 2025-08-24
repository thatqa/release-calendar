package models

import (
	"time"

	"gorm.io/datatypes"
)

type ReleaseStatus string

const (
	StatusPlanned ReleaseStatus = "planned"
	StatusSuccess ReleaseStatus = "success"
	StatusFailed  ReleaseStatus = "failed"
)

type Release struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title"`
	Date      time.Time      `json:"date"`
	Status    ReleaseStatus  `json:"status"`
	Notes     string         `json:"notes"`
	DutyUsers datatypes.JSON `json:"dutyUsers" gorm:"type:JSON;default:'[]'"`
	Links     []Link         `json:"links"`
	Comments  []Comment      `json:"comments"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type Link struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ReleaseID uint      `json:"releaseId" gorm:"index;not null"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ReleaseID uint      `json:"releaseId" gorm:"index;not null"`
	Author    string    `json:"author"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
