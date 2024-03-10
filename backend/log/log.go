package log

import (
	// "time"

	"github.com/sirupsen/logrus"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	// "gorm.io/gorm/logger"
)

func New(logPath string) *logrus.Logger {
	maxSize := 100
	maxAge := 30
	maxBackupss := 5

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
	log.SetLevel(logrus.DebugLevel)

	return log
}

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
