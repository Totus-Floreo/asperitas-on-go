package application

import (
	"context"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/Totus-Floreo/asperitas-on-go/internal/application/helpers"
	"github.com/Totus-Floreo/asperitas-on-go/internal/middleware"
	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService struct {
	postStorage model.IPostStorage
}

func NewPostService(postStorage model.IPostStorage) *PostService {
	return &PostService{
		postStorage: postStorage,
	}
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	posts, err := s.postStorage.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(posts, func(i, j int) bool {
		return helpers.LessByVotesThenViews(posts[i], posts[j])
	})

	return posts, nil
}

func (s *PostService) GetPostByID(ctx context.Context, postID string) (*model.Post, error) {
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, model.ErrInvalidPostID
	}
	post, err := s.postStorage.GetPostByID(ctx, postObjectID)
	if err != nil {
		return nil, err
	}

	if err := s.postStorage.AddView(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPostsByCategory(ctx context.Context, category string) ([]*model.Post, error) {
	posts, err := s.postStorage.GetPostsByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	sort.Slice(posts, func(i, j int) bool {
		return helpers.LessByVotesThenViews(posts[i], posts[j])
	})

	return posts, nil
}

func (s *PostService) GetPostsByUser(ctx context.Context, userName string) ([]*model.Post, error) {
	posts, err := s.postStorage.GetPostsByUser(ctx, userName)
	if err != nil {
		return nil, err
	}

	sort.Slice(posts, func(i, j int) bool {
		return helpers.LessByVotesThenViews(posts[i], posts[j])
	})

	return posts, err
}

func (s *PostService) AddPost(ctx context.Context, post *model.Post) (*model.Post, error) {
	post.Author = ctx.Value(middleware.AuthorContextKey).(*model.Author)

	if post.Url != "" {
		post.Url = strings.TrimSpace(post.Url)
		if work := helpers.CheckLink(post.Url); !work {
			return nil, model.ErrInvalidUrl
		}
	}

	if err := s.postStorage.AddPost(ctx, post); err != nil {
		return nil, err
	} else {
		postCreated, err := s.postStorage.GetPostByID(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		return postCreated, nil
	}
}

func (s *PostService) DeletePost(ctx context.Context, postID string) error {
	author := ctx.Value(middleware.AuthorContextKey).(*model.Author)

	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return model.ErrInvalidPostID
	}

	if post, err := s.postStorage.GetPostByID(ctx, postObjectID); err != nil {
		return err
	} else if post.Author.ID != author.ID {
		return model.ErrUnAuthorized
	}

	err = s.postStorage.DeletePost(ctx, postObjectID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostService) AddComment(ctx context.Context, postID string, body string) (*model.Post, error) {
	author := ctx.Value(middleware.AuthorContextKey).(*model.Author)

	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, model.ErrInvalidPostID
	}

	if utf8.RuneCountInString(body) > 2000 {
		return nil, model.ErrCommentTooLong
	}

	post, err := s.postStorage.GetPostByID(ctx, postObjectID)
	if err != nil {
		return nil, err
	}

	comment := model.NewComment(body, author)
	if err := s.postStorage.AddComment(ctx, post, comment); err != nil {
		return nil, err
	}

	postChanged, err := s.postStorage.GetPostByID(ctx, postObjectID)
	if err != nil {
		return nil, err
	}

	return postChanged, nil

}

func (s *PostService) DeleteComment(ctx context.Context, postID string, commentID string) (*model.Post, error) {
	author := ctx.Value(middleware.AuthorContextKey).(*model.Author)

	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, model.ErrInvalidPostID
	}

	post, err := s.postStorage.GetPostByID(ctx, postObjectID)
	if err != nil {
		return nil, err
	}

	commendIdx := -1
	for idx, comment := range post.Comments {
		if comment.ID.Hex() == commentID {
			commendIdx = idx
		}
	}
	if commendIdx == -1 {
		return nil, model.ErrCommentNotFound
	}

	if post.Comments[commendIdx].Author.ID != author.ID {
		return nil, model.ErrUnAuthorized
	}

	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, model.ErrInvalidCommentID
	}

	if err = s.postStorage.DeleteComment(ctx, post, commentObjectID); err != nil {
		return nil, err
	}

	postChanged, err := s.postStorage.GetPostByID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return postChanged, nil
}

func (s *PostService) Vote(ctx context.Context, postID string, method string) (*model.Post, error) {
	author := ctx.Value(middleware.AuthorContextKey).(*model.Author)

	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, model.ErrInvalidPostID
	}

	post, err := s.postStorage.GetPostByID(ctx, postObjectID)
	if err != nil {
		return nil, err
	}

	switch method {
	case "upvote":
		vote := &model.Vote{
			UserID: author.ID,
			Score:  1,
		}
		if err = s.postStorage.Vote(ctx, post, vote); err != nil {
			return nil, err
		}
	case "downvote":
		vote := &model.Vote{
			UserID: author.ID,
			Score:  -1,
		}
		if err := s.postStorage.Vote(ctx, post, vote); err != nil {
			return nil, err
		}
	case "unvote":
		if err := s.postStorage.UnVote(ctx, post, author.ID); err != nil {
			return nil, err
		}
	default:
		return nil, model.ErrVotesActionNotImplement
	}

	postVoted, err := s.postStorage.GetPostByID(ctx, postObjectID)
	if err != nil {
		return nil, err
	}

	err = s.postStorage.UpdateScore(ctx, postVoted)
	if err != nil {
		return nil, err
	}

	postChanged, err := s.postStorage.GetPostByID(ctx, postObjectID)
	if err != nil {
		return nil, err
	}

	return postChanged, nil
}
