package db

import (
	"log"
	"os"
	"path/filepath"

	"fmt"
	"time"

	// "rockbackup/backend"
	// "github.com/spf13/viper"
	"rockbackup/backend/host"
	"rockbackup/backend/schedulerjob"
	"rockbackup/backend/schedules"
	"rockbackup/backend/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	g *gorm.DB
}

func (d *DB) AutoMigrate() error {
	if err := d.g.AutoMigrate(
		&service.BackupSource{},
		&service.Policy{},
		&schedules.Schedule{},
		&host.Host{},
		&schedulerjob.Job{},
	); err != nil {
		return err
	}

	return nil
}

func InitTest() *DB {
	appHome := "."

	now := time.Now().Unix()
	logPath := filepath.Join(appHome, "logs", fmt.Sprintf("%s-%d.%s", "testing", now, "log"))
	logFh, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, 0600)
	l := logger.New(
		log.New(logFh, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
			// SlowThreshold:             time.Second,   // Slow SQL threshold
			// IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true,          // Don't include params in the SQL log
			// Colorful:                  true,          // Disable color
		},
	)
	l.LogMode(logger.Info)
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(appHome, "test.db")), &gorm.Config{Logger: l})

	if err != nil {
		panic("init db error")
	}

	return &DB{gdb}
}

func (db *DB) GetPolicy(id uint) (service.Policy, error) {
	var p service.Policy

	result := db.g.Joins("BackupSource").First(&p, id)

	if result.Error != nil {
		return service.Policy{}, result.Error
	}

	return p, nil
}

func (db *DB) AddSchedulerJob(job *schedulerjob.Job) error {
	result := db.g.Create(job)
	return result.Error
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
