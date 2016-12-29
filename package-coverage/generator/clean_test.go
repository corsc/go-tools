package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCleanerImpl_Single(t *testing.T) {
	path := "./test/"

	scenarios := []struct {
		desc      string
		setupMock func(mockFs *mockFsWrapper)
	}{
		{
			desc: "file exists, delete should be called",
			setupMock: func(mockFs *mockFsWrapper) {
				mockFs.On("Exists", path+coverageFilename).Once().Return(true)
				mockFs.On("Delete", path+coverageFilename).Once()
			},
		},
		{
			desc: "file does not exist, delete should not be called",
			setupMock: func(mockFs *mockFsWrapper) {
				mockFs.On("Exists", path+coverageFilename).Once().Return(false)
			},
		},
	}

	for _, scenario := range scenarios {
		cleaner, mockFs := newTestCleaner()
		scenario.setupMock(mockFs)

		cleaner.Single(path)
		assert.True(t, mockFs.AssertExpectations(t), scenario.desc)
	}
}

func newTestCleaner() (Cleaner, *mockFsWrapper) {
	mockFs := &mockFsWrapper{}

	return &cleanerImpl{
		fsWrapper: mockFs,
	}, mockFs
}

// mock implementation of the fsWrapper
type mockFsWrapper struct {
	mock.Mock
}

// Exists implements fsWrapper
func (fs *mockFsWrapper) Exists(filename string) bool {
	outputs := fs.Mock.Called(filename)
	return outputs.Bool(0)
}

// Delete implements fsWrapper
func (fs *mockFsWrapper) Delete(filename string) {
	fs.Mock.Called(filename)
}
