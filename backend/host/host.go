package host

import (
	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	ID       int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name     string `gorm:"column:name"`
	HostType string `gorm:"column:host_type"`
	Location string `gorm:"column:location"`
	IsActive bool   `gorm:"column:is_active"`
	Load     int    `gorm:"column:load"`
}

func TableName() string {
	return "hosts"
}
