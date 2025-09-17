package usecase

import (
	"release-calendar/backend/internal/repostiry"
	"time"

	"gorm.io/gorm"
)

type ReleaseDaysGetter struct {
	repo repostiry.ReleaseRepository
}

func NewReleaseDaysGetter(db *gorm.DB) ReleaseDaysGetter {
	return ReleaseDaysGetter{repo: repostiry.NewReleaseRepository(db)}
}

func (u *ReleaseDaysGetter) Get(from, to time.Time) (map[string][]string, error) {
	return u.repo.GetStatusesByRange(from, to)
}
