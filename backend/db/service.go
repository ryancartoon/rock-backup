package db

import (
	"rockbackup/backend/policy"
	"rockbackup/backend/schedules"

	"gorm.io/gorm"
)

func (db *DB) HasSource(id uint) bool {
	return false
}

// SaveService method
func (db *DB) SaveService(src *policy.BackupSource, policy *policy.Policy, schs []schedules.Schedule) error {

	err := db.g.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(src); result.Error != nil {
			return result.Error
		}

		policy.BackupSourceID = src.ID

		if result := tx.Create(policy); result.Error != nil {
			return result.Error
		}

		for i := range schs {
			schs[i].PolicyID = policy.ID
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

func (db *DB) GetPolicies() ([]policy.Policy, error) {
	var ps []policy.Policy
	if result := db.g.Table("policies").InnerJoins("BackupSource").Find(&ps); result.Error != nil {
		return nil, result.Error
	}

	return ps, nil

}

func (db *DB) AddSchedule(sch *schedules.Schedule) (uint, error) {
	var id uint

	return id, nil
}

func (db *DB) GetAllEnabledSchedules() ([]schedules.Schedule, error) {
	var schs []schedules.Schedule

	result := db.g.Table("schedules").Select(
		`schedules.id, schedules.policy_id, schedules.cron, schedules.start_time,
		schedules.duration, schedules.backup_type`,
	).Where(`is_enabled=?`, true).Scan(&schs)

	if result.Error != nil {
		return []schedules.Schedule{}, result.Error
	}

	return schs, nil
}
