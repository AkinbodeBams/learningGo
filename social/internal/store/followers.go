package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)


type Follower struct {
	UserID int64 `json:"user_id"`
	FollowerID string `json:"follower_id"`
	CreatedAt string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func( s *FollowerStore) Follow(ctx context.Context, followerId int64, userId int)error{
	query := `
	INSERT INTO followers (user_id, follower_id) VALUES($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)

	defer cancel()

	_, err := s.db.ExecContext(ctx,query,userId,followerId)
if err != nil {
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code =="23505"{
		return ErrConflict
	}
}


	return nil
}
func( s *FollowerStore) UnFollow(ctx context.Context, followerId int, userId int)error{
	query := `
	DELETE FROM followers 
	WHERE user_id = $1 AND follower_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)

	defer cancel()

	_, err := s.db.ExecContext(ctx,query,userId,followerId)

return err
}