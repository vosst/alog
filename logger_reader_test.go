package alog

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vosst/alog/quirk"
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
	cext := ChainedLoggerAbiExtension{
		Extensions: []LoggerAbiExtension{&quirk.MeizuMx4LoggerAbiExtension{}, &LoggerAbiV2Extension{}},
	}
	readFromLogWorks(LogIdMain, cext, t)
	readFromLogWorks(LogIdRadio, cext, t)
	readFromLogWorks(LogIdEvents, cext, t)
	readFromLogWorks(LogIdSystem, cext, t)
}

func TestLoggerReaderCallsNonNilAbiExtension(t *testing.T) {
	skipIfNoAndroidLoggingFacilities(LogIdMain, t)

	m := make(map[string]interface{})
	mae := &MockLoggerAbiExtension{}

	mae.On("Read", mock.Anything).Return(m, nil)

	lr, err := NewLoggerReader(LogIdMain, mae)
	require.NoError(t, err)

	defer lr.Close()

	_, err = lr.ReadNext()
	assert.NoError(t, err)

	mae.AssertNumberOfCalls(t, "Read", 1)
}

func ExampleLoggerReader() {
	lr, err := NewLoggerReader(LogIdMain, nil)
	if err != nil {
		panic(err)
	}

	// Keep on reading for the next 5 seconds.
	lr.SetDeadline(time.Now().Add(5 * time.Second))

	// Read all available entries from log until an error occurs.
	// Timing out would return an error here.
	for le, err := lr.ReadNext(); err == nil; le, err = lr.ReadNext() {
		fmt.Printf("%s/%s(%d): %s\n", le.Priority, le.Tag, le.Pid, le.Message)
	}
}
