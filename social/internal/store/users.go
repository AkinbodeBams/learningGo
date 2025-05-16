package store

import (
	"context"
	"database/sql"
)
type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Title string `json:"titile"`
	Email int64 `json:"email"`
	Password string `json:"password"`
	CreatedAt string `json:"created_id"`
	

}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context,user *User) error{
	query := `
	INSERT INTO users (username, password , email)
	VALUES ($1,$2,$3)
	returning id, created_at, 
	` 
	err:= s.db.QueryRowContext(
		ctx,query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(&user.ID,&user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}