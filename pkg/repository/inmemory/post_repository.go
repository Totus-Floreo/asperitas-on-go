package repository

import (
	"sync"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"
	uuid "github.com/google/uuid"
)

type PostStorage struct {
	Storage []*model.Post
	mu      *sync.RWMutex
}

func NewPostStorage() *PostStorage {
	return &PostStorage{
		Storage: make([]*model.Post, 0),
		mu:      new(sync.RWMutex),
	}
}

func (s *PostStorage) GetAllPosts() ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Storage, nil
}

func (s *PostStorage) GetPostByID(postID string) (*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, post := range s.Storage {
		if post.ID == postID {
			return post, nil
		}
	}
	return nil, model.ErrPostNotFound
}

func (s *PostStorage) GetPostsByCategory(categoty string) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filtredPosts := Filter(s.Storage, func(post *model.Post) bool {
		return post.Categoty == categoty
	})
	return filtredPosts, nil
}

func (s *PostStorage) GetPostsByUser(userName string) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filtredPosts := Filter(s.Storage, func(post *model.Post) bool {
		return post.Author.Username == userName
	})
	return filtredPosts, nil
}

func (s *PostStorage) AddPost(post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	post.ID = uuid.New().String()
	s.Storage = append(s.Storage, post)

	return nil
}

func (s *PostStorage) DeletePost(postID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var postIdx int
	for idx, post := range s.Storage {
		if post.ID == postID {
			postIdx = idx
		}
	}

	if len(s.Storage) <= 1 {
		s.Storage = []*model.Post{}
		return nil
	}

	if postIdx != len(s.Storage)-1 {
		copy(s.Storage[postIdx:], s.Storage[postIdx+1:])
	}
	s.Storage[len(s.Storage)-1] = nil
	s.Storage = s.Storage[:len(s.Storage)-1]

	return nil
}

func (s *PostStorage) UpdatePost(post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for idx, p := range s.Storage {
		if p.ID == post.ID {
			s.Storage[idx] = post
			return nil
		}
	}
	return model.ErrPostNotFound
}

func (s *PostStorage) AddView(post *model.Post) error {
	post.Views++
	return nil
}

func (s *PostStorage) AddComment(post *model.Post, comment *model.Comment) error {
	post.CM.Lock()
	defer post.CM.Unlock()

	comment.ID = uuid.New().String()
	post.Comments = append(post.Comments, comment)

	return nil
}

func (s *PostStorage) GetCommentIdx(post *model.Post, commentID string) (int, error) {
	var commentIdx int
	for idx, comment := range post.Comments {
		if comment.ID == commentID {
			commentIdx = idx
			return commentIdx, nil
		}
	}
	return 0, model.ErrCommentNotFound
}

func (s *PostStorage) DeleteComment(post *model.Post, commentIdx int) error {
	post.CM.Lock()
	defer post.CM.Unlock()

	if len(post.Comments) == 1 {
		post.Comments = []*model.Comment{}
		return nil
	}

	if commentIdx != len(post.Comments)-1 {
		copy(post.Comments[commentIdx:], post.Comments[commentIdx+1:])
	}
	post.Comments[len(post.Comments)-1] = nil
	post.Comments = post.Comments[:len(post.Comments)-1]

	return nil
}

func (s *PostStorage) Vote(post *model.Post, vote *model.Vote) error {
	post.VM.Lock()
	defer post.VM.Unlock()

	for idx, vt := range post.Votes {
		if vt.UserID == vote.UserID {
			post.Votes[idx] = vote
			return nil
		}
	}

	post.Votes = append(post.Votes, vote)

	return nil
}

func (s *PostStorage) UnVote(post *model.Post, userID string) error {
	post.VM.Lock()
	defer post.VM.Unlock()

	var removedVote int
	for idx, vote := range post.Votes {
		if vote.UserID == userID {
			removedVote = idx
		}
	}

	if len(post.Votes) == 1 {
		post.Votes = []*model.Vote{}
		return nil
	}

	copy(post.Votes[removedVote:], post.Votes[removedVote+1:])
	post.Votes[len(post.Votes)-1] = nil
	post.Votes = post.Votes[:len(post.Votes)-1]

	return nil
}

func (s *PostStorage) UpdateScore(post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	positiveVotes := 0
	post.Score = 0
	for _, vote := range post.Votes {
		switch vote.Score {
		case 1:
			positiveVotes++
			post.Score++
		case -1:
			post.Score--
		}
	}

	post.UpvotePercentage = uint8((float64(positiveVotes) / float64(len(post.Votes))) * 100)
	return nil
}

func Filter(posts []*model.Post, fn func(*model.Post) bool) []*model.Post {
	result := []*model.Post{}
	for _, post := range posts {
		if fn(post) {
			result = append(result, post)
		}
	}
	return result
}
