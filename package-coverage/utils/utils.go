package utils

import "log"

// LogWhenVerbose will write debugging information to stdout when verbose mode is requested
var LogWhenVerbose = func(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}

// VerboseOff will stop all verbose logs
func VerboseOff() {
	LogWhenVerbose = func(msg string, args ...interface{}) {}
}
