package async

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rockbackup/backend/async/taskdef"
	"rockbackup/backend/schedulerjob"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)


func MakeHandleBackupFileTask(config *viper.Viper, db FactoryDB, jobDB schedulerjob.JobDB) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var p taskdef.BackupJobPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}
		log.Printf("handle job id: %d", p.ID)
		factory := &Factory{db}
		factory.StartBackupJobFile(ctx, p, jobDB)
		return nil
	}
}
