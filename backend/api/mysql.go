package api

import (
	"net/http"
	"rockbackup/backend/log"
	"rockbackup/backend/service"

	"github.com/gin-gonic/gin"
)

var logger = log.New("api.log")

type OpenMysqlXtrabackupServiceRequest struct {
	InstanceDataPath  string `json:"instance_data_path"`
	InstanceName      string `json:"instance_name"`
	MysqlVersion      string `json:"mysql_version"`
	XtrabackupVersion string `json:"xtrabackup_version"`
	BackupCycle       uint   `json:"backup_cycle"`
	BaseBackupPolicy
}

func GenOpenMysqlXtrabackupFileServiceHandler(s service.BackupServiceI) func(*gin.Context) {
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
		policyReq.BackupSourceName = r.SourceName
		policyReq.Retention = r.Retention
		policyReq.BackendID = r.BackendID
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
