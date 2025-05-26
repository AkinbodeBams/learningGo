package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/akinbodeBams/social/docs"
	"github.com/akinbodeBams/social/internal/auth"
	"github.com/akinbodeBams/social/internal/store"
	"github.com/akinbodeBams/social/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type config struct {
	addr string
	db dbConfig
	env string
	apiURL string
	mail mailConfig
	auth  authConfig
	redisCfg redisConfig

}
type redisConfig struct {
	addr string
	pw string
	db int
	enabled bool
}

type mailConfig struct {
	exp time.Duration
}
type application struct {
	config config
	store  store.Storage
	cacheStorage cache.Storage
	logger *zap.SugaredLogger
	// mailer mailer.Client
		authenticator auth.Authenticator
}

type authConfig struct{
	basic basicConfig 
	token tokenConfig
}

type tokenConfig struct {
secret string
exp time.Duration
iss string
}

type basicConfig struct {
	user string 
	pass string
}

type dbConfig struct{
	addr string
	maxOpenConns  int
	maxIdleConns int
	maxIdleTime string
}


func (app *application) mount() http.Handler{
	
	docs.SwaggerInfo.Version =version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/v1", func(r chi.Router) {
docsUrl := fmt.Sprintf("http://%s/v1/swagger/doc.json", app.config.apiURL)
r.Mount("/swagger", httpSwagger.Handler(httpSwagger.URL(docsUrl)))
		r.With(app.BasicAuthMiddleware()).Get("/health",app.healthCheckHandler)
		
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/",app.getPostHandler)

				r.Patch("/", app.checkPostOwnership("moderator",app.updatePostHandler))
				r.Delete("/", app.checkPostOwnership("admin",app.deletePostHandler))
			})
		}) 

		r.Route("/users",func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserhandler)
			
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/",app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unFollowUserHandler)
				r.Delete("/",app.deleteUserHandler)
				r.Patch("/", app.updateUserHandler)
				
			})
			r.Group(func(r chi.Router) {
					r.Use(app.AuthTokenMiddleware)
					r.Get("/feed",app.getUserFeedHandler)
				})
		})
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user",app.registerUserHandler)
			r.Post("/token",app.createTokenHandler)
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
app.logger.Infow("Server started", "addr", app.config.addr , "env", app.config.env)
	return srv.ListenAndServe()
}