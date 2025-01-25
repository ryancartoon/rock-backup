package api

import (
	"context"
	"fmt"
	"net/http"
	"rockbackup/backend/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

var (
	NilTime = datatypes.NewTime(0, 0, 0, 0)
)

type BaseBackupPolicy struct {
	Hostname           string `json:"hostname"`
	BackupPlan         uint   `json:"backup_plan"`
	Retention          uint   `json:"retention"`
	FullBackupSchedule string `json:"full_backup_schedule"`
	IncrBackupSchedule string `json:"incr_backup_schedule"`
	LogBackupSchedule  string `json:"log_backup_schedule"` // hours
	StartTime          string `json:"start_time"`
	BackendID          uint   `json:"backend_id"`
	Duration           uint   `json:"duration"`
}

func New(s service.BackupServiceI) *WebAPI {
	return &WebAPI{ServiceEntry: s}
}

type WebAPI struct {
	ServiceEntry service.BackupServiceI
	server       *http.Server
}

func (a *WebAPI) Start() {
	router := a.NewRouter()
	a.server = &http.Server{
		Addr:    "localhost:8000",
		Handler: router,
	}

	// service connections
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("listen: %s\n", err)
	}

}

func (a *WebAPI) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server Shutdown: %s", err)
	}

	select {
	case <-ctx.Done():
		logger.Info("timeout of 5 seconds.")
	}
}

func (a *WebAPI) NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.POST("/service/file/open", GenOpenResticFileServiceHandler(a.ServiceEntry))
	r.POST("/service/mysql/open", GenOpenResticFileServiceHandler(a.ServiceEntry))
	r.GET("/service/file/get", GenGetPolicyHandler(a.ServiceEntry))
	r.POST("/backup/job", GenStartFileBackupJobHandler(a.ServiceEntry))
	r.POST("/restore/job", GenStartFileRestoreJobHandler(a.ServiceEntry))
	// r.POST("/service/db/open", GenOpenDBServiceHandler(a.ServicEntry))

	return r
}

type BadRequestErr struct {
	message string
}

func (e *BadRequestErr) Error() string {
	return e.message
}

type Response struct {
	Result interface{}
}

func convStrToTime(s string) (datatypes.Time, error) {
	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		return NilTime, fmt.Errorf("bad format of time: %v", s)
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return NilTime, err
	}

	min, err := strconv.Atoi(parts[1])
	if err != nil {
		return NilTime, err
	}

	return datatypes.NewTime(hour, min, 0, 0), nil
}
