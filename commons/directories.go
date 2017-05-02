package commons

import (
	"os"
)

// IsDir returns true if the filename is a directory and false otherwise.
func IsDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}
