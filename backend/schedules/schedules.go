package schedules

import (
	"sync"
	"time"

	"gorm.io/datatypes"

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
	ID          uint           `json:"id"          gorm:"column:id;primaryKey;autoIncrement"`
	PolicyID    uint           `json:"policy_id"   gorm:"column:policy_id"`
	Cron        string         `json:"cron"        gorm:"column:cron"`
	StartTime   datatypes.Time `json:"start_time"  gorm:"column:start_time"`
	Duration    time.Duration  `json:"duration"    gorm:"column:duration"`
	BackupType  string         `json:"backup_type" gorm:"column:backup_type"`
	Description string         `json:"description" gorm:"column:description"`
	IsEnabled   bool           `json:"is_enabled"  gorm:"column:is_enabled"`
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
		stopping:   make(chan struct{}),
	}
}

// type BackupSchedule struct {
// 	ID            uint
// 	PolicyID      uint
// 	BackupType    string
// 	Schedule      string
// 	Duration      int
// 	NextStartTIme time.Time
// 	IsEnabled     bool
// }

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
		entryID, err := s.cron.AddJob(sch.Cron, gocronJob{s.starter, sch.PolicyID, sch.ID, sch.BackupType})

		if err != nil {
			return err
		}
		s.ste.update(sch.ID, entryID)
	}

	return nil
}

func (s *TimeScheduler) init() {
	logger.Info("init schedules")
	schs, _ := s.db.GetAllEnabledSchedules()

	for _, sch := range schs {
		if _, err := s.specParser.Parse(sch.Cron); err != nil {
			logger.Errorf("cron of schedule [%d] format [%s] is incorrect", sch.ID, sch.Cron)
		}

		entryID, _ := s.cron.AddJob(
			sch.Cron,
			gocronJob{starter: s.starter, policyID: sch.PolicyID, scheduleID: sch.ID, backupType: sch.BackupType},
		)

		s.ste.update(sch.ID, entryID)
	}

	s.cronInitialized = true
	logger.Info("schedules is initialized")
}

func (s *TimeScheduler) Start() {
	logger.Info("Start the schedules")
	var once sync.Once

	once.Do(s.init)

	// start cron
	s.cron.Start()
	s.started = true

runningLoop:
	for {
		select {
		case <-time.After(time.Hour):
			logger.Info("schedules heart is beating")
		case <-s.stopping:
			logger.Info("stopping schedules")
			s.cron.Stop()
			break runningLoop
		}
	}

	logger.Info("schedules manager is stopped")
}

func (s *TimeScheduler) Stop() {
	logger.Info("stopping time scheduler")
	if s.started {
		s.stopping <- struct{}{}
		s.started = false
	}
}

type gocronJob struct {
	starter    TimerStarter
	policyID   uint
	scheduleID uint
	backupType string
}

func (g gocronJob) Run() {
	operator := "backup scheduler"
	err := g.starter.TimerStartBackup(g.policyID, g.backupType, operator)

	if err != nil {
		logger.Errorf("schedule backup job for schedule:[%d], policy[%d] failed", g.scheduleID, g.policyID)
	}
}
