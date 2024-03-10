package db

import (
	"gorm.io/gorm"
	"rockbackup/backend/schedules"
	"rockbackup/backend/service"
)

func (db *DB) HasSource(id uint) bool {
	return false
}

// SaveService method
func (db *DB) SaveService(src *service.BackupSource, policy *service.Policy, schs []schedules.Schedule) error {

	err := db.g.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(src); result.Error != nil {
			return result.Error
		}

		if result := tx.Create(policy); result.Error != nil {
			return result.Error
		}

		if result := tx.Create(schs); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) AddSchedule(sch *schedules.Schedule) (uint, error) {
	var id uint

	return id, nil
}

func (db *DB) GetAllEnabledSchedules() ([]schedules.Schedule, error) {
	return []schedules.Schedule{}, nil
}
