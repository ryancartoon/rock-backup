package keeperr

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the DB interface
type MockDB struct {
	mock.Mock
	errs []error
}

func (m *MockDB) SaveErr(err error) error {
	args := m.Called(err)
	m.errs = append(m.errs, err)
	return args.Error(0)
}

func TestNewKeepErr(t *testing.T) {
	ke := NewKeepErr()
	assert.NotNil(t, ke)
	assert.NotNil(t, ke.errCh)
}

func TestKeepErr_SaveErr(t *testing.T) {
	mockDB := new(MockDB)
	ke := &KeepErr{
		db:    mockDB,
		errCh: make(chan error, 1000),
	}

	testErr := errors.New("test error")
	mockDB.On("SaveErr", testErr).Return(nil)

	ke.SaveErr(testErr)

	mockDB.AssertCalled(t, "SaveErr", testErr)
}

func TestKeepErr_Start(t *testing.T) {
	mockDB := new(MockDB)
	ke := &KeepErr{
		db:    mockDB,
		errCh: make(chan error, 1000),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ke.Start(ctx)

	testErr := errors.New("test error")
	mockDB.On("SaveErr", testErr).Return(nil)

	ke.SaveErr(testErr)
	time.Sleep(100 * time.Millisecond) // Give some time for the goroutine to process

	mockDB.AssertCalled(t, "SaveErr", testErr)
	assert.Contains(t, mockDB.errs, testErr)
	cancel()
}
