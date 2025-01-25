package api

import (
	"fmt"
	"net/http"
	"rockbackup/backend/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OpenResticServiceRequest struct {
	SourcePath  string `json:"source_path"`
	SourceName  string `json:"source_name"`
	BackupCycle uint   `json:"backup_cycle"`
	BaseBackupPolicy
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

// func verifyOpenServiceRequst(r *OpenResticServiceRequest) *BadRequestErr {
// 	if r.BackupPlan == 0 {
// 		return &BadRequestErr{fmt.Sprintf("backup plan %d", r.BackupPlan)}
// 	}
// 	return nil
// }

func decodeServoceOpenReuqest(c *gin.Context) (OpenResticServiceRequest, error) {
	var request OpenResticServiceRequest

	if err := c.BindJSON(&request); err != nil {
		return OpenResticServiceRequest{}, err
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

func GenOpenResticFileServiceHandler(s service.BackupServiceI) func(*gin.Context) {
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
