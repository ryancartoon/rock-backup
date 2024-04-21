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
		Addr:    ":8000",
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
	r.POST("/service/file/open", GenOpenFileServiceHandler(a.ServiceEntry))
	r.GET("/service/file/get", GenGetPolicyHandler(a.ServiceEntry))
	// r.POST("/service/db/open", GenOpenDBServiceHandler(a.ServicEntry))

	return r
}

type BadRequestErr struct {
	message string
}

func (e *BadRequestErr) Error() string {
	return e.message
}

type OpenServiceRequest struct {
	SourcePath         string `json:"source_path"`
	Hostname           string `json:"hostname"`
	BackupPlan         uint   `json:"backup_plan"`
	Retention          uint   `json:"retention"`
	FullBackupSchedule string `json:"full_backup_schedule"`
	IncrBackupSchedule string `json:"incr_backup_schedule"`
	LogBackupSchedule  string `json:"log_backup_schedule"` // hours
	StartTime          string `json:"start_time"`
	RepositoryID       uint   `json:"repository_id"`
	// Duration           uint   `json:"duration"`
	// BackupCycle        uint   `json:"backup_cycle"`
}

type Response struct {
	Result interface{}
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

	// if err := verifyOpenServiceRequst(&request); err != nil {
	// 	return nil, err
	// }

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
			c.JSON(http.StatusOK, gin.H{"error": err})
			return
		}

		logger.Info("request is decoded.")

		source := &service.BackupSource{
			SourceType: "file",
			SourcePath: r.SourcePath,
		}

		policy := &service.Policy{
			Retention:    r.Retention,
			Status:       service.ServiceStatusServing,
			RepositoryID: r.RepositoryID,
			Hostname:     r.Hostname,
		}

		startTime, err := convStrToTime(r.StartTime)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err})
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

func GenGetPolicyHandler(s service.BackupServiceI) func(*gin.Context) {

	return func(c *gin.Context) {
		logger.Info("request is received.")

		ps, err := s.GetPolicies()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		fmt.Printf("%v", ps)

		c.JSON(http.StatusOK, ps)
	}
}

func GenOpenDBServiceHandler(s service.BackupServiceI) func(*gin.Context) {
	return func(c *gin.Context) {
		r, err := decodeServoceOpenReuqest(c)

		sourceType := "file"
		name := sourceType + r.SourcePath
		source := &service.BackupSource{
			SourceType: sourceType,
			SourcePath: r.SourcePath,
			SourceName: name,
		}

		startTime, err := convStrToTime(r.StartTime)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{})
		}

		full := schedules.Schedule{Cron: r.FullBackupSchedule, StartTime: startTime}
		incr := schedules.Schedule{Cron: r.IncrBackupSchedule, StartTime: startTime}
		log := schedules.Schedule{Cron: r.LogBackupSchedule, StartTime: startTime}

		var policy *service.Policy

		err = s.OpenFile(source, policy, []schedules.Schedule{full, incr, log})

		if err != nil {
			c.JSON(http.StatusOK, gin.H{})
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func GenCloseServiceHandler(s service.BackupServiceI) func(c *gin.Context) {
	return func(c *gin.Context) {
		// s.Close()
	}
}
