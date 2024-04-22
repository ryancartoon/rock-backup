package schedulerjob

import (
	"rockbackup/backend/service"
)

type BackupJobFile struct {
	// logger
	policy service.Policy
}

func (j *BackupJobFile) Run() {

}

func LoadBackupJobFile(id uint) (*BackupJobFile, error) {
	return nil
}
