package main

import (
	"log"

	env "github.com/akinbodeBams/social/internal"
	"github.com/akinbodeBams/social/internal/db"
	"github.com/akinbodeBams/social/internal/store"
)

func main (){
	addr:= env.GetString("DB_ADDR", "postgres://postgres:1855@localhost/social?sslmode=disable")
	conn ,err:= db.New(addr,"15m",3,3)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	store:= store.NewStorage(conn)
	db.Seed(store,conn)
}