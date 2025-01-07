package common

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger(logFile string) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.WithError(err).Fatal("ログファイルのオープンに失敗しました")
	}

	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}
