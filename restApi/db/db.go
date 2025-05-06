package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var DB  *sql.DB

func InitDb(){

var err error
DB, err = sql.Open("sqlite3", "api.db")

if err!= nil {
	fmt.Printf("%v",err)
	panic("Could not connect to db")
}

DB.SetMaxOpenConns(10)
DB.SetConnMaxIdleTime(5)
createTables()
}



func createTables(){
	createEventTables := `
	CREATE TABLE IF NOT EXISTS events(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	dateTime DATETIME NOT NULL,
	user_id INTEGER
	)
	`

fmt.Println("log it all out ", createEventTables)


	_, err := DB.Exec(createEventTables)

	if err != nil {
		fmt.Printf("Could not create event table: %v", err)
	}
	// n:=0
	// fmt.Print(n)
	var m sync.Mutex
}