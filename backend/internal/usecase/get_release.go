package usecase

import (
	"release-calendar/backend/internal/models"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type ReleaseGetter struct {
	releaseRepository repostiry.ReleaseRepository
}

func NewReleaseGetter(db *gorm.DB) ReleaseGetter {
	return ReleaseGetter{
		releaseRepository: repostiry.NewReleaseRepository(db),
	}
}

func (u *ReleaseGetter) GetRelease(ID uint) (*models.Release, error) {
	release, err := u.releaseRepository.GetById(ID)
	if err != nil {
		return nil, err
	}

	return release, nil
}
