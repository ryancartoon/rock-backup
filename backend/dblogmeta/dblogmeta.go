package dblogmeta

import "time"

type FileMeta struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"column:name"`
	Path         string    `gorm:"column:path"`
	Size         int64     `gorm:"column:size"`
	ModTime      time.Time `gorm:"column:mod_time"`
	Hash         string    `gorm:"column:hashid;uniqueIndex"`
	InstanceName string    `gorm:"column:instance_name"`
	RepoID       uint      `gorm:"column:repo_id"`
	CreatedAt    time.Time
	UpdateAt     time.Time
}
