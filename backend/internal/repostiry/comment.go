package repostiry

import (
	"release-calendar/backend/internal/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	GetById(ID uint) (*models.Comment, error)
	GetByReleaseID(releaseID uint) ([]*models.Comment, error)
	Add(comment *models.Comment) error
	Update(comment *models.Comment) error
	DeleteById(ID uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(
	db *gorm.DB,
) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) GetByReleaseID(releaseID uint) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := r.db.
		Where("release_id = ?", releaseID).
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *commentRepository) Add(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetById(ID uint) (*models.Comment, error) {
	var comment *models.Comment
	err := r.db.Where("id = ?", ID).First(&comment).Error
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) DeleteById(ID uint) error {
	return r.db.Delete(&models.Comment{ID: ID}).Error
}
