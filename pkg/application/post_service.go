package application

import (
	"errors"
	"net/http"
	"sort"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"
)

type PostService struct {
	postStorage model.IPostStorage
}

func NewPostService(postStorage model.IPostStorage) *PostService {
	return &PostService{
		postStorage: postStorage,
	}
}

func (s *PostService) GetAllPosts() ([]*model.Post, error) {
	posts, err := s.postStorage.GetAllPosts()
	if err != nil {
		return nil, err
	}

	sort.Slice(posts, func(i, j int) bool {
		return lessByVotesThenViews(posts[i], posts[j])
	})

	return posts, nil
}

func (s *PostService) GetPostByID(postID string) (*model.Post, error) {
	post, err := s.postStorage.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	if err := s.postStorage.AddView(post); err != nil {
		return post, err
	}
	return post, nil
}

func (s *PostService) GetPostsByCategory(category string) ([]*model.Post, error) {
	posts, err := s.postStorage.GetPostsByCategory(category)
	if err != nil {
		return nil, err
	}

	sort.Slice(posts, func(i, j int) bool {
		return lessByVotesThenViews(posts[i], posts[j])
	})

	return posts, nil
}

func (s *PostService) GetPostsByUser(userName string) ([]*model.Post, error) {
	posts, err := s.postStorage.GetPostsByUser(userName)
	if err != nil {
		return nil, err
	}

	sort.Slice(posts, func(i, j int) bool {
		return lessByVotesThenViews(posts[i], posts[j])
	})

	return posts, err
}

func (s *PostService) AddPost(post *model.Post) (*model.Post, error) {
	if post.Url != "" {
		if work := checkLink(post.Url); !work {
			return nil, model.ErrInvalidUrl
		}
	}
	if err := s.postStorage.AddPost(post); err != nil {
		return nil, err
	} else {
		postCreated, err := s.postStorage.GetPostByID(post.ID)
		if err != nil {
			return nil, err
		}
		return postCreated, nil
	}
}

func (s *PostService) DeletePost(postID string, author *model.Author) (string, error) {
	if post, err := s.postStorage.GetPostByID(postID); err != nil {
		return `{"message": "post not found"}`, err
	} else if post.Author.ID != author.ID {
		return `{"message":"unauthorized"}`, model.ErrUnAuthorized
	}
	err := s.postStorage.DeletePost(postID)
	if err != nil {
		return `{"message":"bad connection to db"}`, err
	}
	return `{"message":"success"}`, nil
}

func (s *PostService) AddComment(postID string, body string, author *model.Author) (*model.Post, error) {
	post, err := s.postStorage.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	comment := model.NewComment(body, author)
	if err := s.postStorage.AddComment(post, comment); err != nil {
		return nil, err
	}

	postChanged, err := s.postStorage.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	return postChanged, nil

}

func (s *PostService) DeleteComment(postID string, commendID string, author *model.Author) (*model.Post, error) {
	post, err := s.postStorage.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	commendIdx, err := s.postStorage.GetCommentIdx(post, commendID)
	if err != nil {
		return nil, err
	}

	if post.Comments[commendIdx].Author.ID != author.ID {
		return nil, model.ErrUnAuthorized
	}

	if err = s.postStorage.DeleteComment(post, commendIdx); err != nil {
		return nil, err
	}

	postChanged, err := s.postStorage.GetPostByID(post.ID)
	if err != nil {
		return nil, err
	}

	return postChanged, nil
}

func (s *PostService) Vote(postID string, author *model.Author, method string) (*model.Post, error) {
	post, err := s.postStorage.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	switch method {
	case "upvote":
		vote := &model.Vote{
			UserID: author.ID,
			Score:  1,
		}
		if err = s.postStorage.Vote(post, vote); err != nil {
			return nil, err
		}
	case "downvote":
		vote := &model.Vote{
			UserID: author.ID,
			Score:  -1,
		}
		if err := s.postStorage.Vote(post, vote); err != nil {
			return nil, err
		}
	case "unvote":
		if err := s.postStorage.UnVote(post, author.ID); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("not implement vote action")
	}

	err = s.postStorage.UpdateScore(post)
	if err != nil {
		return nil, err
	}

	postChanged, err := s.postStorage.GetPostByID(postID)
	if err != nil {
		return nil, err
	}
	return postChanged, nil
}

func checkLink(link string) bool {
	resp, err := http.Get(link)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func lessByVotesThenViews(p1, p2 *model.Post) bool {
	if p1.Score != p2.Score {
		return p1.Score > p2.Score
	}
	return p1.Views > p2.Views
}
