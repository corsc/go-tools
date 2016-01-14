package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateCoverage(t *testing.T) {
	expected := map[string]*coverage{
		"sage42.org/go-tools/package-coverage/": {
			selfStatements: 87,
			selfCovered:    55,

			childStatements: 0,
			childCovered:    0,
		},
	}

	result := calculateCoverage(sampleCoverageFileContents)
	converted := map[string]*coverage(result)
	assert.Equal(t, expected, converted)
}

var sampleCoverageFileContents = `
mode: set
sage42.org/go-tools/package-coverage/file_parser.go:19.50,23.21 3 1
sage42.org/go-tools/package-coverage/file_parser.go:28.2,31.29 3 1
sage42.org/go-tools/package-coverage/file_parser.go:38.2,40.9 2 1
sage42.org/go-tools/package-coverage/file_parser.go:45.2,45.15 1 1
sage42.org/go-tools/package-coverage/file_parser.go:23.21,26.3 2 0
sage42.org/go-tools/package-coverage/file_parser.go:31.29,32.29 1 1
sage42.org/go-tools/package-coverage/file_parser.go:36.3,36.32 1 1
sage42.org/go-tools/package-coverage/file_parser.go:32.29,33.12 1 1
sage42.org/go-tools/package-coverage/file_parser.go:41.2,41.16 0 1
sage42.org/go-tools/package-coverage/file_parser.go:48.40,50.2 1 1
sage42.org/go-tools/package-coverage/file_parser.go:52.92,55.12 2 1
sage42.org/go-tools/package-coverage/file_parser.go:66.2,66.15 1 1
sage42.org/go-tools/package-coverage/file_parser.go:55.12,56.36 1 1
sage42.org/go-tools/package-coverage/file_parser.go:61.3,63.16 2 1
sage42.org/go-tools/package-coverage/file_parser.go:56.36,59.4 2 1
sage42.org/go-tools/package-coverage/file_parser.go:69.77,71.9 2 1
sage42.org/go-tools/package-coverage/file_parser.go:75.2,75.14 1 1
sage42.org/go-tools/package-coverage/file_parser.go:71.9,74.3 2 1
sage42.org/go-tools/package-coverage/file_parser.go:78.62,80.22 2 1
sage42.org/go-tools/package-coverage/file_parser.go:80.22,82.3 1 1
sage42.org/go-tools/package-coverage/file_parser.go:85.55,86.46 1 1
sage42.org/go-tools/package-coverage/file_parser.go:86.46,87.47 1 1
sage42.org/go-tools/package-coverage/file_parser.go:87.47,88.28 1 1
sage42.org/go-tools/package-coverage/file_parser.go:92.4,92.35 1 1
sage42.org/go-tools/package-coverage/file_parser.go:88.28,89.13 1 1
sage42.org/go-tools/package-coverage/file_parser.go:92.35,95.5 2 1
sage42.org/go-tools/package-coverage/file_parser.go:100.45,102.2 1 1
sage42.org/go-tools/package-coverage/line_parser.go:15.37,23.2 3 1
sage42.org/go-tools/package-coverage/line_parser.go:25.40,27.21 2 1
sage42.org/go-tools/package-coverage/line_parser.go:31.2,31.30 1 1
sage42.org/go-tools/package-coverage/line_parser.go:27.21,28.69 1 1
sage42.org/go-tools/package-coverage/line_parser.go:34.45,36.21 2 1
sage42.org/go-tools/package-coverage/line_parser.go:40.2,43.23 3 1
sage42.org/go-tools/package-coverage/line_parser.go:36.21,37.83 1 1
sage42.org/go-tools/package-coverage/line_parser.go:46.40,48.16 2 1
sage42.org/go-tools/package-coverage/line_parser.go:51.2,51.19 1 1
sage42.org/go-tools/package-coverage/line_parser.go:48.16,49.13 1 1
sage42.org/go-tools/package-coverage/line_parser.go:54.38,56.2 1 1
sage42.org/go-tools/package-coverage/main.go:10.13,11.15 1 0
sage42.org/go-tools/package-coverage/main.go:17.2,25.35 5 0
sage42.org/go-tools/package-coverage/main.go:11.15,12.31 1 0
sage42.org/go-tools/package-coverage/main.go:12.31,14.4 1 0
sage42.org/go-tools/package-coverage/main.go:28.27,32.20 3 0
sage42.org/go-tools/package-coverage/main.go:36.2,36.17 1 0
sage42.org/go-tools/package-coverage/main.go:32.20,33.48 1 0
sage42.org/go-tools/package-coverage/main.go:39.46,41.16 2 0
sage42.org/go-tools/package-coverage/main.go:45.2,45.25 1 0
sage42.org/go-tools/package-coverage/main.go:41.16,42.13 1 0
sage42.org/go-tools/package-coverage/main.go:48.68,51.32 2 0
sage42.org/go-tools/package-coverage/main.go:55.2,57.15 2 0
sage42.org/go-tools/package-coverage/main.go:51.32,53.3 1 0
sage42.org/go-tools/package-coverage/main.go:60.70,63.27 2 0
sage42.org/go-tools/package-coverage/main.go:73.2,73.15 1 0
sage42.org/go-tools/package-coverage/main.go:63.27,72.3 5 0
`
