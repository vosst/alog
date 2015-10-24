package alog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var testTag = Tag("Test")

func testLoggerWriterWorks(logId LogId, t *testing.T) {
	skipIfNoAndroidLoggingFacilities(logId, t)

	writer, err := NewLoggerWriter(logId)
	require.NoError(t, err)

	defer writer.Close()

	writer.Write(PriorityDebug, testTag, "42")

	reader, err := NewLoggerReader(logId, LoggerAbiV1)
	require.NoError(t, err)

	defer reader.Close()

	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, err := reader.ReadNext()

	require.NoError(t, err)

	assert.Equal(t, testTag, entry.Tag)
	assert.Equal(t, PriorityDebug, entry.Priority)
	assert.Equal(t, "42", entry.Message)
}

func TestLoggerWriteWorks(t *testing.T) {
	testLoggerWriterWorks(LogIdMain, t)
	testLoggerWriterWorks(LogIdRadio, t)
	testLoggerWriterWorks(LogIdEvents, t)
	testLoggerWriterWorks(LogIdSystem, t)
}
