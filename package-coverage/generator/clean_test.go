// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCleanerImpl_Single(t *testing.T) {
	path := "./test-data/"

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
