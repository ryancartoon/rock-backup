package async

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rockbackup/backend/async/taskdef"

	"github.com/hibiken/asynq"
)

func HandleBackupFileTask(ctx context.Context, t *asynq.Task) error {
	var p taskdef.BackupJobPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("handle job id: %d", p.ID)
	factory := &Factory{DB}
	factory.StartBackupJobFile(ctx, p, DB)
	return nil
}
