package db

import (
	"rockbackup/backend/schedules"
	"rockbackup/backend/service"

	"gorm.io/gorm"
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

// func (db *DB) GetPolicies() ([]service.PolicyView, error) {
func (db *DB) GetPolicies() ([]service.PolicyView, error) {

	var ps []service.PolicyView

	result := db.g.Table("policies").Select(
		`
		backup_sources.id, 
		backup_sources.source_name,
		backup_sources.source_type,
		policies.id, 
		policies.retention,
		policies.status,
		policies.schedule_desc,
		hosts.name as hostname
		`,
	).Joins(
		"left join backup_sources on policies.backup_source_id = backup_sources.id",
	).Joins(
		"left join hosts on policies.hostname = hosts.name",
	).Scan(&ps)

	if result.Error != nil {
		return []service.PolicyView{}, result.Error
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
