package alog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readFromLogWorks(log LogId, t *testing.T) {
	if _, err := os.Open(filepath.Join("/dev", "alog", log.String())); err != nil {
		t.Skipf("Android logging facilities are not accessible [%s]", err)
	}

	lr, err := NewLoggerReader(log)
	assert.Nil(t, err)

	atLeastOne := false

	for entry, err := lr.ReadNext(); err == nil; entry, err = lr.ReadNext() {
		atLeastOne = true
		t.Logf("%+v\n", entry)
	}

	assert.True(t, atLeastOne, "Log ", log, " is empty")
}

func TestReadFromLogsWorks(t *testing.T) {
	readFromLogWorks(LogIdMain, t)
	readFromLogWorks(LogIdRadio, t)
	readFromLogWorks(LogIdEvents, t)
	readFromLogWorks(LogIdSystem, t)
}
