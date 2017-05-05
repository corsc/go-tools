package fiximports

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUAT(t *testing.T) {
	scenarios := []struct {
		desc           string
		files          []string
		resultFilename string
	}{
		{
			desc:           "single file",
			files:          []string{"./test-data/test-a.in"},
			resultFilename: "./test-data/test-a.out",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			outputCapture := io.Writer(buffer)

			ProcessFiles(scenario.files, outputCapture)

			expected, err := ioutil.ReadFile(scenario.resultFilename)
			assert.Nil(t, err, scenario.desc)
			assert.Equal(t, expected, buffer.Bytes(), scenario.desc)
		})
	}

}
