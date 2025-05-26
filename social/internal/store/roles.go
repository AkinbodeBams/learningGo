package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Level int64 `json:"level"`
	Description  string `json:"description"`
}
 

type RoleStore struct {
	db *sql.DB}


func (s *RoleStore) GetByName(ctx context.Context ,roleName string) (*Role,error){
	var role Role
	query :=`
	SELECT name FROM ROLES WHERE name=$1
	` 
	ctx,cancel:= context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
 err := s.db.QueryRowContext(ctx, query, role).Scan(
	&role.Name,
)
 if err != nil {
	return nil, err
 }

 return &role,nil
 }
