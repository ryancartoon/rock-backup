package db

import "gorm.io/gorm"
import "rockbackup/backend/schedulerjob"
import "rockbackup/backend/backupset"

func (db *DB) SaveBackupError(id uint, errMessage string) {
	db.g.Model(&schedulerjob.Job{}).Where("id = ?", id).Updates(
		map[string]interface{}{"error_message": errMessage, "status": "failed"},
	)
}

func (db *DB) SaveBackupResult(id uint, bsetID uint, snapID string, size int64, fileNum int64) error {

	var err error

	err = db.g.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&schedulerjob.Job{}).Where("id = ?", id).Updates(
			map[string]interface{}{
				"status": "completed",
			},
		).Error; err != nil {
			return err
		}

		return nil

	})

	err = db.g.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&backupset.Backupset{}).Where("id = ?", bsetID).Updates(
			map[string]interface{}{
				"external_backupset_id": snapID,
				"size":                  size,
				"file_num":              fileNum,
			},
		).Error; err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil

}
