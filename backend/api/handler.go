package api

import (
	"context"
	"fmt"
	"net/http"
	"rockbackup/backend/log"
	"rockbackup/backend/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/datatypes"
)

var logger = log.New("api.log")

// var logger *logrus.Entry

// func init() {
// 	logName := "backend.log"
// 	l := log.New(logName)
// 	logger = l.WithField("compoent", "api")
// }

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
	Duration           uint   `json:"duration"`
	BackupCycle        uint   `json:"backup_cycle"`
}

type StartBackupJobRequest struct {
	PolicyID   string `json:"policy_id"`
	BackupType string `json:"backup_type"`
}

type StartFileRestoreJobRequest struct {
	PolicyID    uint   `json:"policy_id"`
	BackupsetID uint   `json:"backupset_id"`
	TargetPath  string `json:"target_path"`
}

var (
	NilTime = datatypes.NewTime(0, 0, 0, 0)
)

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
		Addr:    "0.0.0.0:8000",
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
		logger.Fatalf("Server Shutdown:", err)
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/service/file/open", GenOpenFileServiceHandler(a.ServiceEntry))
	r.POST("/service/mysql/open", GenOpenFileServiceHandler(a.ServiceEntry))
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

func verifyOpenServiceRequst(r *OpenServiceRequest) *BadRequestErr {
	if r.BackupPlan == 0 {
		return &BadRequestErr{fmt.Sprintf("backup plan %d", r.BackupPlan)}
	}
	return nil
}

func decodeServoceOpenReuqest(c *gin.Context) (OpenServiceRequest, error) {
	var request OpenServiceRequest

	if err := c.BindJSON(&request); err != nil {
		return OpenServiceRequest{}, err
	}

	// if err := verifyOpenServiceRequst(&request); err != nil {
	// 	return nil, err
	// }

	return request, nil
}

func decocdeStartBackupJobRequest(c *gin.Context) (StartBackupJobRequest, error) {
	var r StartBackupJobRequest

	if err := c.BindJSON(&r); err != nil {
		return StartBackupJobRequest{}, err
	}

	return r, nil
}

func decocdeStartFileRestoreJobRequest(c *gin.Context) (StartFileRestoreJobRequest, error) {
	var r StartFileRestoreJobRequest

	if err := c.BindJSON(&r); err != nil {
		return StartFileRestoreJobRequest{}, err
	}

	return r, nil
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

type PolicyRequest struct {
	Retention          uint           `json:"album_name" binding:"required"`
	BackupSourcePath   string         `json:"backup_source_path" binding:"required"`
	Hostname           string         `json:"hostname" binding:"required"`
	RepositoryID       uint           `json:"repository_id" binding:"required"`
	BackupSourceID     uint           `json:"backup_source_id" binding:"required"`
	FullBackupSchedule string         `json:"full_backup_schedule" binding:"required"`
	IncrBackupSchedule string         `json:"incr_backup_schedule" binding:"required"`
	ScheduleDesc       string         `json:"schedule_desc" binding:"required"`
	StartTime          datatypes.Time `json:"start_time" binding:"required"`
	BackupCycle        uint           `json:"backup_cycle" binding:"required"`
}

// CreateBackupPolicy godoc
//
//	@Summary		Create backup policy
//	@Description	create backup policy
//	@Tags			Policy
//	@Accept			json
//	@Produce		json
//	@Router			/service/file/open [post]
//	@Param			policy	body		PolicyRequest	true	"backup policy params"
//	@Success		200		{object}	Response
func GenOpenFileServiceHandler(s service.BackupServiceI) func(*gin.Context) {
	return func(c *gin.Context) {
		var policyReq service.PolicyRequest

		logger.Info("request is received.")
		r, err := decodeServoceOpenReuqest(c)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err})
			return
		}

		logger.Info("request is decoded.")

		policyReq.BackupSourcePath = r.SourcePath
		policyReq.Hostname = r.Hostname
		policyReq.Retention = r.Retention
		policyReq.RepositoryID = r.RepositoryID
		policyReq.BackupCycle = r.BackupCycle

		//TODO:  verify schedules
		policyReq.FullBackupSchedule = r.FullBackupSchedule
		policyReq.IncrBackupSchedule = r.IncrBackupSchedule

		//veirfy time
		policyReq.StartTime, err = convStrToTime(r.StartTime)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err})
			return
		}

		err = s.OpenFile(policyReq)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

type PolicyView struct {
	ID             uint   `json:"id"`
	SourceID       uint   `json:"source_id"`
	SourceType     string `json:"source_type"`
	SourcePath     string `json:"source_path"`
	SourceHost     string `json:"source_host"`
	Retention      uint   `json:"retention"`
	BackupSourceID uint   `json:"backup_source_id"`
	Hostname       string `json:"hostname"`
	Status         string `json:"status"`
	RepositoryID   uint   `json:"repository_id"`
}

// GetBackupPolicy godoc
//
//	@Summary		Get backup policy
//	@Description	get backup policy
//	@Tags			Policy
//	@Accept			json
//	@Produce		json
//	@Router			/service/file/get [get]
//	@Success		200	{object} PolicyView
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

// StartBackupJob godoc
//
//	@Summary		Start backup job
//	@Description	start backup job
//	@Tags			BackupJob
//	@Accept			json
//	@Produce		json
//	@Router			/backup/job [post]
//	@Param			policy	body		StartBackupJobRequest	true	"backup job params"
//	@Success		200		{object}	Response
func GenStartFileBackupJobHandler(s service.BackupServiceI) func(*gin.Context) {
	return func(c *gin.Context) {
		logger.Info("request is received.")

		r, err := decocdeStartBackupJobRequest(c)
		logger.Infof("%v\n", r)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err})
			return
		}

		logger.Info("request is decoded.")

		policyID, err := strconv.ParseUint(r.PolicyID, 10, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err = s.StartBackupJob(uint(policyID), r.BackupType)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, nil)

	}
}

func GenStartFileRestoreJobHandler(s service.BackupServiceI) func(*gin.Context) {
	return func(c *gin.Context) {
		logger.Info("request is received.")

		r, err := decocdeStartFileRestoreJobRequest(c)
		logger.Infof("%v\n", r)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err})
			return
		}

		logger.Info("request is decoded.")

		err = s.StartRestoreJob(r.PolicyID, r.BackupsetID, r.TargetPath)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, nil)

	}
}

// func GenStartRestoreJobHandler(s service.BackupServiceI) func(*gin.Context) {
// 	return func(c *gin.Context) {
// 		logger.Info("request is received.")
//
// 		r, err := decocdeStartRestoreJobRequest(c)
// 		logger.Infof("%v\n", r)
// 		if err != nil {
// 			c.JSON(http.StatusOK, gin.H{"error": err})
// 			return
// 		}
//
// 		logger.Info("request of restore is decoded.")
//
// 		policyID, err := strconv.ParseUint(r.PolicyID, 10, 0)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err})
// 			return
// 		}
//
// 		err = s.StartBackupJob(uint(policyID), r.BackupType)
//
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err})
// 			return
// 		}
//
// 		c.JSON(http.StatusOK, nil)
//
// 	}
// }

func GenCloseServiceHandler(s service.BackupServiceI) func(c *gin.Context) {
	return func(c *gin.Context) {
		// s.Close()
	}
}
