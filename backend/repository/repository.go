package repository

import (
	"path/filepath"
	"rockbackup/backend/backupset"
)

type BackendI interface {
	Link()
}

type RepositoryI interface {
	SaveSnapshot() error
}

type RepositoryDB interface {
	AddBackupset(RepoName string, BackendID uint, jobID uint, backupTypte string) (*backupset.Backupset, error)
}

// func NewRepository(name string, backend *Backend, db RepositoryDB) *Repository {
// 	return &Repository{
// 		db:      db,
// 		Name:    name,
// 		Backend: backend,
// 	}
// }

type Repository struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Name      string
	IsActive  bool     `gorm:"default:true"`
	Backend   *Backend `gorm:"_"`
	BackendID uint     `gorm:"Column:backend_id;not null"`
	PolicyID  uint     `gorm:"column:policy_id"`
	// db        RepositoryDB `gorm:"_"`
}

func (r *Repository) AddBackupset(db RepositoryDB, jobID uint, backupType string) (*backupset.Backupset, error) {
	return db.AddBackupset(r.Name, r.Backend.ID, jobID, backupType)
}

func (r *Repository) GetPath() string {
	return filepath.Join(r.Backend.Path, r.Name)
}

// func (r *Repository) SetExpireAt(bSetID uint, retention uint) erorr {
// 	var expireAt *time.Time = time.Now().Add(retention * time.Day)
// 	return r.db.SetExpireAt(expireAt)
// }
