package models

import (
	"time"
)

type ReleaseStatus string

const (
	StatusPlanned ReleaseStatus = "planned"
	StatusSuccess ReleaseStatus = "success"
	StatusFailed  ReleaseStatus = "failed"
)

type Release struct {
	ID        uint          `json:"id"          gorm:"column:id"`
	Title     string        `json:"title"       gorm:"column:title"`
	Date      time.Time     `json:"date"        gorm:"column:date"`
	Status    ReleaseStatus `json:"status"      gorm:"column:status"`
	Notes     string        `json:"notes"       gorm:"column:notes"`
	DutyUsers []string      `json:"dutyUsers"   gorm:"column:duty_users;serializer:json"`
	CreatedAt time.Time     `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time     `json:"updatedAt" gorm:"column:updated_at"`
	Links     []*Link       `json:"links"    gorm:"foreignKey:ReleaseID;references:ID"`
	Comments  []Comment     `json:"comments" gorm:"foreignKey:ReleaseID;references:ID"`
}

func (Release) TableName() string { return "releases" }
