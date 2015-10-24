package alog

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func skipIfNoAndroidLoggingFacilities(log LogId, t *testing.T) {
	if _, err := os.Open(filepath.Join("/dev", "alog", log.String())); err != nil {
		t.Skipf("Android logging facilities are not accessible [%s]", err)
	}
}

func readFromLogWorks(log LogId, abiVersion int, t *testing.T) {
	skipIfNoAndroidLoggingFacilities(log, t)

	lr, err := NewLoggerReader(log, abiVersion)
	require.NoError(t, err)

	defer lr.Close()

	lr.SetDeadline(time.Now().Add(500 * time.Millisecond))

	atLeastOne := false

	for entry, err := lr.ReadNext(); err == nil; entry, err = lr.ReadNext() {
		lr.SetDeadline(time.Now().Add(500 * time.Millisecond))
		atLeastOne = true
		t.Logf("%+v\n", entry)
	}

	assert.True(t, atLeastOne, "Log ", log, " is empty")
}

func TestReadFromLogsWorksForAbiV1(t *testing.T) {
	readFromLogWorks(LogIdMain, LoggerAbiV1, t)
	readFromLogWorks(LogIdRadio, LoggerAbiV1, t)
	readFromLogWorks(LogIdEvents, LoggerAbiV1, t)
	readFromLogWorks(LogIdSystem, LoggerAbiV1, t)
}

func TestReadFromLogsWorksForAbiV2(t *testing.T) {
	readFromLogWorks(LogIdMain, LoggerAbiV2, t)
	readFromLogWorks(LogIdRadio, LoggerAbiV2, t)
	readFromLogWorks(LogIdEvents, LoggerAbiV2, t)
	readFromLogWorks(LogIdSystem, LoggerAbiV2, t)
}
