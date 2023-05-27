package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Body    string             `json:"body" bson:"body"`
	Created string             `json:"created" bson:"created"`
	Author  *Author            `json:"author" bson:"author"`
}

const layout = "2006-01-02T15:04:05.000Z"

func NewComment(body string, author *Author) *Comment {
	return &Comment{
		Body:    body,
		Created: time.Now().UTC().Format(layout),
		Author:  author,
	}
}
