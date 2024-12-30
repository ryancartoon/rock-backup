package db

import (
	"github.com/stretchr/testify/assert"
	"rockbackup/backend/schedules"
	"rockbackup/backend/service"
	"testing"
)

func initServiceDB(t *testing.T) *DB {
	db := InitTest()
	if err := db.g.AutoMigrate(&service.BackupSource{}, &service.Policy{}, &schedules.Schedule{}); err != nil {
		t.Fatal(err)
	}

	return db
}

func TestSaveServiceDB(t *testing.T) {

	db := initServiceDB(t)

	src := service.BackupSource{}
	plc := service.Policy{}
	schs := []schedules.Schedule{}
	err := db.SaveService(&src, &plc, schs)

	if err != nil {
		assert.NoError(t, err)
	}
}
