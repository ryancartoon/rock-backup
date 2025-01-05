package db

import (
	"rockbackup/backend/agentd"
	"rockbackup/backend/host"
	"rockbackup/backend/policy"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
)

func (db *DB) LoadJob(id uint) (*schedulerjob.Job, error) {
	var job schedulerjob.Job

	if result := db.g.First(&job, id); result.Error != nil {
		return nil, result.Error
	}

	return &job, nil
}

func (db *DB) LoadAgent(hostname string) (*agentd.Agent, error) {
	var host host.Host
	if result := db.g.Where("name = ?", hostname).Find(&host); result.Error != nil {
		return nil, result.Error
	}

	agent := &agentd.Agent{Host: hostname, Port: host.AgentPort}

	return agent, nil
}

func (db *DB) LoadBackend(id uint) (*repository.Backend, error) {
	var b repository.Backend
	if result := db.g.First(&b, id); result.Error != nil {
		return nil, result.Error
	}

	return &b, nil
}

func (db *DB) LoadPolicy(id uint) (*policy.Policy, error) {
	var p policy.Policy
	if result := db.g.InnerJoins("BackupSource").First(&p, id); result.Error != nil {
		return nil, result.Error
	}
	return &p, nil
}
