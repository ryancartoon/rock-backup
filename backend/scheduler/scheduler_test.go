package scheduler

import (
	"github.com/spf13/viper"
	"testing"
)

type mdb struct{}
type mhandler struct{}

func (h *mhandler) StartBackup(policyID uint) error     { return nil }
func (h *mhandler) StartRestore(backupsetID uint) error { return nil }

func TestSchedulerStart(t *testing.T) {
	db := &mdb{}
	handler := &mhandler{}
	conf := &viper.Viper{}
	schr := New(conf, db, handler)
	go schr.Start()
	schr.Stop()
}
