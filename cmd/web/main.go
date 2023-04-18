package main

import (
	"flag"
	"net/http"

	"github.com/gopher-bell/site/log"
	"go.uber.org/zap"
)

type application struct {
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	log.Init()
	defer log.Quit()

	log.ZapLogger.Info("starting http server...", zap.String("addr", *addr))

	app := &application{}
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.ZapLogger.Error("failed to start http server", zap.Error(err))
	}
}
