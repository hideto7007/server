package common

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	file, err := os.OpenFile("takuwaeru.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}
