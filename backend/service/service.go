package service

import (
	"errors"
	"rockbackup/backend/schedules"

	"gorm.io/datatypes"
)

const (
	BackupTypeFull = "full"
	BackupTypeIncr = "incremental"
)

var (
	ResourceAlreadyExistError = errors.New("error: already exist")
)

type PolicyView struct {
	ID             uint   `json:"id"`
	SourceID       uint   `json:"source_id"`
	SourceType     string `json:"source_type"`
	SourcePath     string `json:"source_path"`
	SourceHost     string `json:"source_host"`
	Retention      uint   `json:"retention"`
	BackupSourceID uint   `json:"backup_source_id"`
	Hostname       string `json:"hostname"`
	Status         string `json:"status"`
	RepositoryID   uint   `json:"repository_id"`
	// ScheduleDesc   string `json:"schedule_desc"`
	// todo
	// FullDay uint `json:"full_day"`
	// Repository     repository.Repository
	// Schedules      []schedules.Schedule
}

type PolicyRequest struct {
	Retention          uint
	BackupSourcePath   string
	Hostname           string
	RepositoryID       uint
	BackupSourceID     uint
	FullBackupSchedule string
	IncrBackupSchedule string
	ScheduleDesc       string
	StartTime          datatypes.Time
	BackupCycle        uint
}

type PolicyWithSource struct {
	Policy
	BackupSource
}

type BackupServiceI interface {
	OpenFile(PolicyRequest) error
	GetPolicies() ([]PolicyView, error)
	// OpenDB(src *BackupSource, policy *Policy, schs []schedules.Schedule) error Close(srcID uint) error }
}

type JobStarter interface {
	StartBackup(policyID uint, backupType string, opertor string) error
	StartRestore()
}

type DB interface {
	SaveService(*BackupSource, *Policy, []schedules.Schedule) error
	GetPolicies() ([]Policy, error)
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

func (s *BackupService) OpenFile(req PolicyRequest) error {

	var schs []schedules.Schedule

	sourceType := "file"

	src := &BackupSource{
		SourceType: sourceType,
		SourcePath: req.BackupSourcePath,
	}

	policy := &Policy{
		Retention:    req.Retention,
		Status:       ServiceStatusServing,
		RepositoryID: req.RepositoryID,
		Hostname:     req.Hostname,
	}

	full := schedules.Schedule{Cron: req.FullBackupSchedule, StartTime: req.StartTime}
	incr := schedules.Schedule{Cron: req.IncrBackupSchedule, StartTime: req.StartTime}
	schs = []schedules.Schedule{full, incr}

	// src.SourceName = src.SourcePath

	// if s.hasSource(src.SourceName) {
	// 	return ResourceAlreadyExistError
	// }

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

	if s.hasSource(src.SourceName) {
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

func (s *BackupService) GetPolicies() ([]PolicyView, error) {
	var pvs []PolicyView
	ps, err := s.db.GetPolicies()

	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		pvs = append(pvs, PolicyView{
			ID:         p.ID,
			SourceType: p.BackupSource.SourceType,
			SourcePath: p.BackupSource.SourcePath,
			SourceHost: p.BackupSource.SourceName,
			Hostname:   p.Hostname,
			Retention:  p.Retention,
			Status:     p.Status,
		})
	}

	return pvs, nil
}
