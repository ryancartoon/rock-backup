package repository

import (
	"gorm.io/gorm"
)

type RepositoryI interface {
	SaveSnapshot() error
}

type Repository struct {
	gorm.Model
	ID         uint `gorm:"primaryKey;autoIncrement"`
	Name       string
	Type       string
	MountPoint string
}
