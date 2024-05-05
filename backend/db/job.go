package db

import "rockbackup/backend/schedulerjob"

func (db *DB) SaveBackupError(id uint, errMessage string) {
	db.g.Model(&schedulerjob.Job{}).Where("id = ?", id).Updates(
		map[string]interface{}{"error_message": errMessage},
	)
}

func (db *DB) SaveBackupResult(id uint, snapID string, size int64, fileNum int64) error {
	return nil
}
