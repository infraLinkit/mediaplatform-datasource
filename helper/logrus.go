package helper

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type MyFormatter struct{}

var levelList = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

func (mf *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	level := levelList[int(entry.Level)]
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1]
	b.WriteString(fmt.Sprintf("%s - %s - [line:%d] - %s - %s\n",
		entry.Time.Format("2006-01-02 15:04:05,678"), fileName,
		entry.Caller.Line, level, entry.Message))
	return b.Bytes(), nil
}

func MakeLogger(filename string, display bool, level string) *logrus.Logger {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		panic(err.Error())
	}
	logger := logrus.New()
	if display {
		logger.SetOutput(io.MultiWriter(os.Stdout, f))
	} else {
		logger.SetOutput(io.MultiWriter(f))
	}
	logger.SetReportCaller(true)
	//logger.SetFormatter(&MyFormatter{})
	logger.SetFormatter(&logrus.JSONFormatter{})

	level = strings.ToUpper(level)
	if level == "INFO" {
		logger.Level = logrus.InfoLevel
	} else if level == "DEBUG" {
		logger.Level = logrus.DebugLevel
	} else if level == "WARN" {
		logger.Level = logrus.WarnLevel
	} else if level == "ERROR" {
		logger.Level = logrus.ErrorLevel
	} else if level == "FATAL" {
		logger.Level = logrus.FatalLevel
	} else if level == "PANIC" {
		logger.Level = logrus.PanicLevel
	} else if level == "TRACE" {
		logger.Level = logrus.TraceLevel
	}

	/* rotatelog.NewHook("./"+filename+".%Y%m%d",
	rotatelog.WithMaxAge(24*time.Hour),
	rotatelog.WithRotationTime(time.Hour),
	rotatelog.WithClock(rotatelog.Local)) */

	/* if err != nil {
		log.Hooks.Add(hook)
	} */

	return logger
}
