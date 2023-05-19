package model

import (
	"sync"
	"time"
)

type Post struct {
	Categoty string `json:"category"`
	Text     string `json:"text"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Url      string `json:"url"`

	ID     string  `json:"id"`
	Author *Author `json:"author"`

	Created          string `json:"created"`
	UpvotePercentage uint8  `json:"upvotePercentage"`
	Views            int64  `json:"views"`
	Score            int64  `json:"score"`

	Comments []*Comment  `json:"comments"`
	CM       *sync.Mutex `json:"-"`

	Votes []*Vote     `json:"votes"`
	VM    *sync.Mutex `json:"-"`
}

func NewPost() *Post {
	return &Post{
		Created:          time.Now().UTC().Format(layout),
		Views:            0,
		Score:            0,
		UpvotePercentage: 0,
		Comments:         []*Comment{},
		CM:               new(sync.Mutex),
		Votes:            []*Vote{},
		VM:               new(sync.Mutex),
	}
}

type Vote struct {
	UserID string `json:"user"`
	Score  int8   `json:"vote"`
}

type IPostStorage interface {
	GetAllPosts() ([]*Post, error)
	GetPostByID(string) (*Post, error)
	GetPostsByCategory(string) ([]*Post, error)
	GetPostsByUser(string) ([]*Post, error)
	AddPost(*Post) error
	DeletePost(string) error
	UpdatePost(*Post) error
	AddView(*Post) error
	AddComment(*Post, *Comment) error
	GetCommentIdx(*Post, string) (int, error)
	DeleteComment(*Post, int) error
	Vote(*Post, *Vote) error
	UnVote(*Post, string) error
	UpdateScore(*Post) error
}
