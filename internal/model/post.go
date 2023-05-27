package model

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Category string `json:"category" bson:"category"`
	Text     string `json:"text" bson:"text"`
	Title    string `json:"title" bson:"title"`
	Type     string `json:"type" bson:"type"`
	Url      string `json:"url" bson:"url"`

	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Author *Author            `json:"author" bson:"author"`

	Created          string `json:"created" bson:"created"`
	UpvotePercentage int64  `json:"upvotePercentage" bson:"upvotepercentage"`
	Views            int64  `json:"views" bson:"views"`
	Score            int64  `json:"score" bson:"score"`

	Comments []*Comment  `json:"comments" bson:"comments"`
	CM       *sync.Mutex `json:"-" bson:"-"`

	Votes []*Vote     `json:"votes" bson:"votes"`
	VM    *sync.Mutex `json:"-" bson:"-"`
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
	UserID string `json:"user" bson:"user"`
	Score  int64  `json:"vote" bson:"vote"`
}

type IPostStorage interface {
	GetAllPosts(context.Context) ([]*Post, error)
	GetPostByID(context.Context, primitive.ObjectID) (*Post, error)
	GetPostsByCategory(context.Context, string) ([]*Post, error)
	GetPostsByUser(context.Context, string) ([]*Post, error)
	AddPost(context.Context, *Post) error
	DeletePost(context.Context, primitive.ObjectID) error
	AddView(context.Context, *Post) error
	AddComment(context.Context, *Post, *Comment) error
	DeleteComment(context.Context, *Post, primitive.ObjectID) error
	Vote(context.Context, *Post, *Vote) error
	UnVote(context.Context, *Post, string) error
	UpdateScore(context.Context, *Post) error
}
