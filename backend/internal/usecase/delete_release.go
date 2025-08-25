package usecase

import (
	"fmt"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type ReleaseDeleter struct {
	releaseRepository repostiry.ReleaseRepository
}

func NewReleaseDeleter(db *gorm.DB) ReleaseDeleter {
	return ReleaseDeleter{
		releaseRepository: repostiry.NewReleaseRepository(db),
	}
}

func (u *ReleaseDeleter) DeleteRelease(ID uint) error {
	if err := u.releaseRepository.DeleteById(ID); err != nil {
		return fmt.Errorf("failed to delete release")
	}
	return nil
}
