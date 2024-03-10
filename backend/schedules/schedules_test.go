package schedules

import (
	gocron "github.com/robfig/cron/v3"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// https://dev.to/salesforceeng/mocks-in-go-tests-with-testify-mock-6pd

type mockdb struct {
	mock.Mock
}

func (d *mockdb) AddSchedule(*Schedule) (uint, error) {
	return 0, nil
}

func (d *mockdb) GetAllEnabledSchedules() ([]Schedule, error) {
	args := d.Called()
	return args.Get(0).([]Schedule), args.Error(1)
}

type starter struct{}

func (j *starter) TimerStartBackup(policyID uint, backupType string, operator string) error {
	return nil
}

func newMock(db DB) *TimeScheduler {
	jb := &starter{}
	config := &viper.Viper{}

	cron := gocron.New()
	sm := New(config, db, jb, cron)

	return sm
}

func TestScheduleManagerStart(t *testing.T) {
	db := &mockdb{}
	sm := newMock(db)
	go sm.Start()
	sm.Stop()
}

func TestScheduleManagerInit(t *testing.T) {
	db := &mockdb{}
	schs := []Schedule{{ID: 1, PolicyID: 1, Cron: "* * * * *"}, {ID: 2, PolicyID: 2, Cron: "* * * * *"}}

	db.On("GetAllEnabledSchedules").Return(schs, nil)

	sm := newMock(db)
	sm.init()

	assert.Equal(t, 2, len(sm.cron.Entries()))
	assert.Equal(t, true, sm.cronInitialized)
}
