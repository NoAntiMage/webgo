package logx

import (
	"fmt"
	"goweb/config"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Trace(args ...any)
	Tracef(format string, args ...any)
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Print(args ...any)
	Printf(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Panic(args ...any)
	Panicf(format string, args ...any)
}

var _ Logger = (*logrus.Logger)(nil)

var Loggerx *logrus.Logger

//DEBUG
type CostomFormatter struct{}

func (f *CostomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	message := fmt.Sprintf("%s %s :\"%s\"\n",
		entry.Time.Format(time.DateTime),
		entry.Level.String(),
		entry.Message,
		//TODO Runtime
	)

	return []byte(message), nil
}

func LoggerInit() {
	Loggerx = logrus.New()
	logConf := config.GetGlobalConfig().LogConf

	writers := make([]io.Writer, 0)
	writers = append(writers, os.Stdout)

	switch logConf.Output {
	case "file":
		_, err := os.Stat(logConf.Path)
		if err != nil {
			err = os.MkdirAll(logConf.Path, 0755)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("log path not exist.\nMake dir at: " + logConf.Path)
		}
		logFile := fmt.Sprintf("%v/app.log", logConf.Path)
		fileWriter, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			fmt.Println(err)
		} else {
			writers = append(writers, fileWriter)
		}
	default:
	}
	Loggerx.SetOutput(io.MultiWriter(writers...))

	switch logConf.Level {
	case "trace":
		Loggerx.SetLevel(logrus.TraceLevel)
	case "debug":
		Loggerx.SetLevel(logrus.DebugLevel)
	case "info":
		Loggerx.SetLevel(logrus.InfoLevel)
	default:
		Loggerx.SetLevel(logrus.DebugLevel)
	}

	Loggerx.SetFormatter(&CostomFormatter{})
}

func GetLogger() Logger {
	return Loggerx
}
