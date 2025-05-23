package main

import (
	"log"
	"net/http"
	"time"

	"github.com/akinbodeBams/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type config struct {
	addr string
	db dbConfig
	env string
}
type application struct {
	config config
	store  store.Storage
}

type dbConfig struct{
	addr string
	maxOpenConns  int
	maxIdleConns int
	maxIdleTime string
}


func (app *application) mount() http.Handler{
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/v1", func(r chi.Router) {

		r.Get("/health",app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Get("/{postID}", app.getPostHandler)
		})
	})


return r
}

func (app *application) run(mux http.Handler) error {
	
	srv := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}
log.Printf("Server started at %s", app.config.addr)
	return srv.ListenAndServe()
}