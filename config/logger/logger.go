package logger

import (
    "os"
    "github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger(mode string) {
    Log = logrus.New()
    Log.Out = os.Stdout

    if mode == "production" {
        Log.SetFormatter(&logrus.JSONFormatter{})
        Log.SetLevel(logrus.WarnLevel)
    } else {
        Log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
        Log.SetLevel(logrus.DebugLevel)
    }
}
