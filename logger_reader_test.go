package alog

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func readFromLogWorks(log LogId, t *testing.T) {
	if _, err := os.Open(filepath.Join("/dev", "alog", log.String())); err != nil {
		t.Skipf("Android logging facilities are not accessible [%s]", err)
	}

	lr, err := NewLoggerReader(log)
	assert.Nil(t, err)

	lr.SetDeadline(time.Now().Add(500 * time.Millisecond))

	atLeastOne := false

	for entry, err := lr.ReadNext(); err == nil; entry, err = lr.ReadNext() {
		lr.SetDeadline(time.Now().Add(500 * time.Millisecond))
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
