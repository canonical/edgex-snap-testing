package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// global setup

	code := m.Run()

	// global teardown

	os.Exit(code)
}
