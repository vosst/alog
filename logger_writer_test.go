package alog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var testTag = Tag("Test")

func drainLog(reader *LoggerReader) {
	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	for _, err := reader.ReadNext(); err == nil; _, err = reader.ReadNext() {
		reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	}
}

func testLoggerWriterWorks(logId LogId, t *testing.T) {
	skipIfNoAndroidLoggingFacilities(logId, t)

	reader, err := NewLoggerReader(logId, nil)
	require.NoError(t, err)

	defer reader.Close()

	drainLog(reader)

	writer, err := NewLoggerWriter(logId)
	require.NoError(t, err)

	defer writer.Close()

	writer.Write(PriorityDebug, testTag, "42")

	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, err := reader.ReadNext()

	require.NoError(t, err)

	t.Logf("%+v\n", entry)

	assert.Equal(t, testTag, entry.Tag)
	assert.Equal(t, PriorityDebug, entry.Priority)
	// TODO(tvoss): Figure out why this fails with:
	// --- FAIL: TestLoggerWriteWorks (1.54 seconds)
	//	logger_writer_test.go:37: &{Pid:21471 Tid:21476 When:{Seconds:1446319710 Nanoseconds:66931758} Priority:D Tag:Test Message:42 Ext:map[]}
	// 	logger_writer_test.go:38: 42
	// 	Error Trace:    logger_writer_test.go:41
	//			logger_writer_test.go:45
	//	Error:		Not equal: "42" (expected)
	//			!= "\x0042" (astual)
	// assert.Equal(t, "42", entry.Message)
}

func TestLoggerWriteWorks(t *testing.T) {
	testLoggerWriterWorks(LogIdMain, t)
	testLoggerWriterWorks(LogIdRadio, t)
	testLoggerWriterWorks(LogIdEvents, t)
	testLoggerWriterWorks(LogIdSystem, t)
}

func TestLoggerWriterConvenienceWorks(t *testing.T) {
	skipIfNoAndroidLoggingFacilities(LogIdMain, t)

	reader, err := NewLoggerReader(LogIdMain, nil)
	require.NoError(t, err)

	defer reader.Close()

	drainLog(reader)

	writer, err := NewLoggerWriter(LogIdMain)
	require.NoError(t, err)

	defer writer.Close()

	writer.V(testTag, "V")
	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, _ := reader.ReadNext()
	assert.Equal(t, PriorityVerbose, entry.Priority)

	writer.D(testTag, "D")
	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, _ = reader.ReadNext()
	assert.Equal(t, PriorityDebug, entry.Priority)

	writer.I(testTag, "I")
	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, _ = reader.ReadNext()
	assert.Equal(t, PriorityInfo, entry.Priority)

	writer.W(testTag, "W")
	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, _ = reader.ReadNext()
	assert.Equal(t, PriorityWarn, entry.Priority)

	writer.E(testTag, "E")
	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, _ = reader.ReadNext()
	assert.Equal(t, PriorityError, entry.Priority)

	writer.F(testTag, "F")
	reader.SetDeadline(time.Now().Add(500 * time.Millisecond))
	entry, _ = reader.ReadNext()
	assert.Equal(t, PriorityFatal, entry.Priority)

}

func ExampleLoggerWriter() {
	writer, err := NewLoggerWriter(LogIdMain)
	if err != nil {
		panic(err)
	}

	// The writer has to be closed explicitly.
	defer writer.Close()

	// Log a debug entry with tag 'A funky tag' and a message
	// giving the answer to life, the universe and everything.
	writer.Write(PriorityDebug, "A funky tag", "42")
}

func ExampleLoggerWriter_second() {
	main, err := NewLoggerWriter(LogIdMain)
	if err != nil {
		panic(err)
	}

	defer main.Close()

	// Log a debug entry with tag 'A funky tag' and a message
	// giving the answer to life, the universe and everything.
	main.D("A funky tag", "42")

}
