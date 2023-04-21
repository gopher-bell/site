package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"net/http"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/gopher-bell/site/internal/models"
	"github.com/gopher-bell/site/log"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root:1234@tcp(localhost:3306)/snippetbox?parseTime=True", "MySQL data source name")
	flag.Parse()

	log.Init()
	defer log.Quit()

	log.ZapLogger.Info("starting http server", zap.String("addr", *addr))

	db, err := openDB(*dsn)
	if err != nil {
		log.ZapLogger.Fatal("failed to open db connection", zap.Error(err))
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		log.ZapLogger.Fatal("failed to cache template", zap.Error(err))
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		log.ZapLogger.Fatal("failed to start http server", zap.Error(err))
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
