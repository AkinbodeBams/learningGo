package cache

import (
	"context"

	"github.com/akinbodeBams/social/internal/store"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(ctx context.Context, userId int64)(*store.User, error)
		Set(ctx context.Context, user *store.User)( error)
	}
}

func NewRedisStorage(rdb *redis.Client) Storage{
	return Storage{
Users: &UserStore{rdb:rdb},
	}
}