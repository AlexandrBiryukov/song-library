// logger/logger.go
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	Log.Out = os.Stdout
	Log.SetLevel(logrus.DebugLevel)
}
