package helpers

import "github.com/Totus-Floreo/asperitas-on-go/internal/model"

func FindCommentIdx(post *model.Post, commentID string) (int, error) {
	for idx, comment := range post.Comments {
		if comment.ID.Hex() == commentID {
			return idx, nil
		}
	}
	return 0, model.ErrCommentNotFound
}
