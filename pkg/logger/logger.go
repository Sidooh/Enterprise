package logger

import (
	"enterprise.sidooh/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ClientLog = &log.Logger{
	Out: nil,
}

func Init() {
	ClientLog = log.New()

	env := viper.GetString("APP_ENV")
	if env != "TEST" {

		file := utils.GetLogFile("client.log")
		ClientLog.SetOutput(file)

	}
}
