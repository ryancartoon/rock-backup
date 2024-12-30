package logbackup

import (
	"rockbackup/backend/dblogmeta"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) AddFileMeta(fileMeta dblogmeta.FileMeta) error {
	args := m.Called(fileMeta)
	return args.Error(0)
}

func (m *MockDB) GetLastUpdateTime(instanceName string) (time.Time, error) {
	args := m.Called(instanceName)
	return args.Get(0).(time.Time), args.Error(1)
}

// func (m *MockDB) GetAllPolicy() ([]policy.Policy, error) {
// 	args := m.Called()
// 	return args.Get(0).([]policy.Policy), args.Error(1)
// }

// type MockAgent struct {
// 	mock.Mock
// }

// func (m *MockAgent) Scan(ctx context.Context, path string, startTime time.Time) ([]dblogmeta.FileMeta, error) {
// 	args := m.Called(ctx, path, startTime)
// 	return args.Get(0).([]dblogmeta.FileMeta), args.Error(1)
// }

// type MockAgentd struct {
// 	mock.Mock
// }

// func (m *MockAgentd) GetAgent(host string) (agentd.Agent, error) {
// 	args := m.Called(host)
// 	return args.Get(0).(agentd.Agent), args.Error(1)
// }

// func TestNewDBLogWatch(t *testing.T) {
// 	mockKeeperr := &keeperr.KeepErr{}
// 	mockAgentd := &MockAgentd{}
// 	logWatcher := NewDBLogWatch(mockKeeperr, mockAgentd)

// 	assert.NotNil(t, logWatcher)
// 	assert.Equal(t, mockAgentd, logWatcher.agentd)
// 	assert.Equal(t, mockKeeperr, logWatcher.keeperr)
// }

// func TestLogWatcher_LoopPolices(t *testing.T) {
// 	mockDB := new(MockDB)
// 	mockAgentd := new(MockAgentd)
// 	mockAgent := new(MockAgent)
// 	mockKeeperr := &keeperr.KeepErr{}
// 	logWatcher := &LogWatcher{
// 		db:            mockDB,
// 		agentd:        mockAgentd,
// 		newFileMetaCh: make(chan dblogmeta.FileMeta, 1),
// 		keeperr:       mockKeeperr,
// 	}

// 	policies := []policy.Policy{
// 		{
// 			LogHostname: "host1",
// 			BackupSource: policy.BackupSource{
// 				SourceName: "source1",
// 			},
// 			Repository: policy.Repository{
// 				MountPoint: "/mnt",
// 			},
// 		},
// 	}

// 	mockDB.On("GetAllPolicy").Return(policies, nil)
// 	mockDB.On("GetLastUpdateTime", "source1").Return(time.Now(), nil)
// 	mockAgentd.On("GetAgent", "host1").Return(mockAgent, nil)
// 	mockAgent.On("Scan", mock.Anything, "/mnt/logs/source1", mock.Anything).Return([]dblogmeta.FileMeta{}, nil)

// 	logWatcher.LoopPolices()

// 	mockDB.AssertExpectations(t)
// 	mockAgentd.AssertExpectations(t)
// 	mockAgent.AssertExpectations(t)
// }

// func TestLogWatcher_Watch(t *testing.T) {
// 	mockDB := new(MockDB)
// 	mockAgentd := new(MockAgentd)
// 	mockKeeperr := &keeperr.KeepErr{}
// 	logWatcher := &LogWatcher{
// 		db:            mockDB,
// 		agentd:        mockAgentd,
// 		newFileMetaCh: make(chan dblogmeta.FileMeta, 1),
// 		stoppingCh:    make(chan struct{}),
// 		keeperr:       mockKeeperr,
// 	}

// 	go func() {
// 		time.Sleep(1 * time.Second)
// 		logWatcher.stoppingCh <- struct{}{}
// 	}()

// 	err := logWatcher.Watch()
// 	assert.NoError(t, err)
// }

// func TestLogWatcher_Stop(t *testing.T) {
// 	mockDB := new(MockDB)
// 	mockAgentd := new(MockAgentd)
// 	mockKeeperr := &keeperr.KeepErr{}
// 	logWatcher := &LogWatcher{
// 		db:            mockDB,
// 		agentd:        mockAgentd,
// 		newFileMetaCh: make(chan dblogmeta.FileMeta, 1),
// 		stoppingCh:    make(chan struct{}),
// 		keeperr:       mockKeeperr,
// 	}

// 	go func() {
// 		time.Sleep(1 * time.Second)
// 		logWatcher.Stop()
// 	}()

// 	select {
// 	case <-logWatcher.stoppingCh:
// 		assert.True(t, true)
// 	case <-time.After(2 * time.Second):
// 		assert.Fail(t, "Stop() did not send to stoppingCh in time")
// 	}
// }
