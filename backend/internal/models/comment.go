package models

import "time"

type Comment struct {
	ID        uint      `json:"id"        gorm:"column:id"`
	ReleaseID uint      `json:"releaseId" gorm:"column:release_id"`
	Author    string    `json:"author"    gorm:"column:author"`
	Message   string    `json:"message"   gorm:"column:message"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Comment) TableName() string { return "comments" }
