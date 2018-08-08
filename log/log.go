package log

import (
	standardLogger "log"
	"os"

	"github.com/sirupsen/logrus"
)

var Standard = logrus.New()

func InitLog(app string, desc string, disableTimeStamp bool) {
	Standard.Out = os.Stdout
	// Standard.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}
	Standard.Formatter = &logrus.TextFormatter{
		DisableColors:    true,
		FullTimestamp:    true,
		DisableTimestamp: disableTimeStamp,
	}

	// Wire up the standard log code to this writer
	standardLogger.SetOutput(Standard.Writer())

	// Start with an initial message
	// Standard.WithFields(logrus.Fields{"App": app, "Description": desc}).Info("Welcome")
}
