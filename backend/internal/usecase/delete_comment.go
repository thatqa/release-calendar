package usecase

import (
	"fmt"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type CommentDeleter struct {
	commentRepository repostiry.CommentRepository
}

func NewCommentDeleter(db *gorm.DB) CommentDeleter {
	return CommentDeleter{
		commentRepository: repostiry.NewCommentRepository(db),
	}
}

func (u *CommentDeleter) DeleteComment(ID uint) error {
	if err := u.commentRepository.DeleteById(ID); err != nil {
		return fmt.Errorf("failed to delete comment")
	}
	return nil
}
