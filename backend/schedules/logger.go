package schedules

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var logger = initLog()

func initLog() *logrus.Logger {
	appHome := "."
	logPath := filepath.Join(appHome, "logs", "schedules.log")
	logFh, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, 0600)

	log := logrus.New()
	log.Out = logFh

	return log
}
