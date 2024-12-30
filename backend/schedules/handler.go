package schedules

type JobStarter interface {
	AddSchedulerJobBackup(policyID uint, backupType string, operator string) error
}

func NewHandler(starter JobStarter) *Handler {
	return &Handler{jobStarter: starter}
}

type Handler struct {
	jobStarter JobStarter
}

func (h *Handler) TimerStartBackup(policyID uint, backupType string, operator string) error {
	logger.Infof("Starting job for policy %v", policyID)
	// send job to scheduler by a channel
	return h.jobStarter.AddSchedulerJobBackup(policyID, backupType, operator)
}
