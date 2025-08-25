package usecase

import (
	"fmt"
	"release-calendar/backend/internal/models"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type CommentCreator struct {
	commentRepository repostiry.CommentRepository
}

func NewCommentCreator(db *gorm.DB) CommentCreator {
	return CommentCreator{
		commentRepository: repostiry.NewCommentRepository(db),
	}
}

func (u *CommentCreator) AddComment(releaseID uint, cmd CommentCMD) (*models.Comment, error) {
	comment := &models.Comment{
		ReleaseID: releaseID,
		Author:    cmd.Author,
		Message:   cmd.Message,
	}
	if err := u.commentRepository.Add(comment); err != nil {
		return nil, fmt.Errorf("failed to add comment")
	}

	return comment, nil
}
