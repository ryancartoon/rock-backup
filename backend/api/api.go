package api

import (
	"context"
	"fmt"
	"net/http"
	"rockbackup/backend/schedules"
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

func New(service service.BackupServiceI) *WebAPI {
	return &WebAPI{ServiceEntry: service}
}

type WebAPI struct {
	ServiceEntry service.BackupServiceI
	server       *http.Server
}

func (a *WebAPI) Start() {
	router := a.NewRouter()
	a.server = &http.Server{
		Addr:    ":8080",
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
		logger.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		logger.Println("timeout of 5 seconds.")
	}
}

func (a *WebAPI) NewRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/service/file/open", GenOpenFileServiceHandler(a.ServiceEntry))
	// r.POST("/service/db/open", GenOpenDBServiceHandler(a.ServicEntry))

	return r
}

type BadRequestErr struct {
	message string
}

func (e *BadRequestErr) Error() string {
	return e.message
}

type BackupSource struct {
	// SourceType string `json:"source_type"`
	Name     string `json:"name"`
	DataPath string `json:"data_path"`
	Hostname string `json:"hostname"`
}

type OpenServiceRequest struct {
	Source             BackupSource `json:"source"`
	BackupPlan         uint         `json:"backup_plan"`
	FullBackupSchedule string       `json:"full_backup_schedule"`
	IncrBackupSchedule string       `json:"incr_backup_schedule"`
	LogBackupSchedule  string       `json:"log_backup_schedule"` // hours
	Retention          uint         `json:"retention"`
	BackupCycle        uint         `json:"backup_cycle"`
	StartTime          string       `json:"start_time"`
	Duration           uint         `json:"duration"`
	RepositoryID       uint         `json:"repository_id"`
}

type Response struct {
}

func verifyOpenServiceRequst(r *OpenServiceRequest) *BadRequestErr {
	if r.BackupPlan == 0 {
		return &BadRequestErr{fmt.Sprintf("backup plan %d", r.BackupPlan)}
	}
	return nil
}

func decodeServoceOpenReuqest(c *gin.Context) (*OpenServiceRequest, error) {
	var request OpenServiceRequest

	if err := c.BindJSON(&request); err != nil {
		return nil, err
	}

	if err := verifyOpenServiceRequst(&request); err != nil {
		return nil, err
	}

	return &request, nil
}

type CloseServiceRequest struct {
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

func GenOpenFileServiceHandler(s service.BackupServiceI) func(*gin.Context) {
	return func(c *gin.Context) {
		logger.Info("request is received.")
		r, err := decodeServoceOpenReuqest(c)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{})
			return
		}

		logger.Info("request is decoded.")

		source := &service.BackupSource{
			SourceType: "file",
			Name:       r.Source.Name,
			DataPath:   r.Source.DataPath,
		}

		policy := &service.Policy{
			Retention:    30,
			Status:       service.ServiceStatusServing,
			RepositoryID: r.RepositoryID,
		}

		startTime, err := convStrToTime(r.StartTime)
		if err != nil {
			c.JSON(http.StatusOK, policy)
			return
		}

		full := schedules.Schedule{Cron: r.FullBackupSchedule, StartTime: startTime}
		incr := schedules.Schedule{Cron: r.IncrBackupSchedule, StartTime: startTime}

		err = s.OpenFile(source, policy, []schedules.Schedule{full, incr})

		if err != nil {
			c.JSON(http.StatusOK, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

// func GenOpenDBServiceHandler(s service.BackupServiceI) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		r, err := decodeServoceOpenReuqest(c)
//
// 		source := &service.BackupSource{
// 			SourceType: r.BackupSource.SourceType,
// 			Name:       r.BackupSource.SourceName,
// 			DataPath:   r.BackupSource.DataPath,
// 		}
//
// 		startTime, err := convStrToTime(r.StartTime)
// 		if err != nil {
// 			c.JSON(http.StatusOK, gin.H{})
// 		}
//
// 		full := schedules.Schedule{Cron: r.FullBackupSchedule, StartTime: startTime}
// 		incr := schedules.Schedule{Cron: r.IncrBackupSchedule, StartTime: startTime}
// 		log := schedules.Schedule{Cron: r.LogBackupSchedule, StartTime: startTime}
//
// 		err = s.OpenFile(source, []schedules.Schedule{full, incr, log}, r.Retention)
//
// 		if err != nil {
// 			c.JSON(http.StatusOK, gin.H{})
// 		}
//
// 		c.JSON(http.StatusOK, gin.H{})
// 	}
// }

func GenCloseServiceHandler(s service.BackupServiceI) func(c *gin.Context) {
	return func(c *gin.Context) {
		// s.Close()
	}
}
