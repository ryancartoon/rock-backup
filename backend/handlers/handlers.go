package handlers

func New() *Handler {
	return &Handler{}
}

type Handler struct{}

func (h *Handler) TimerStartBackup(policyID uint, backupType string, operator string) error {
	return nil
}
