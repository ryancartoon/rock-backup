package keeperr

import (
	"rockbackup/backend/log"
)

var logger *log.Logger

func init() {
	logger = log.New("keeperr")
}

type KeepErr struct {
	db         DB
	errCh      chan error
	stoppingCh chan struct{}
}

type DB interface {
	SaveErr(err error) error
}

func NewKeepErr() *KeepErr {
	return &KeepErr{
		errCh:      make(chan error, 1000),
		stoppingCh: make(chan struct{}),
	}
}

func (k *KeepErr) SaveErr(er error) {
	err := k.db.SaveErr(er)

	if err != nil {
		logger.Errorf("save err failed: %v", err)
	}
}

func (k *KeepErr) Start() {

RunningLoop:
	for {
		select {
		case <-k.stoppingCh:
			break RunningLoop
		case err := <-k.errCh:
			k.SaveErr(err)
		}
	}
}

func (k *KeepErr) Stop() {
	k.stoppingCh <- struct{}{}
}
