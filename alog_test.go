package alog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testLogPrintWorks(t *testing.T, logId LogId, prio Priority, tag Tag) {
	skipIfNoAndroidLoggingFacilities(logId, t)

	reader, err := NewLoggerReader(logId, nil)
	require.NoError(t, err)

	defer reader.Close()

	drainLog(reader)

	l, err := NewLogger(logId, prio, tag, 0)
	require.NoError(t, err)

	l.Print("42")

	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, err := reader.ReadNext()
	require.NoError(t, err)

	assert.Equal(t, prio, entry.Priority)
	assert.Equal(t, tag, entry.Tag)
}

func TestLogPrintWorks(t *testing.T) {
	testLogPrintWorks(t, LogIdMain, PriorityDebug, testTag)
	testLogPrintWorks(t, LogIdMain, PriorityInfo, testTag)
	testLogPrintWorks(t, LogIdMain, PriorityWarn, testTag)
	testLogPrintWorks(t, LogIdMain, PriorityError, testTag)

	testLogPrintWorks(t, LogIdRadio, PriorityDebug, testTag)
	testLogPrintWorks(t, LogIdRadio, PriorityInfo, testTag)
	testLogPrintWorks(t, LogIdRadio, PriorityWarn, testTag)
	testLogPrintWorks(t, LogIdRadio, PriorityError, testTag)

	testLogPrintWorks(t, LogIdEvents, PriorityDebug, testTag)
	testLogPrintWorks(t, LogIdEvents, PriorityInfo, testTag)
	testLogPrintWorks(t, LogIdEvents, PriorityWarn, testTag)
	testLogPrintWorks(t, LogIdEvents, PriorityError, testTag)

	testLogPrintWorks(t, LogIdSystem, PriorityDebug, testTag)
	testLogPrintWorks(t, LogIdSystem, PriorityInfo, testTag)
	testLogPrintWorks(t, LogIdSystem, PriorityWarn, testTag)
	testLogPrintWorks(t, LogIdSystem, PriorityError, testTag)
}

func ExampleNewLogger() {
	l, err := NewLogger(LogIdMain, PriorityInfo, "MyTag", 0)
	if err != nil {
		panic(err)
	}

	l.Print("Test")
}
