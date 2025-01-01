package logbackup

import (
	"context"
	"fmt"
	"path/filepath"
	"rockbackup/backend/agentd"
	"rockbackup/backend/dblogmeta"
	"rockbackup/backend/keeperr"
	"rockbackup/backend/log"
	"rockbackup/backend/policy"
	"time"
)

var logger *log.Logger

type DB interface {
	AddFileMeta(FileMeta dblogmeta.FileMeta) error
	GetLastUpdateTime(instanceName string) (time.Time, error)
	GetAllPolicy() ([]policy.Policy, error)
}

type Agentd interface {
	GetAgent(host string) (agentd.Agent, error)
}

func NewDBLogWatch(keeperr *keeperr.KeepErr, agentd Agentd) *LogWatcher {
	return &LogWatcher{
		agentd:        agentd,
		newFileMetaCh: make(chan dblogmeta.FileMeta),
		stoppingCh:    make(chan struct{}),
		keeperr:       keeperr,
	}
}

type LogWatcher struct {
	newFileMetaCh chan dblogmeta.FileMeta
	db            DB
	stoppingCh    chan struct{}
	agentd        Agentd
	keeperr       *keeperr.KeepErr
}

func genLogPath(mp string, instanceName string) string {
	return filepath.Join(mp, "logs", instanceName)
}

func (w *LogWatcher) LoopPolices(ctx context.Context) {
	polices, err := w.db.GetAllPolicy()
	if err != nil {
		logger.Errorf("get all policy error: %v", err)
		return
	}

	for _, p := range polices {
		metas, err := w.scanPolicy(ctx, p)

		if err != nil {
			logger.Errorf("scan error: %v", err)
			continue
		}

		for _, m := range metas {
			w.newFileMetaCh <- m
		}
	}
}

func (w *LogWatcher) scanPolicy(ctx context.Context, p policy.Policy) (metas []dblogmeta.FileMeta, err error) {
	var hostname string

	if p.LogHostname != "" {
		hostname = p.LogHostname
	} else if p.Hostname != "" {
		hostname = p.Hostname
	} else {
		logger.Errorf("policy hostname is not set")
		return nil, fmt.Errorf("policy hostname is not set")
	}

	agent, err := w.agentd.GetAgent(hostname)

	if err != nil {
		logger.Errorf("agent %s is not found", hostname)
		return nil, fmt.Errorf("agent %s is not found", hostname)
	}

	startTime, err := w.db.GetLastUpdateTime(p.BackupSource.SourceName)
	if err != nil {
		logger.Errorf("get last update time error: %v", err)
		return nil, err
	}
	path := genLogPath(p.Repository.Backend.Path, p.BackupSource.SourceName)
	metas, err = agent.Scan(ctx, path, startTime)

	if err != nil {
		logger.Errorf("scan error: %v", err)
		return nil, err
	}

	return metas, nil
}

func (w *LogWatcher) Watch() error {

	ctx, cancel := context.WithCancel(context.Background())
RunningLoop:
	for {
		select {
		case <-time.After(5 * time.Second):
			// scan for each policy
			w.LoopPolices(ctx)
		case fm := <-w.newFileMetaCh:
			if err := w.db.AddFileMeta(fm); err != nil {
				logger.Errorf("%v", err)
			}

		case <-time.After(5 * time.Second):
			logger.Info("heart beat")

		case <-w.stoppingCh:
			cancel()
			break RunningLoop
		}
	}

	return nil
}

func (w *LogWatcher) Stop() {
	logger.Info("stopping scheduler")
	w.stoppingCh <- struct{}{}
}
