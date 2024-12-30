package repository

import (
	"rockbackup/backend/backupset"

	"gorm.io/gorm"
)

type RepositoryI interface {
	SaveSnapshot() error
}

type RepositoryDB interface {
	LoadRepository(id uint) (*Repository, error)
	AddBackupset(repoID uint, jobID uint, backupTypte string) (*backupset.Backupset, error)
}

func LoadRepo(db RepositoryDB, id uint) (*Repository, error) {
	repo, err := db.LoadRepository(id)

	if err != nil {
		return nil, err
	}

	repo.db = db

	return repo, nil
}

type Repository struct {
	gorm.Model
	db         RepositoryDB `gorm:"-"`
	ID         uint         `gorm:"primaryKey;autoIncrement"`
	Name       string
	Type       string
	MountPoint string
}

func (r *Repository) AddBackupset(jobID uint, backupType string) (*backupset.Backupset, error) {
	return r.db.AddBackupset(r.ID, jobID, backupType)
}

// func (r *Repository) SetExpireAt(bSetID uint, retention uint) erorr {
// 	var expireAt *time.Time = time.Now().Add(retention * time.Day)
// 	return r.db.SetExpireAt(expireAt)
// }
