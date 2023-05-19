package model

import "time"

type Comment struct {
	ID      string  `json:"id"`
	Body    string  `json:"body"`
	Created string  `json:"created"`
	Author  *Author `json:"author"`
}

const layout = "2006-01-02T15:04:05.000Z"

func NewComment(body string, author *Author) *Comment {
	return &Comment{
		Body:    body,
		Created: time.Now().UTC().Format(layout),
		Author:  author,
	}
}
