package main

import (
	"time"

	_ "github.com/akinbodeBams/social/docs"
	env "github.com/akinbodeBams/social/internal"
	"github.com/akinbodeBams/social/internal/auth"
	"github.com/akinbodeBams/social/internal/db"
	"github.com/akinbodeBams/social/internal/store"
	"go.uber.org/zap"
)
const version = "0.0.1"

type Error struct {
	Error string `json:"error"`
}
//	@title			Social API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html


//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description
func main() {

	cfg := config{
		addr: env.GetString("ADDR",":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{addr: env.GetString("DBADDR",":postgres://postgres:1855@localhost/social?sslmode=disable"),
		
				maxOpenConns:  env.GetInt("DB_Max_OPEN_CONNS",30),
				maxIdleTime:  env.GetString("DB_Max_OPEN_CONNS","15m"),
				maxIdleConns:  env.GetInt("DB_Max_OPEN_CONNS",30)},
				
		env: env.GetString("env","dev"),mail: mailConfig{
			exp: time.Hour * 24 *3 ,// 3 days 
		},
		auth: authConfig{
			basic:basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_USER", "admin"),
			},
			token: tokenConfig{
				secret: "example",
				exp: time.Hour * 24 * 3,
				iss: "qwerty",
			},
		},
}
				
// loger
logger:= zap.Must((zap.NewProduction())).Sugar()
defer logger.Sync()

db,err:=db.New(cfg.db.addr,cfg.db.maxIdleTime,cfg.db.maxOpenConns,cfg.db.maxIdleConns)

if err != nil {
	logger.Fatal(err)
}

defer db.Close()

logger.Info("database connection pool established")
	store:= store.NewStorage(db)

	jwtAuthenticator:=auth.NewJWTAuthenticator(cfg.auth.token.secret,cfg.auth.token.iss, cfg.auth.token.iss )
	app := &application{
		config: cfg,
		store: store,
		logger: logger,
		authenticator:jwtAuthenticator,
	}

	
	mux:= app.mount()

	logger.Fatal(app.run(mux))
}