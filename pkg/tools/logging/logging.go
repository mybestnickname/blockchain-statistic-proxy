package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/adjust/rmq/v4"
	"github.com/sirupsen/logrus"
)

// InitLogger init the logrus.Logger for defined log level.
func InitLogger(level logrus.Level) *logrus.Logger {
	return &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.JSONFormatter{
			CallerPrettyfier: caller,
			TimestampFormat:  "2006-01-02 15:04:05",
		},
		ReportCaller: true,
		Level:        level,
	}
}

type Logger struct {
	Logger *logrus.Logger
}

// InitLoggerTwo
func InitLoggerNew(level logrus.Level, out io.Writer) *Logger {

	return &Logger{
		Logger: &logrus.Logger{
			Out: out,
			Formatter: &logrus.JSONFormatter{
				CallerPrettyfier: caller,
				TimestampFormat:  "2006-01-02 15:04:05",
			},
			ReportCaller: true,
			Level:        level,
		},
	}

}

// Info
func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

// Infof
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Infoln
func (l *Logger) Infoln(args ...interface{}) {
	l.Logger.Infoln(args...)
}

// Error
func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

// Errorf
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Errorln
func (l *Logger) Errorln(args ...interface{}) {
	l.Logger.Errorln(args...)
}

// Debug
func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

// Debugf
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

// Debugln
func (l *Logger) Debugln(args ...interface{}) {
	l.Logger.Debugln(args...)
}

// Fatal
func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

// caller returns string presentation of log caller which is formatted as.
func caller(f *runtime.Frame) (function string, file string) {
	_, filename := path.Split(f.File)
	return "", fmt.Sprintf("%s:%d", filename, f.Line)
}

func LogErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *rmq.HeartbeatError:
			if err.Count == rmq.HeartbeatErrorLimit {
				log.Print("heartbeat error (limit): ", err)
			} else {
				log.Print("heartbeat error: ", err)
			}
		case *rmq.ConsumeError:
			log.Print("consume error: ", err)
		case *rmq.DeliveryError:
			log.Print("delivery error: ", err.Delivery, err)
		default:
			log.Print("other error: ", err)
		}
	}
}
