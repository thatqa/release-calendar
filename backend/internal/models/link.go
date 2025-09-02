package models

import "time"

type Link struct {
	ID        uint      `json:"id"        gorm:"column:id"`
	ReleaseID uint      `json:"releaseId" gorm:"column:release_id"`
	Name      string    `json:"name"      gorm:"column:name"`
	URL       string    `json:"url"       gorm:"column:url"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Link) TableName() string { return "links" }
