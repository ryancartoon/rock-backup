package xtrabackup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAgent struct {
	mock.Mock
}

func (m *MockAgent) RunCmd(ctx context.Context, name string, args []string, envs []string) (int, []byte, []byte, error) {
	call := m.Called(ctx, name, args, envs)
	return call.Int(0), call.Get(1).([]byte), call.Get(2).([]byte), call.Error(1)
}

func TestXtraBackupBackup(t *testing.T) {
	version := "8.0"
	BinaryPath := "/usr/bin/xtrabackup"
	x := NewXtrabackup(version, BinaryPath)

	agent := &MockAgent{}
	ctx := context.Background()

	backupType := "full"
	instance := "mysql-instance"
	repo := "backup-repo"
	targetPath := "/backup/path"
	err := x.Backup(ctx, agent, backupType, instance, repo, targetPath)
	assert.NoError(t, err)
}
