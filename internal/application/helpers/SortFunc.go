package helpers

import "github.com/Totus-Floreo/asperitas-on-go/internal/model"

func LessByVotesThenViews(p1, p2 *model.Post) bool {
	if p1.Score != p2.Score {
		return p1.Score > p2.Score
	}
	return p1.Views > p2.Views
}
