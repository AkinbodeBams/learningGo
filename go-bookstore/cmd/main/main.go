package main

import (
	"log"
	"net/http"

	"github.com/akinbodeBams/go-bookstore/pkg/routes"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main(){
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/",r)
	log.Fatal(http.ListenAndServe("localhost:8080",r))
}