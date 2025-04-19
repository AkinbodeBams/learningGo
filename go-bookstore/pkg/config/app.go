package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (

	db * gorm.DB
)

func Connect(){
	dsn := "host=localhost port=5432 user=postgres dbname=bookrest password=1855 sslmode=disable"
	d,err := gorm.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB{
	return db
}