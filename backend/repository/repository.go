package repository

import (
	"path/filepath"
	"rockbackup/backend/backupset"
)

type RepositoryI interface {
	SaveSnapshot() error
}

type RepositoryDB interface {
	AddBackupset(RepoName string, BackendID uint, jobID uint, backupTypte string) (*backupset.Backupset, error)
}

func NewRepository(name string, backend *Backend, db RepositoryDB) *Repository {
	return &Repository{
		db:      db,
		Name:    name,
		Backend: backend,
	}
}

type Repository struct {
	db      RepositoryDB
	Name    string
	Backend *Backend
}

func (r *Repository) AddBackupset(jobID uint, backupType string) (*backupset.Backupset, error) {
	return r.db.AddBackupset(r.Name, r.Backend.ID, jobID, backupType)
}

func (r *Repository) GetTarget() string {
	return filepath.Join(r.Backend.Path, r.Name)
}

// func (r *Repository) SetExpireAt(bSetID uint, retention uint) erorr {
// 	var expireAt *time.Time = time.Now().Add(retention * time.Day)
// 	return r.db.SetExpireAt(expireAt)
// }
