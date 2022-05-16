package api_common

import (
	"os"
	"strings"

	"go.elastic.co/ecslogrus"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes settings for logrus logger
func InitLogger(loglevel string) {

	formatter := ecslogrus.Formatter{
		DisableHTMLEscape: true,
		DataKey:           "variables",
		CallerPrettyfier:  nil,
		PrettyPrint:       true,
	}
	logrus.SetFormatter(&formatter)
	logrus.SetReportCaller(true)
	switch strings.ToUpper(os.Getenv(loglevel)) {
	case "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
