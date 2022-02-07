package backend

import "context"

type Backend interface {
	GetUserName(context.Context) (string, error)
}
