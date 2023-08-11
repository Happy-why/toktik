package repo

import "context"

type UserRepo interface {
	GetUserByUsername(c context.Context, username string) (bool, error)
}
