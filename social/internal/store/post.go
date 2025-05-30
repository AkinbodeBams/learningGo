package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)


type Post struct {
	ID int64 `json:"id"`
	Content string `json:"content"`
	Title string `json:"title"`
	UserId int64 `json:"user_id"`
	Tags []string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Version int `json:"version"`
	Comments []Comment `json:"comment"`
	User User `json:"user"`

}
 

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}
type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context,post *Post) error{
	query := `
	INSERT INTO posts (content,title,user_id,tags)
	VALUES ($1, $2, $3 , $4) RETURNING id, created_at, updated_at
	`
	ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

err:=	s.db.QueryRowContext(ctx, 
		query, 
		post.Content,
		post.Title,
		post.UserId, 
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err!= nil{
		return err
	}

	return nil
}

func (s *PostStore) GetById(ctx context.Context,id int) (*Post,error){
	var post Post
	query := `
	SELECT id,content,title,user_id,tags,created_at,updated_at FROM posts WHERE id=$1
	`
ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx,query,id,
	).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserId,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err!= nil{
		switch {
		case errors.Is(err,sql.ErrNoRows):
			return nil , ErrNotFound
		default:
			return nil , err
		}
		
	}



	return &post,nil
}

func(s *PostStore)  Delete(ctx context.Context,id int) (error){
	query := `
	DELETE FROM posts WHERE id=$1
	`
	ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return  err
	}

	rows, err := res.RowsAffected()

	if err != nil {	
		return  err
	}

	if rows == 0 {
		return fmt.Errorf("delete failed: %w", ErrNotFound)
	}
	
	return nil
}

func (s *PostStore) Update(ctx context.Context,post *Post)error {
	query := `
	UPDATE posts
	SET title = $1, content = $2 , version = version + 1
	WHERE id = $3 AND version = $4 RETURNING version
	`
	ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	var newVersion int
err := s.db.QueryRowContext(ctx, query, post.Title,post.Content,post.ID,post.Version).Scan(post.Version)
	if err != nil {
		switch {
		case errors.Is(err,sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
		
	}
	post.Version = newVersion
	return nil
} 

func (s *PostStore)GetUserFeed(ctx context.Context,userID int, fq PaginatedFeedQuery)([]PostWithMetadata, error){
 query := `
  SELECT 
		p.id, p.user_id, p.title ,p.content, p.created_at , p.version , p.tags,u.username ,
		COUNT(c.id) AS comment_count
	FROM posts p 
	LEFT JOIN comments c ON c.post_id = p.id
	LEFT JOIN users u ON p.user_id = u.id
	JOIN followers f ON f.follower_id = p.user_id  OR p.user_id = $1
	WHERE 
		f.user_id =$1  AND
		(p.title ILike '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
		(p.tags @> $5 OR $5 = '{}')
	GROUP BY p.id, u.username
	ORDER BY p.created_at ` + fq.Sort + `
	LIMIT $2 OFFSET $3
	`
 
ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	rows,err := s.db.QueryContext(ctx,query,userID,fq.Limit,fq.Offset,fq.Search,pq.Array(fq.Tags))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
var feed []PostWithMetadata
for rows.Next(){
	var post PostWithMetadata
	err := rows.Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.Version,
		pq.Array(&post.Tags),
		&post.User.Username,
		&post.CommentsCount,


	)

	if err != nil {
		return nil, err
	}
	feed=append(feed, post)
}
return feed,nil
}
 