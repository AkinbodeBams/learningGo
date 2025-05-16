package main

import (
	"log"

	env "github.com/akinbodeBams/social/internal"
	"github.com/akinbodeBams/social/internal/db"
	"github.com/akinbodeBams/social/internal/store"
)
const version = "0.0.1"
func main() {

	cfg := config{
		addr: env.GetString("ADDR",":8080"),
		db: dbConfig{addr: env.GetString("DBADDR",":postgres://postgres:1855@localhost/social?sslmode=disable"),
				maxOpenConns:  env.GetInt("DB_Max_OPEN_CONNS",30),
				maxIdleTime:  env.GetString("DB_Max_OPEN_CONNS","15m"),
				maxIdleConns:  env.GetInt("DB_Max_OPEN_CONNS",30)},
				
		env: env.GetString("env","dev"),}
				
		
db,err:=db.New(cfg.db.addr,cfg.db.maxIdleTime,cfg.db.maxOpenConns,cfg.db.maxIdleConns)

if err != nil {
	log.Panic(err)
}
	store:= store.NewStorage(db)
	app := &application{
		config: cfg,
		store: store,
	}

	
	mux:= app.mount()

	log.Fatal(app.run(mux))
}