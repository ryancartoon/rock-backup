package db

import (
	"rockbackup/backend/agentd"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
	"rockbackup/backend/service"
)

func (db *DB) LoadJob(id uint) (schedulerjob.Job, error) {
	var job schedulerjob.Job

	if result := db.g.First(&job, id); result.Error != nil {
		return job, result.Error
	}

	return job, nil
}

func (db *DB) LoadAgent(hostname string) (*agentd.Agent, error) {
	return nil, nil
}

func (db *DB) LoadRepository(id uint) (*repository.Repository, error) { return nil, nil }

func (db *DB) LoadPolicy(id uint) (service.Policy, error) {
	return service.Policy{}, nil
}
