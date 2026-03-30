package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"time"
)

var Logger *log.Logger

func Init(logPath string) *rotatelogs.RotateLogs {
	if logPath == "" {
		logPath = "./logs"
	}
	Logger = log.New()
	info, err := rotatelogs.New(
		logPath+"/info"+".%Y%m%d",
		rotatelogs.WithLinkName(logPath+"info"),
		rotatelogs.WithMaxAge(5*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(100*1024*1024),
	)
	warn, err := rotatelogs.New(
		logPath+"/warn"+".%Y%m%d",
		rotatelogs.WithLinkName(logPath+"warn"),
		rotatelogs.WithMaxAge(5*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(100*1024*1024),
	)
	errWriter, err := rotatelogs.New(
		logPath+"/error"+".%Y%m%d",
		rotatelogs.WithLinkName(logPath+"error"),
		rotatelogs.WithMaxAge(5*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(100*1024*1024),
	)
	fatal, err := rotatelogs.New(
		logPath+"/fatal"+".%Y%m%d",
		rotatelogs.WithLinkName(logPath+"fatal"),
		rotatelogs.WithMaxAge(5*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(100*1024*1024),
	)
	panicWriter, err := rotatelogs.New(
		logPath+"panic"+".%Y%m%d",
		rotatelogs.WithLinkName(logPath+"panic"),
		rotatelogs.WithMaxAge(5*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(100*1024*1024),
	)
	if err != nil {
		log.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}
	Logger.SetFormatter(&log.TextFormatter{})
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.InfoLevel:  info,
		log.WarnLevel:  warn,
		log.ErrorLevel: errWriter,
		log.FatalLevel: fatal,
		log.PanicLevel: panicWriter,
	}, &log.TextFormatter{})
	Logger.AddHook(lfHook)
	return info
}
