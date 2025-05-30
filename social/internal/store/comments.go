package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID int64 `json:"id"`
	PostId int64 `json:"post_id"`
	UserId int64 `json:"user_id"`
	Content string  `json:"content"`
	CreatedAt string  `json:"created_at"`
	User User`json:"user"`

}


type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context,comment *Comment) error{
	query := `
	INSERT INTO comments (post_id, user_id , content)
	VALUES ($1,$2,$3)
	returning id, created_at
	` 
	ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err:= s.db.QueryRowContext(
		ctx,query,
		comment.PostId,
		comment.UserId,
		comment.Content, ).
	Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func(s *CommentStore) GetByPostID(ctx context.Context, postId int64)([]Comment, error){
query:= `
SELECT c.id, c.post_id, c.user_id , c.content, c.created_at ,  
users.username,users.id
FROM comments c
JOIN users on users.id = c.user_id
WHERE c.post_id = $1
Order BY c.created_at DESC;
`
ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
rows, err := s.db.QueryContext(ctx, query, postId)
 if err != nil {
	return nil, err
 }

 defer rows.Close()

 comments := []Comment{}

 for rows.Next(){
	var c Comment
	c.User = User{}
	err :=  rows.Scan(&c.ID,&c.PostId,&c.UserId,&c.Content,&c.CreatedAt,&c.User.Username,&c.User.ID)
	if err != nil {
		return nil, err
	}
	comments = append(comments,c)
	
 }
 return comments,nil
}

