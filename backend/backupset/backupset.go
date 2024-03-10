package backupset

import (
	"gorm.io/gorm"
	"time"
)

type Backupset struct {
	gorm.Model
	ID        uint
	Size      uint
	ExpiredAt time.Time
}
