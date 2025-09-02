package usecase

import (
	"fmt"
	"release-calendar/backend/internal/models"
)

func (u *ReleaseGetter) ListRelease(date string, status string, duty string) ([]*models.Release, error) {
	releases, err := u.releaseRepository.GetByDateAndStatus(date, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get releases")
	}

	if duty != "" {
		filtered := make([]*models.Release, 0, len(releases))
		for _, r := range releases {
			for _, u := range r.DutyUsers {
				if u == duty {
					filtered = append(filtered, r)
					break
				}
			}
		}
		releases = filtered
	}

	return releases, nil
}
