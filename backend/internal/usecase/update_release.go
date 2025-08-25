package usecase

import (
	"fmt"
	"release-calendar/backend/internal/models"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type ReleaseUpdater struct {
	releaseRepository repostiry.ReleaseRepository
	linkRepository    repostiry.LinkRepository
}

func NewReleaseUpdater(db *gorm.DB) ReleaseUpdater {
	return ReleaseUpdater{
		releaseRepository: repostiry.NewReleaseRepository(db),
		linkRepository:    repostiry.NewLinkRepository(db),
	}
}

func (u *ReleaseUpdater) UpdateRelease(ID uint, request ReleaseCMD) (*models.Release, error) {
	release, err := u.releaseRepository.GetById(ID)
	if err != nil {
		return nil, fmt.Errorf("release not found")
	}

	release.Title = request.Title
	release.Date = request.Date
	release.Status = request.Status
	release.Notes = request.Notes
	release.DutyUsers = request.DutyUsers

	if err := u.releaseRepository.Update(release); err != nil {
		return nil, fmt.Errorf("failed to update release")
	}

	newLinks := make([]*models.Link, 0)
	existing := map[uint]*models.Link{}
	for _, l := range release.Links {
		existing[l.ID] = l
	}
	seen := map[uint]bool{}
	for _, rl := range request.Links {
		if rl.ID == 0 {
			newLinks = append(newLinks, &models.Link{ReleaseID: release.ID, Name: rl.Name, URL: rl.URL})
		} else {
			if l, ok := existing[rl.ID]; ok {
				l.Name = rl.Name
				l.URL = rl.URL
				if err := u.linkRepository.Update(l); err != nil {
					return nil, fmt.Errorf("failed to update link")
				}
				seen[rl.ID] = true
			}
		}
	}
	if len(newLinks) > 0 {
		if err := u.linkRepository.Add(newLinks); err != nil {
			return nil, err
		}
	}
	deletedLinks := make([]*models.Link, 0)
	for id, l := range existing {
		if !seen[id] {
			deletedLinks = append(deletedLinks, l)
		}
	}
	if len(deletedLinks) > 0 {
		if err := u.linkRepository.Delete(deletedLinks); err != nil {
			return nil, fmt.Errorf("failed to delete links")
		}
	}

	return release, nil
}
