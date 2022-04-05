package utils

import (
	"log"
	"testing"
)

// testingFatal can be set to true to reduce the severity from fatal and exit to
// a log message. This is useful when testing a fatal behavior.
// Note: changing the value is not thread-safe. Run tests sequentially by setting "-p 1" flag
var testingFatal = false

func logf(t *testing.T, format string, args ...interface{}) {
	if t != nil {
		t.Logf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func fatalf(t *testing.T, format string, args ...interface{}) {
	if testingFatal {
		logf(t, "fatal error: "+format, args...)
		return
	}

	if t != nil {
		t.Fatalf(format, args...)
	} else {
		log.Fatalf(format, args...)
	}
}
