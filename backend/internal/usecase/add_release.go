package usecase

import (
	"fmt"
	"release-calendar/backend/internal/models"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type ReleaseCreator struct {
	db *gorm.DB
}

func NewReleaseCreator(db *gorm.DB) ReleaseCreator {
	return ReleaseCreator{
		db: db,
	}
}

func (u *ReleaseCreator) AddRelease(request ReleaseCMD) (*models.Release, error) {
	var createdRelease *models.Release

	err := u.db.Transaction(func(tx *gorm.DB) error {
		releaseRepository := repostiry.NewReleaseRepository(tx)
		linkRepository := repostiry.NewLinkRepository(tx)

		release := &models.Release{
			Title:     request.Title,
			Date:      request.Date,
			Status:    request.Status,
			Notes:     request.Notes,
			DutyUsers: request.DutyUsers,
		}
		if err := releaseRepository.Add(release); err != nil {
			return fmt.Errorf("failed to add release")
		}

		linksCount := len(request.Links)
		if linksCount > 0 {
			links := make([]*models.Link, linksCount)
			for i, link := range request.Links {
				links[i] = &models.Link{ReleaseID: release.ID, Name: link.Name, URL: link.URL}
			}
			release.Links = links
			if err := linkRepository.Add(release.Links); err != nil {
				return fmt.Errorf("failed to add links")
			}
		}
		createdRelease = release
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}
