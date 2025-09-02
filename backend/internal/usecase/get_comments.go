package usecase

import (
	"fmt"
	"release-calendar/backend/internal/models"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type CommentGetter struct {
	commentRepository repostiry.CommentRepository
}

func NewCommentGetter(db *gorm.DB) CommentGetter {
	return CommentGetter{
		commentRepository: repostiry.NewCommentRepository(db),
	}
}

func (u *CommentGetter) GetComments(releaseID uint) ([]*models.Comment, error) {
	comments, err := u.commentRepository.GetByReleaseID(releaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments")
	}

	return comments, nil
}
