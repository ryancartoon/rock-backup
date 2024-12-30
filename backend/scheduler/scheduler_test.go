package scheduler

import (
	"rockbackup/backend/policy"
	"rockbackup/backend/schedulerjob"
	// "testing"
	//
	// "github.com/spf13/viper"
)

type mdb struct{}
type mhandler struct{}

func (h *mhandler) Handle(job JobInSchedule) error      { return nil }
func (m *mdb) AddSchedulerJob(*schedulerjob.Job) error  { return nil }
func (m *mdb) GetPolicy(uint) (policy.Policy, error)    { return policy.Policy{}, nil }
func (m *mdb) GetOnGoingJobs() ([]JobInSchedule, error) { return nil, nil }
func (m *mdb) StartJob(id uint) error

// func TestSchedulerStart(t *testing.T) {
// 	db := &mdb{}
// 	handler := &mhandler{}
// 	conf := &viper.Viper{}
// 	schr := New(conf, db, handler)
// 	go schr.Start()
// 	schr.Stop()
// }
