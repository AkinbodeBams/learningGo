package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		GetById(ctx context.Context,id int) (*Post,error)
		Create(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *User) error
	}

}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
	}
}