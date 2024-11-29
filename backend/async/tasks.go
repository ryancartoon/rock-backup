package async

import (
	"context"
	"encoding/json"
	"fmt"
	"rockbackup/backend/async/taskdef"
	"rockbackup/backend/log"
	"rockbackup/backend/schedulerjob"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

var logger = log.New("worker-tasks.log")

func MakeHandleBackupFileTask(config *viper.Viper, db DB, jobDB schedulerjob.JobDB) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var p taskdef.BackupJobPayload

		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}
		logger.Infof("handle job id: %d", p.ID)
		logger.Infof("job input: %v", p)
		factory := &Factory{db}
		factory.StartBackupJobFile(ctx, p, jobDB)

		return nil
	}
}
