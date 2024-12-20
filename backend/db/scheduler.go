package db

import (
	"rockbackup/backend/backupset"
	"rockbackup/backend/scheduler"
	"rockbackup/backend/schedulerjob"
)

func (db *DB) GetOnGoingJobs() (jobs []scheduler.JobInSchedule, err error) {
	var sjobs []schedulerjob.Job

	result := db.g.Model(&sjobs).Where(
		"status = ? OR status =? OR status = ?",
		schedulerjob.SchedulerJobStatusQueued,
		schedulerjob.SchedulerJobStatusRunning,
		schedulerjob.SchedulerJobStatusCreated,
	).Find(&sjobs)

	if result.Error != nil {
		return nil, err
	}

	for _, job := range sjobs {
		jobs = append(jobs, scheduler.JobInSchedule{Job: job})
	}

	return jobs, nil
}

func (db *DB) GetJobsInschedule() (jobs []scheduler.JobInSchedule, err error) {
	var sjobs []schedulerjob.Job

	result := db.g.Model(&sjobs).Where(
		"in_schedule = ?", true,
	).Find(&sjobs)

	if result.Error != nil {
		return nil, err
	}

	for _, job := range sjobs {
		jobs = append(jobs, scheduler.JobInSchedule{Job: job})
	}

	return jobs, nil
}

func (db *DB) StartJob(id uint) error {
	if result := db.g.Model(&schedulerjob.Job{}).Where("id = ?", id).Updates(
		map[string]interface{}{"status": schedulerjob.SchedulerJobStatusRunning},
	); result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) GetBackupset(id uint) (bset backupset.Backupset, err error) {
	return bset, nil
}
