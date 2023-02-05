package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"student-management/handler"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/nosurf"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var sessionManager *scs.SessionManager

var schema = `
CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
	status BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL,

	PRIMARY KEY(id),
	UNIQUE(username),
	UNIQUE(email)
);

CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
`

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	decoder := form.NewDecoder()
	
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetString("database.host"),
		config.GetString("database.port"),
		config.GetString("database.user"),
		config.GetString("database.password"),
		config.GetString("database.dbname"),
		config.GetString("database.sslmode"),
	))
	if err != nil {
		log.Fatalln(err)
	}

	res := db.MustExec(schema)
	row, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}

	if row < 0 {
		log.Fatalln("failed to run schema")
	}

	lt := config.GetDuration("session.lifetime")
	it := config.GetDuration("session.idletime")
	sessionManager = scs.New()
	sessionManager.Lifetime = lt * time.Hour
	sessionManager.IdleTimeout = it * time.Minute
	sessionManager.Cookie.Name = "web-session"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = true
	sessionManager.Store = NewSQLXStore(db)
	

	chi := handler.NewHandler(sessionManager, decoder, db)
	p := config.GetInt("server.port")
	newChi := nosurf.New(chi)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", p), newChi); err != nil {
		log.Fatalf("%#v", err)
	}
}
