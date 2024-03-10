package service

import (
	"errors"
	"rockbackup/backend/schedules"
)

const (
	BackupTypeFull = "full"
	BackupTypeIncr = "incremental"
)

var (
	ResourceAlreadyExistError = errors.New("error: already exist")
)

type BackupServiceI interface {
	OpenFile(src *BackupSource, policy *Policy, schs []schedules.Schedule) error
	// OpenDB(src *BackupSource, policy *Policy, schs []schedules.Schedule) error
	Close(srcID uint) error
}

type JobStarter interface {
	StartBackup(policyID uint, backupType string, opertor string) error
	StartRestore()
}

type DB interface {
	SaveService(*BackupSource, *Policy, []schedules.Schedule) error
	HasSource(ID uint) bool
}

func New(db DB, sched *schedules.TimeScheduler) *BackupService {
	return &BackupService{db: db, timeSched: sched}
}

type BackupService struct {
	timeSched *schedules.TimeScheduler
	db        DB
}

//	func New(sm ScheduleMan, sched Scheduler, db DB) {
//		return &BackupService{sm, sched, db}
//	}

func (s *BackupService) OpenFile(src *BackupSource, policy *Policy, schs []schedules.Schedule) error {

	src.Name = src.Host + "-" + src.DataPath

	if s.hasSource(src.Name) {
		return ResourceAlreadyExistError
	}

	// save source, policy, schedules to get ID
	if err := s.db.SaveService(src, policy, schs); err != nil {
		return err
	}

	if err := s.timeSched.AddSchedules(schs); err != nil {
		return err
	}

	return nil
}

func (s *BackupService) OpenDB(src BackupSource, policy *Policy, schs []schedules.Schedule) error {

	if s.hasSource(src.Name) {
		return ResourceAlreadyExistError
	}

	// save source, policy, schedules to get ID
	if err := s.db.SaveService(&src, policy, schs); err != nil {
		return err
	}

	if err := s.timeSched.AddSchedules(schs); err != nil {
		return err
	}

	return nil
}

func (s *BackupService) hasSource(name string) bool {
	return false
}

func (s *BackupService) Close(srcID uint) error {
	return nil
}
