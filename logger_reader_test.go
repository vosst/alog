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

func readFromLogWorks(log LogId, abiExtension LoggerAbiExtension, t *testing.T) {
	skipIfNoAndroidLoggingFacilities(log, t)

	lr, err := NewLoggerReader(log, abiExtension)
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
	readFromLogWorks(LogIdMain, nil, t)
	readFromLogWorks(LogIdRadio, nil, t)
	readFromLogWorks(LogIdEvents, nil, t)
	readFromLogWorks(LogIdSystem, nil, t)
}

func TestReadFromLogsWorksForAbiV2(t *testing.T) {
	readFromLogWorks(LogIdMain, LoggerAbiV2Extension{}, t)
	readFromLogWorks(LogIdRadio, LoggerAbiV2Extension{}, t)
	readFromLogWorks(LogIdEvents, LoggerAbiV2Extension{}, t)
	readFromLogWorks(LogIdSystem, LoggerAbiV2Extension{}, t)
}
