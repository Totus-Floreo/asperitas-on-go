package model

import "context"

type ICursor interface {
	All(context.Context, interface{}) error
	Decode(interface{}) error
	TryNext(context.Context) bool
	Err() error
	Close(context.Context) error
}
