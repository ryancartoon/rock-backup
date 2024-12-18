package logbackup

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dbMock struct {
	lastUpdateTime time.Time
	fileMetas      []FileMeta
}

func (m *dbMock) AddFileMeta(meta FileMeta) error {
	m.fileMetas = append(m.fileMetas, meta)
	return nil
}

func (m *dbMock) UpdateLastUpdateTime(lastUpdateTime time.Time) error {
	return nil
}

func (m *dbMock) GetLastUpdateTime() (time.Time, error) {
	return m.lastUpdateTime, nil
}

func TestWatcherWatch(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// This should involve creating a mock database and watcher,
	// setting up a test directory with files, starting the watcher, and then
	// making sure it detects the file changes and updates the database accordingly.

	// Set up a test directory with some files.
	testDir1 := t.TempDir()
	testDir2 := t.TempDir()
	err := createTestFiles(testDir1, []string{"1"})
	assert.NoError(t, err)
	err = createTestFiles(testDir2, []string{"1"})
	assert.NoError(t, err)

	db := &dbMock{}
	watcher := NewLogWatcher(db)
	err = watcher.AddPath(testDir1)
	assert.NoError(t, err)
	err = watcher.AddPath(testDir2)
	assert.NoError(t, err)

	go watcher.Watch(ctx)

	time.Sleep(time.Second * 1)

	assert.Equal(t, len(db.fileMetas), 0) // Assuming two files were created initially.

	createTestFiles(testDir1, []string{"2"})
	time.Sleep(time.Second * 1)
	assert.Equal(t, len(db.fileMetas), 1)

	assert.Equal(t, "2", db.fileMetas[0].Name)
	assert.Equal(t, filepath.Join(testDir1, "2"), db.fileMetas[0].Path)
	assert.Equal(t, string("6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"), db.fileMetas[0].SHA256)

	createTestFiles(testDir2, []string{"2"})
	time.Sleep(time.Second * 1)
	assert.Equal(t, len(db.fileMetas), 2)

	assert.Equal(t, "2", db.fileMetas[0].Name)
	assert.Equal(t, filepath.Join(testDir2, "2"), db.fileMetas[1].Path)
	assert.Equal(t, string("6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"), db.fileMetas[1].SHA256)

}

func createTestFiles(dir string, fileNames []string) error {
	for _, fileName := range fileNames {
		err := os.WriteFile(filepath.Join(dir, fileName), []byte("test content"), 0644)

		if err != nil {
			return err
		}
	}

	return nil
}
