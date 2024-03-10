package schedules

import (
	"gorm.io/datatypes"
	"sync"
	"time"

	gocron "github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type TimeScheduler struct {
	db              DB
	cron            *gocron.Cron
	specParser      gocron.Parser
	ste             *ScheduleToEntry
	starter         TimerStarter
	config          *viper.Viper
	cronInitialized bool
	started         bool
	stopping        chan struct{}
}

type Schedule struct {
	gorm.Model
	ID          uint
	PolicyID    uint
	Cron        string
	StartTime   datatypes.Time
	Duration    time.Duration
	BackupType  string
	Description string
}

type DB interface {
	// AddSchedule(policyID uint, backupType string, cronStr string, duration uint) (uint, error)
	AddSchedule(*Schedule) (uint, error)
	GetAllEnabledSchedules() ([]Schedule, error)
}

type TimerStarter interface {
	TimerStartBackup(policyID uint, backupType string, operator string) error
}

func New(config *viper.Viper, db DB, starter TimerStarter, cron *gocron.Cron) *TimeScheduler {
	return &TimeScheduler{
		db:         db,
		cron:       cron,
		specParser: gocron.NewParser(gocron.Minute | gocron.Hour | gocron.Dom | gocron.Month | gocron.Dow),
		starter:    starter,
		config:     config,
		ste:        &ScheduleToEntry{m: make(map[uint]gocron.EntryID)},
	}
}

type BackupSchedule struct {
	ID            uint
	PolicyID      uint
	BackupType    string
	Schedule      string
	Duration      int
	NextStartTIme time.Time
	IsEnabled     bool
}

type ScheduleToEntry struct {
	mu sync.Mutex
	m  map[uint]gocron.EntryID
}

func (e *ScheduleToEntry) update(schedID uint, entryID gocron.EntryID) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.m[schedID] = entryID
}

// func (s *TimeScheduler) AddSchedule(policyID uint, backupType string, cronStr string, duration uint) error {
func (s *TimeScheduler) AddSchedules(scheds []Schedule) error {

	for i := range scheds {
		sch := scheds[i]
		entryID, err := s.cron.AddJob(sch.Cron, queueJob{s, sch.PolicyID, sch.ID, sch.BackupType})

		if err != nil {
			return err
		}
		s.ste.update(sch.ID, entryID)
	}

	return nil
}

func (s *TimeScheduler) init() {
	schs, _ := s.db.GetAllEnabledSchedules()

	for _, sch := range schs {
		if _, err := s.specParser.Parse(sch.Cron); err != nil {
			logger.Errorf("cron of schedule [%d] format [%s] is incorrect", sch.ID, sch.Cron)
		}

		entryID, _ := s.cron.AddJob(
			sch.Cron,
			queueJob{s: s, policyID: sch.PolicyID, scheduleID: sch.ID, backupType: sch.BackupType},
		)

		s.ste.update(sch.ID, entryID)
	}

	s.cronInitialized = true
}

func (s *TimeScheduler) Start() {
	var once sync.Once

	once.Do(s.init)

	// start cron
	s.cron.Start()

runningLoop:
	for {
		select {
		case <-time.After(time.Hour):
			logger.Debug("schedule manmger heart is beating")
		case <-s.stopping:
			logger.Info("stopping schedules maanger")
			s.cron.Stop()
			break runningLoop
		}
	}
	logger.Info("schedules manager is stopped")
}

func (s *TimeScheduler) Stop() {
	if s.started {
		s.stopping <- struct{}{}
		s.started = false
	}
}

type queueJob struct {
	s          *TimeScheduler
	policyID   uint
	scheduleID uint
	backupType string
}

func (q queueJob) Run() {
	operator := "backup scheduler"
	err := q.s.starter.TimerStartBackup(q.policyID, q.backupType, operator)

	if err != nil {
		logger.Errorf("schedule backup job for schedule:[%d], policy[%d] failed", q.scheduleID, q.policyID)
	}
}
