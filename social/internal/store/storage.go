package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)
var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("duplicate") 
QueryTimeOutDuration = time.Second * 5
ErrDuplicateEmail =  errors.New("duplicate Email") 
ErrDuplicateUsername =  errors.New("duplicate Username") 


)



type Storage struct {
	Posts interface {
		GetById(ctx context.Context,id int) (*Post,error)
		Delete(ctx context.Context,id int) (error)
		Create(context.Context, *Post) error
		Update(context.Context, *Post)error
		GetUserFeed(ctx context.Context, userID int, pq PaginatedFeedQuery)([]PostWithMetadata, error)
	}

	Users interface {
		Create(context.Context,*sql.Tx, *User) error
		GetById(ctx context.Context,id int64) (*User,error)
		GetByEmail(ctx context.Context,email string) (*User,error)
		Delete(ctx context.Context,id int) (error)
		Update(ctx context.Context,user *User) error
		CreateAndInvite(ctx context.Context , user *User , token string, exp time.Duration) error
		Activate (ctx context.Context,token string) error
		GetUserByRoleID(ctx context.Context ,roleId int64) (*User,error)
	}

	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(ctx context.Context, postId int64)([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerUser int64, userId int)error
		UnFollow(ctx context.Context, followerUser int, userId int)error
	}

	Roles interface {
		GetByName(ctx context.Context ,role string) (*Role,error)
		
	}

}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
		Comments: &CommentStore{db},
		Followers: &FollowerStore{db},
		Roles: &RoleStore{db},
	}
}

func withTx(db *sql.DB , ctx context.Context ,  fn func(*sql.Tx)error)error{
tx, err := db.BeginTx(ctx,nil)
if err != nil {
	return err
}
if err:= fn(tx); err != nil {
	_ = tx.Rollback()
	return err
}
return tx.Commit()
}