package repostiry

import (
	"release-calendar/backend/internal/models"

	"gorm.io/gorm"
)

type LinkRepository interface {
	Add(links []*models.Link) error
	Update(link *models.Link) error
	Delete(links []*models.Link) error
}

type linkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(
	db *gorm.DB,
) LinkRepository {
	return &linkRepository{
		db: db,
	}
}

func (r *linkRepository) Add(links []*models.Link) error {
	return r.db.Create(links).Error
}

func (r *linkRepository) Update(link *models.Link) error {
	return r.db.Save(link).Error
}

func (r *linkRepository) Delete(links []*models.Link) error {
	return r.db.Delete(links).Error
}
