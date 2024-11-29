package db

import (
	"rockbackup/backend/backupset"
)

func (db *DB) AddBackupset(repoID uint, jobID uint, backup_type string) (*backupset.Backupset, error) {
	bset := &backupset.Backupset{}
	var err error

	bset.RepositoryID = repoID
	bset.JobID = jobID
	bset.BackupType = backup_type

	if err = db.g.Create(bset).Error; err != nil {
		return nil, err
	}

	return bset, nil
}
