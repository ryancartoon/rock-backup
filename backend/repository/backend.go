package repository

import "gorm.io/gorm"

type Backend struct {
	gorm.Model
	ID   uint `gorm:"primaryKey;autoIncrement"`
	Name string
	Type string
	Path string
}

func (b *Backend) GetTargetRoot() string {
	return b.Path
}
