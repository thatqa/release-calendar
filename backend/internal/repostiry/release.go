package repostiry

import (
	"release-calendar/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReleaseRepository interface {
	Add(release *models.Release) error
	GetById(ID uint) (*models.Release, error)
	GetByDateAndStatus(date string, status string) ([]*models.Release, error)
	Update(release *models.Release) error
	DeleteById(ID uint) error
	GetStatusesByRange(from, to time.Time) (map[string][]string, error)
}

type releaseRepository struct {
	db *gorm.DB
}

func NewReleaseRepository(
	db *gorm.DB,
) ReleaseRepository {
	return &releaseRepository{
		db: db,
	}
}

func (r *releaseRepository) Add(release *models.Release) error {
	return r.db.Create(release).Error
}

func (r *releaseRepository) GetById(ID uint) (*models.Release, error) {
	var release models.Release
	err := r.db.
		Preload("Links").
		Preload("Comments").
		Where("id = ?", ID).
		First(&release).Error
	if err != nil {
		return nil, err
	}

	return &release, nil
}

func (r *releaseRepository) Update(release *models.Release) error {
	return r.db.Save(release).Error
}

func (r *releaseRepository) GetByDateAndStatus(date string, status string) ([]*models.Release, error) {
	var releases []*models.Release
	q := r.db.
		Preload("Links").
		Order("date ASC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if date != "" {
		if t, err := time.Parse("2006-01-02", date); err == nil {
			start := t
			end := t.Add(24 * time.Hour)
			q = q.Where("date >= ? AND date < ?", start, end)
		}
	}
	err := q.Find(&releases).Error
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func (r *releaseRepository) DeleteById(ID uint) error {
	return r.db.Delete(&models.Release{ID: ID}).Error
}

func (r *releaseRepository) GetStatusesByRange(from, to time.Time) (map[string][]string, error) {
	type result struct {
		Day    time.Time
		Status string
	}

	var rows []result
	err := r.db.
		Model(&models.Release{}).
		Select("DATE(`date`) as day, status").
		Where("date >= ? AND date < ?", from, to).
		Group("day, status").
		Order("day ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	res := make(map[string][]string)
	seen := make(map[string]map[string]bool)

	for _, row := range rows {
		key := row.Day.Format("2006-01-02")
		if seen[key] == nil {
			seen[key] = map[string]bool{}
		}
		if !seen[key][row.Status] {
			res[key] = append(res[key], row.Status)
			seen[key][row.Status] = true
		}
	}

	return res, nil
}
