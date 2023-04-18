package main

import (
	"database/sql"
	"flag"
	"net/http"

	"github.com/gopher-bell/site/internal/models"
	"github.com/gopher-bell/site/log"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	snippets *models.SnippetModel
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

	app := &application{
		snippets: &models.SnippetModel{DB: db},
	}
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
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