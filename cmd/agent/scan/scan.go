package scan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"rockbackup/backend/log"
	"time"
)

var logger *log.Logger

func init() {
	logName := "log-watcher"
	logger = log.New(logName)
}

type FileMeta struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex column:name"`
	Path      string    `gorm:"column:path"`
	Size      int64     `gorm:"column:size"`
	ModTime   time.Time `gorm:"column:mod_time"`
	Hash      string    `gorm:"column:hashid"`
	User      string    `gorm:"column:user"`
	Group     string    `gorm:"column:group"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdateAt  time.Time `gorm:"column:update_at"`
}

func calSHA(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func NewLogScaner() *LogScaner {
	return &LogScaner{}
}

type LogScaner struct{}

func (w *LogScaner) Scan(ctx context.Context, path string, t time.Time) ([]FileMeta, error) {

	var metas []FileMeta

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if info.ModTime().Before(t) {
				return nil
			}

			sha256Sum, err := calSHA(path)
			if err != nil {
				logger.Error(err)
			}

			file := FileMeta{
				Name:    info.Name(),
				Path:    path,
				Size:    info.Size(),
				ModTime: info.ModTime(),
				Hash:    sha256Sum,
			}

			metas = append(metas, file)
		}

		if err != nil {
			logger.Error(err)
		}

		return nil
	})

	if err != nil {
		logger.Printf("Error walking folder %s: %v", path, err)
	}

	return metas, nil
}

// func (w *LogWatcher) Watch(ctx context.Context) error {
// 	err := w.beforeWatch([]string{})

// 	if err != nil {
// 		return err
// 	}

// RunningLoop:
// 	for {
// 		select {
// 		case event, ok := <-w.watcher.Events:
// 			if !ok {
// 				logger.Errorf("%v", event)
// 				continue RunningLoop
// 			}

// 			if event.Op&fsnotify.Create == fsnotify.Create {
// 				logger.Printf("New file detected: %s", event.Name)

// 				info, err := os.Stat(event.Name)
// 				if err != nil {
// 					logger.Printf("Failed to stat file: %v", err)
// 					continue
// 				}

// 				if !info.IsDir() {
// 					sha256Sum, err := calSHA(event.Name)

// 					if err != nil {
// 						logger.Printf("Failed to compute SHA256 for file %s: %v", event.Name, err)
// 						continue
// 					}

// 					fm := FileMeta{
// 						Name:    info.Name(),
// 						Path:    event.Name,
// 						Size:    info.Size(),
// 						ModTime: info.ModTime(),
// 						SHA256:  sha256Sum,
// 					}

// 					w.newFileMetaCh <- fm

// 				}
// 			}
// 		case fm := <-w.newFileMetaCh:
// 			if err := w.AddNewFile(fm); err != nil {
// 				logger.Errorf("%v", err)
// 			}

// 		case path := <-w.addPathCh:
// 			if err := w.AddPath(path); err != nil {
// 				logger.Errorf("%v", err)
// 			}

// 		case <-time.After(5 * time.Second):
// 			logger.Info("heart beat")

// 		case err, ok := <-w.watcher.Errors:
// 			if !ok {
// 				logger.Error("watcher closed")
// 			}
// 			logger.Errorf("error: %s", err)
// 		case <-w.stoppingCh:
// 			w.watcher.Close()
// 			break RunningLoop
// 		case <-ctx.Done():
// 			w.watcher.Close()
// 			break RunningLoop
// 		}
// 	}

// 	return nil
// }

// func (w *LogWatcher) AddNewFile(file FileMeta) error {
// 	return w.db.AddFileMeta(file)
// }

// func (w *LogWatcher) AddPath(path string) error {
// 	err := w.watcher.Add(path)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (w *LogWatcher) RemovePath(path string) error {
// 	err := w.watcher.Remove(path)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// type DB interface {
// 	AddFileMeta(FileMeta FileMeta) error
// 	UpdateLastUpdateTime(t time.Time) error
// 	GetLastUpdateTime() (time.Time, error)
// }
