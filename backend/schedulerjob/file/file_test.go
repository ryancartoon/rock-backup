package filejob

type MockDB struct{}

func (d *MockDB) SaveBackupResult(id uint, SnapID string, Size int64, FileNum int64) error {
	return nil
}
func (d *MockDB) SaveBackupError(id uint, err string) {}

// func TestAgentRunCmd(t *testing.T) {
// 	agent := &agentd.Agent{Host: "localhost", Port: 50001}
// 	src := service.BackupSource{SourcePath: "/home/ryan/codes/rock-backup/proto"}
// 	policy := service.Policy{BackupSource: src}
// 	repo := &repository.Repository{MountPoint: "/tmp/repo"}
// 	db := &MockDB{}
// 	ctx := context.Background()
// 	schedulerJob := schedulerjob.Job{}
// 	job := NewFileBackupSchedulerJob(schedulerJob)
// 	// err := job.Run(ctx, policy, "Full", repo, agent)
// 	err := job.Run(ctx, db, policy, repo, agent)

// 	assert.NoError(t, err)

// }
