package alog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testLogPrintWorks(t *testing.T) {
	skipIfNoAndroidLoggingFacilities(LogIdMain, t)

	reader, err := NewLoggerReader(LogIdMain, nil)
	require.NoError(t, err)

	defer reader.Close()

	drainLog(reader)

	l, err := NewLogger(LogIdMain, PriorityDebug, testTag, 0)
	require.NoError(t, err)

	l.Print("42")

	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, err := reader.ReadNext()
	require.NoError(t, err)

	assert.Equal(t, PriorityDebug, entry.Priority)
	assert.Equal(t, testTag, entry.Tag)
}
