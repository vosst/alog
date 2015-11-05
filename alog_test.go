package alog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type MockWriter struct {
	mock.Mock
}

func (self *MockWriter) SetDeadline(t time.Time) error {
	return self.Called(t).Error(0)
}

func (self *MockWriter) Close() error {
	return self.Called().Error(0)
}

func (self *MockWriter) Write(prio Priority, tag Tag, msg string) error {
	args := self.Called(prio, tag, msg)
	return args.Error(0)
}

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

func TestConvenienceFunctionsCallIntoWriter(t *testing.T) {
	mw := &MockWriter{}

	mw.On("Write", PriorityVerbose, testTag, "42").Return(nil)
	V(mw, testTag, "42")
	mw.AssertExpectations(t)

	mw.On("Write", PriorityDebug, testTag, "42").Return(nil)
	D(mw, testTag, "42")
	mw.AssertExpectations(t)

	mw.On("Write", PriorityInfo, testTag, "42").Return(nil)
	I(mw, testTag, "42")
	mw.AssertExpectations(t)

	mw.On("Write", PriorityWarn, testTag, "42").Return(nil)
	W(mw, testTag, "42")
	mw.AssertExpectations(t)

	mw.On("Write", PriorityError, testTag, "42").Return(nil)
	E(mw, testTag, "42")
	mw.AssertExpectations(t)
}

func TestConvenienceFunctionsWorkWithGlobalInstances(t *testing.T) {
	skipIfNoAndroidLoggingFacilities(LogIdMain, t)
	skipIfNoAndroidLoggingFacilities(LogIdRadio, t)
	skipIfNoAndroidLoggingFacilities(LogIdEvents, t)
	skipIfNoAndroidLoggingFacilities(LogIdSystem, t)

	assert.NoError(t, V(Main, testTag, "42"))
	assert.NoError(t, V(Radio, testTag, "42"))
	assert.NoError(t, V(Events, testTag, "42"))
	assert.NoError(t, V(System, testTag, "42"))

	assert.NoError(t, D(Main, testTag, "42"))
	assert.NoError(t, D(Radio, testTag, "42"))
	assert.NoError(t, D(Events, testTag, "42"))
	assert.NoError(t, D(System, testTag, "42"))

	assert.NoError(t, I(Main, testTag, "42"))
	assert.NoError(t, I(Radio, testTag, "42"))
	assert.NoError(t, I(Events, testTag, "42"))
	assert.NoError(t, I(System, testTag, "42"))

	assert.NoError(t, W(Main, testTag, "42"))
	assert.NoError(t, W(Radio, testTag, "42"))
	assert.NoError(t, W(Events, testTag, "42"))
	assert.NoError(t, W(System, testTag, "42"))

	assert.NoError(t, E(Main, testTag, "42"))
	assert.NoError(t, E(Radio, testTag, "42"))
	assert.NoError(t, E(Events, testTag, "42"))
	assert.NoError(t, E(System, testTag, "42"))
}

func ExampleNewLogger() {
	l, err := NewLogger(LogIdMain, PriorityInfo, "MyTag", 0)
	if err != nil {
		panic(err)
	}

	l.Print("Test")
}
