package usecase

import (
	"fmt"
	"release-calendar/backend/internal/models"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type CommentUpdater struct {
	commentRepository repostiry.CommentRepository
}

func NewCommentUpdater(db *gorm.DB) CommentUpdater {
	return CommentUpdater{
		commentRepository: repostiry.NewCommentRepository(db),
	}
}

func (u *CommentUpdater) UpdateComment(ID uint, cmd CommentCMD) (*models.Comment, error) {
	comment, err := u.commentRepository.GetById(ID)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}

	comment.Author = cmd.Author
	comment.Message = cmd.Message

	if err := u.commentRepository.Update(comment); err != nil {
		return nil, fmt.Errorf("failed to update comment")
	}

	return comment, nil
}
