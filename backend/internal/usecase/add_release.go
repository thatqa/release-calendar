package usecase

import (
	"fmt"
	"release-calendar/backend/internal/models"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type ReleaseCreator struct {
	releaseRepository repostiry.ReleaseRepository
	linkRepository    repostiry.LinkRepository
}

func NewReleaseCreator(db *gorm.DB) ReleaseCreator {
	return ReleaseCreator{
		releaseRepository: repostiry.NewReleaseRepository(db),
		linkRepository:    repostiry.NewLinkRepository(db),
	}
}

func (u *ReleaseCreator) AddRelease(request ReleaseCMD) (*models.Release, error) {
	release := &models.Release{
		Title:     request.Title,
		Date:      request.Date,
		Status:    request.Status,
		Notes:     request.Notes,
		DutyUsers: request.DutyUsers,
	}
	if err := u.releaseRepository.Add(release); err != nil {
		return nil, fmt.Errorf("failed to add release")
	}

	linksCount := len(request.Links)
	if linksCount > 0 {
		links := make([]*models.Link, linksCount)
		for i, link := range request.Links {
			links[i] = &models.Link{ReleaseID: release.ID, Name: link.Name, URL: link.URL}
		}
		release.Links = links
		if err := u.linkRepository.Add(release.Links); err != nil {
			return nil, fmt.Errorf("failed to add links")
		}
	}

	return release, nil
}
