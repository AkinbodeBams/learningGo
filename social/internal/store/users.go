package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)
type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password Password `json:"-"`
	CreatedAt string `json:"created_at"`
	IsActive bool `json:"is_active"`
	Version int `json:"version"`
	RoleID int `json:"role_id"`
	Role Role `json:"role"`
}
type Password struct {
text *string
hash []byte
}

func (p *Password) Set(text string)error{
hash , err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
if err != nil {
	return err
}
p.text = &text
p.hash = hash
return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context,tx *sql.Tx,user *User) error{
	query := `
	INSERT INTO users (username, password , email,role_id)
	VALUES ($1,$2,$3,(SELECT id FROM roles WHERE name = $4))
	returning id, created_at
	` 
	ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	role:= user.Role.Name

	if role == ""{
		role = "user"
	}
	err:= s.db.QueryRowContext(
		ctx,query,
		user.Username,
		user.Password.hash,
		user.Email,
		role, 
	).Scan(&user.ID,&user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique costraint "user_email_key`:
			return  ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique costraint "user_username_key`:
			return  ErrDuplicateUsername
		default:
			return err
		
		}
	}

	return nil
}

func (s *UserStore) GetById(ctx context.Context,id int64) (*User,error){
	var user User
	query := `
	SELECT id,email,username,created_at FROM users WHERE id=$1
	`
ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx,query,id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
	)
	if err!= nil{
		switch {
		case errors.Is(err,sql.ErrNoRows):
			return nil , ErrNotFound
		default:
			return nil , err
		}
	}
	return &user,nil
}

func (s *UserStore) GetByEmail(ctx context.Context,email string) (*User,error){
	query := `
	SELECT id, username, email , password , created_at FROM users
	WHERE email = $1 AND is_active = true
	`;
	var user User
	ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx,query,email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)
	if err!= nil{
		switch {
		case errors.Is(err,sql.ErrNoRows):
			return nil , ErrNotFound
		default:
			return nil , err
		}
	}
	return &user,nil
}

func(s *UserStore)  Delete(ctx context.Context,id int) (error){
	query := `
	DELETE FROM users WHERE id=$1
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

func (s *UserStore) Update(ctx context.Context,user *User)error {
	return withTx(s.db, ctx,func(tx *sql.Tx) error {
query := `
	UPDATE users
	SET email = $1, username = $2 , version = version + 1, is_active = $3
	WHERE id = $4 AND version = $5 RETURNING version
	`
	ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	var newVersion int
err := s.db.QueryRowContext(ctx, query,
	user.Email,
	user.Username,
	user.IsActive,
	user.ID,
	user.Version,
).Scan(&newVersion)

if err != nil {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	default:
		return err
	}
}

user.Version = newVersion
return nil

	})
	
} 


func (s *UserStore) CreateAndInvite(ctx context.Context , user *User , token string , invitationExp time.Duration) error {
	return withTx(s.db,ctx,func(tx *sql.Tx) error {
		if err:= s.Create(ctx,tx,user); err!= nil{
			return err
		}

		err := s.createUserInvitation(ctx,tx, token , invitationExp,user.ID);
		if err != nil {
			return  err
		}
		return nil
	})
}


func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string,exp time.Duration, userID int64)error {
query := `INSERT INTO user_invitations (token , user_id , expiry)
			VALUES ($1,$2,$3)`
ctx, cancel :=  context.WithTimeout(ctx, QueryTimeOutDuration)

defer cancel()

_,err := tx.ExecContext(ctx,query,token,userID,time.Now().Add(exp))

if err != nil {
	return err
}

return nil
}

func (s *UserStore) Activate (ctx context.Context,token string)error{
return withTx(s.db, ctx , func(tx *sql.Tx) error {
	user, err:= s.getUserFromInvitation(ctx, tx , token)
	if err != nil {
		return  err
	}

user.IsActive = true

if err := s.Update(ctx, user); err!= nil {
	return err
}

if err:= s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
	return err
}
	return nil
})
}


func (s *UserStore) getUserFromInvitation(ctx context.Context , tx *sql.Tx , token string ) (*User, error){
 query := `SELECT u.id, u.username , u.email , u.created_at , u.is_active 
 FROM  users u
 JOIN user_invitations ui on u.id = ui.user_id
 where ui.token = $1 AND ui.expiry > $2`
 hash:= sha256.Sum256([]byte(token))
hashToken := hex.EncodeToString(hash[:])
 ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)

 defer cancel()
exp:= time.Now()


 user := &User{ }
 err := tx.QueryRowContext(ctx, query, hashToken,exp ).Scan(
	&user.ID,
	&user.Username,
	&user.Email,
	&user.CreatedAt,
	&user.IsActive,
 )

 if err != nil {
	switch err {
	case sql.ErrNoRows:
		return nil , ErrNotFound
	default:
		return nil, err
	}
 }

 return user, nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context , tx *sql.Tx , userId int64 ) ( error){
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	
	ctx, cancel :=  context.WithTimeout(ctx,QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,userId)

	if err != nil {
		return  err
	}
	return nil
}

func (s *UserStore) GetUserByRoleID(ctx context.Context , roleID int64 ) (*User,error){
	var user User
	query := `SELECT users.id, role_id, roles.*  FROM users
	JOIN roles ON (users.role_id = roles.id)
	WHERE role_id = $1`
	
ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
 err := s.db.QueryRowContext(ctx, query, roleID).Scan(
	&user.ID,
	&user.RoleID,
	&user.Role.ID,
	&user.Role.Name,
	&user.Role.Level,
	&user.Role.Description,
)
 if err != nil {
	return nil, err
 }

 return &user,nil
}