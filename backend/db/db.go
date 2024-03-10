package db

import (
	"log"
	"os"
	"path/filepath"
	"time"

	// "rockbackup/backend"
	// "github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"rockbackup/backend/schedules"
	"rockbackup/backend/service"
)

type DB struct {
	g *gorm.DB
}

func (d *DB) AutoMigrate() error {
	if err := d.g.AutoMigrate(
		&service.BackupSource{},
		&service.Policy{},
		&schedules.Schedule{},
	); err != nil {
		return err
	}

	return nil
}

func InitTest() *DB {
	appHome := "."

	logPath := filepath.Join(appHome, "logs", "testing.log")
	logFh, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, 0600)
	logger := logger.New(
		log.New(logFh, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(appHome, "test.db")), &gorm.Config{Logger: logger})

	if err != nil {
		panic("init db error")
	}

	return &DB{gdb}
}

// func InitDB(appHome string, config *viper.Viper, logPath string) *DB {
// 	var db *gorm.DB
// 	logger := initDBLogger(backend.AppHome)
// 	dialect := config.GetString("database.dialect")
// 	port := config.GetString("database.port")
// 	dbname := config.GetString("database.dbname")
// 	host := configGetString("database.host")
// 	dsn := fmt.Sprinft(config.GetString("database.dsn"), user, pass, host, port, dbname)
//
// 	if dialect == "postgres" {
// 		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	}
//
// 	return db
// }
//
// func initDBLogger(logBasePath string) logger.Interface {
// 	logName = "gorm.log"
// 	logPath := filepath.Join(logBasePath, logName)
// 	logFh, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, 0600)
// 	logger := logger.New(
// 		log.New(logFh, "\r\n", log.LstdFlags), // io writer
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
