package logger

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

var ClientLog = &log.Logger{
	Out: nil,
}

func Init() {
	ClientLog = log.New()

	env := viper.GetString("APP_ENV")
	if env != "TEST" {

		filename := "pkg/logger/client-" + time.Now().Format("2006-01-02") + ".log"
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			ClientLog.SetOutput(os.Stderr)
		} else {
			ClientLog.SetOutput(file)
		}

	}
}
