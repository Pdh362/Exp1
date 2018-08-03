package log

import (
	standardLogger "log"
	"os"

	"github.com/sirupsen/logrus"
)

var Standard = logrus.New()

func InitLog(logfile string, app string, desc string) {
	Standard.Out = os.Stdout
	Standard.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}

	// logrus.SetFormatter(&logrus.JSONFormatter{})

	// If logfile specified, create the file to log to,
	// and inform the logger. Switch to text format as well.
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			Standard.Out = file
			Standard.Formatter = &logrus.TextFormatter{}
		}
	}

	// Wire up the standard log code to this writer
	standardLogger.SetOutput(Standard.Writer())

	// End with an initial message
	Standard.WithFields(logrus.Fields{"App": app, "Description": desc}).Info("Welcome")
}
