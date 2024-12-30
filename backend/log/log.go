package log

import (
	// "time"
	"fmt"
	"os"

	"path/filepath"

	"github.com/sirupsen/logrus"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Log *logrus.Logger
}

type Fields map[string]interface{}

func New(logName string) *Logger {
	maxSize := 100
	maxAge := 30
	maxBackupss := 5

	appHome, _ := filepath.Abs(".")

	logPath := filepath.Join(appHome, "logs", logName)

	log := logrus.New()
	output := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackupss,
		Compress:   true,
		LocalTime:  true,
	}

	log.SetOutput(output)
	// log.SetLevel(logrus.DebugLevel)

	return &Logger{log}
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Log.Debug(args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Log.Info(args...)
}

func (logger *Logger) WithFields(fields logrus.Fields) *Logger {
	logger.Log.WithFields(fields)
	return logger
}

//	func (logger *Logger) Print(args ...interface{}) {
//		entry := logger.newEntry()
//		entry.Print(args...)
//		logger.releaseEntry(entry)
//	}
//
//	func (logger *Logger) Warn(args ...interface{}) {
//		logger.Log(WarnLevel, args...)
//	}
//
//	func (logger *Logger) Warning(args ...interface{}) {
//		logger.Warn(args...)
//	}
func (logger *Logger) Error(args ...interface{}) {
	logger.Log.Error(args...)
}

//	func (logger *Logger) Fatal(args ...interface{}) {
//		logger.Log(FatalLevel, args...)
//		logger.Exit(1)
//	}
//
//	func (logger *Logger) Panic(args ...interface{}) {
//		logger.Log(PanicLevel, args...)
//	}
//
//	func (logger *Logger) Tracef(format string, args ...interface{}) {
//		logger.Logf(TraceLevel, format, args...)
//	}
//
//	func (logger *Logger) Debugf(format string, args ...interface{}) {
//		logger.Logf(DebugLevel, format, args...)
//	}
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Log.Logf(logrus.InfoLevel, format, args...)
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

//	func (logger *Logger) Warnf(format string, args ...interface{}) {
//		logger.Logf(WarnLevel, format, args...)
//	}
//
//	func (logger *Logger) Warningf(format string, args ...interface{}) {
//		logger.Warnf(format, args...)
//	}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.Log.Logf(logrus.ErrorLevel, format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.Log.Logf(logrus.FatalLevel, format, args...)
	os.Exit(1)
}

//
// func (logger *Logger) Panicf(format string, args ...interface{}) {
// 	logger.Logf(PanicLevel, format, args...)
// }
//
// func NewGormLoggerForTest() *logrus.Logger {
// 	logger := logrus.New(
// 		logger.Config{
// 			SlowThreshold:             time.Second,   // Slow SQL threshold
// 			LogLevel:                  logger.Silent, // Log level
// 			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
// 			ParameterizedQueries:      true,          // Don't include params in the SQL log
// 			Colorful:                  false,         // Disable color
// 		},
// 	)
//
// 	return logger
// }
