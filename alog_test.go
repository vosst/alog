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
}
