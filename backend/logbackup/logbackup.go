package logbackup

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"rockbackup/backend/log"
	"time"

	"github.com/fsnotify/fsnotify"
)

var logger *log.Logger

func init() {
	logName := "log-watcher"
	logger = log.New(logName)
}

type DB interface {
	AddFileMeta(FileMeta FileMeta) error
	UpdateLastUpdateTime(t time.Time) error
}

type FileMeta struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	Path      string    `gorm:"not null"`
	Size      int64     `gorm:"not null"`
	ModTime   time.Time `gorm:"not null"`
	SHA256    string    `gorm:"not null"`
	CreatedAt time.Time
	UpdateAt  time.Time
}

func NewWatcher() *LogWatcher {
	fsWatcher, err := fsnotify.NewWatcher()

	if err != nil {
		panic("Failed to create watcher")
	}

	filemetaCh := make(chan FileMeta)

	return &LogWatcher{
		watcher:    fsWatcher,
		fileMetaCh: filemetaCh,
		stoppingCh: make(chan struct{}),
		addPathCh:  make(chan string),
	}
}

type LogWatcher struct {
	watcher      *fsnotify.Watcher
	fileMetaCh   chan FileMeta
	stoppingCh   chan struct{}
	addPathCh    chan string
	removePathCh chan string
	db           DB
}

func (w *LogWatcher) AddPath(path string) error {
	err := w.watcher.Add(path)

	if err != nil {
		return err
	}

	return nil
}

func (w *LogWatcher) RemovePath(path string) error {
	err := w.watcher.Remove(path)

	if err != nil {
		return err
	}

	return nil
}

func (w *LogWatcher) beforeWatch(paths []string) error {
	// scan paths

	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				sha256Sum, err := calSHA(path)
				if err != nil {
					logger.Error(err)
				}

				file := FileMeta{
					Name:    info.Name(),
					Path:    path,
					Size:    info.Size(),
					ModTime: info.ModTime(),
					SHA256:  sha256Sum,
				}

				w.db.AddFileMeta(file)
			}

			return nil
		})

		if err != nil {
			logger.Printf("Error walking folder %s: %v", path, err)
		}
	}

	return nil
}

func (w *LogWatcher) Watch() error {
	err := w.beforeWatch()

	if err != nil {
		return err
	}

RunningLoop:
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				logger.Errorf("%v", event)
				continue RunningLoop
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				logger.Printf("New file detected: %s", event.Name)

				info, err := os.Stat(event.Name)
				if err != nil {
					logger.Printf("Failed to stat file: %v", err)
					continue
				}

				if !info.IsDir() {
					sha256Sum, err := calSHA(event.Name)

					if err != nil {
						logger.Printf("Failed to compute SHA256 for file %s: %v", event.Name, err)
						continue
					}

					fm := FileMeta{
						Name:    info.Name(),
						Path:    event.Name,
						Size:    info.Size(),
						ModTime: info.ModTime(),
						SHA256:  sha256Sum,
					}

					w.fileMetaCh <- fm

				}
			}
		case fm := <-w.fileMetaCh:
			if err := w.AddNewFile(fm); err != nil {
				logger.Errorf("%v", err)
			}

		case path := <-w.addPathCh:
			w.AddPath(path)

		case <-time.After(5 * time.Second):
			logger.Info("heart beat")

		case err, _ := <-w.watcher.Errors:
			logger.Printf("Watcher error: %v", err)
		case <-w.stoppingCh:
			break RunningLoop
		}
	}

	return nil

}

func (w *LogWatcher) AddNewFile(file FileMeta) error {
	return w.db.AddFileMeta(file)
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
