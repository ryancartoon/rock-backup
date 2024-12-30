package scan

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWatcherWatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a test directory with some files.
	testDir1 := t.TempDir()
	err := createTestFiles(testDir1, []string{"1"})
	assert.NoError(t, err)

	scanner := NewLogScaner()

	time.Sleep(10 * time.Millisecond)
	now := time.Now()

	metas, err := scanner.Scan(ctx, testDir1, now)
	assert.NoError(t, err)

	assert.Equal(t, len(metas), 0)

	time.Sleep(10 * time.Millisecond)
	err = createTestFiles(testDir1, []string{"2"})
	assert.NoError(t, err)

	metas, err = scanner.Scan(ctx, testDir1, now)
	assert.NoError(t, err)

	assert.Equal(t, len(metas), 1)
	assert.Equal(t, "2", metas[0].Name)
	assert.Equal(t, filepath.Join(testDir1, "2"), metas[0].Path)
	assert.Equal(t, string("6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"), metas[0].Hash)

	err = createTestFiles(testDir1, []string{"3"})
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(testDir1, "3"), []byte("hello world"), 0644)
	assert.NoError(t, err)

	metas, err = scanner.Scan(ctx, testDir1, now)
	assert.NoError(t, err)
	assert.Equal(t, len(metas), 2)
	assert.Equal(t, "3", metas[1].Name)
	assert.Equal(t, filepath.Join(testDir1, "3"), metas[1].Path)
	assert.Equal(t, string("b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"), metas[1].Hash)
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

// type dbMock struct {
// 	lastUpdateTime time.Time
// 	fileMetas      []FileMeta
// }

// func (m *dbMock) AddFileMeta(meta FileMeta) error {
// 	m.fileMetas = append(m.fileMetas, meta)
// 	return nil
// }

// func (m *dbMock) UpdateLastUpdateTime(lastUpdateTime time.Time) error {
// 	return nil
// }

// func (m *dbMock) GetLastUpdateTime() (time.Time, error) {
// 	return m.lastUpdateTime, nil
// }
