package schedules

import "rockbackup/backend/log"

// var logger = initLog()

// func initLog() *logrus.Logger {
// 	appHome := "."
// 	logPath := filepath.Join(appHome, "logs", "schedules.log")
// 	logFh, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, 0600)

// 	log := logrus.New()
// 	log.Out = logFh

// 	return log
// }

var logger = log.New("schedules.log")
