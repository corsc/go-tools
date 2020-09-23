package log

import (
	"fmt"
	"log"
	"os"
)

func Error(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	if err != nil {
		log.Fatalf("failed to write to stdErr with err: %s", err)
	}
}
