package log

import (
	"log"

	"go.uber.org/zap"
)

var ZapLogger *zap.Logger

func Init() {
	var err error
	ZapLogger, err = zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}
}

func Quit() {
	ZapLogger.Sync()
}
