package alog

import "log"

var (
	Main, _   = NewLoggerWriter(LogIdMain)   // Global Writer for accessing log Main.
	Radio, _  = NewLoggerWriter(LogIdRadio)  // Global Writer for accessing log Radio.
	Events, _ = NewLoggerWriter(LogIdEvents) // Global Writer for accessing log Events.
	System, _ = NewLoggerWriter(LogIdSystem) // Global Writer for accessing log System.
)

type ioWriterWrapper struct {
	prio   Priority
	tag    Tag
	writer Writer
}

func (self ioWriterWrapper) Write(b []byte) (int, error) {
	if err := self.writer.Write(self.prio, self.tag, string(b)); err != nil {
		return 0, err
	}
	return len(b), nil
}

// NewLogger returns a log.Logger instance, sending entries to the Android log
// specificed by logId, with priority prio and tag. logFlags are passed on to the
// log.Logger created by this function, but we disable all time-related flags.
//
// Returns an error if accessing the Android logging facilities fails.
func NewLogger(logId LogId, prio Priority, tag Tag, logFlags int) (*log.Logger, error) {
	// We disable all time-related flags as entries are timestamped by the
	// the Android logging facilities.
	logFlags = logFlags & ^log.Ldate
	logFlags = logFlags & ^log.Ltime
	logFlags = logFlags & ^log.Lmicroseconds

	w, err := NewLoggerWriter(logId)
	if err != nil {
		return nil, err
	}

	iow := &ioWriterWrapper{prio, tag, w}
	return log.New(iow, "", logFlags), nil
}

// V logs message under tag with priority PriorityVerbose to w.
func V(w Writer, tag Tag, message string) error {
	return w.Write(PriorityVerbose, tag, message)
}

// D logs message under tag with priority PriorityDebug to w.
func D(w Writer, tag Tag, message string) error {
	return w.Write(PriorityDebug, tag, message)
}

// I logs message under tag with priority PriorityInfo to w.
func I(w Writer, tag Tag, message string) error {
	return w.Write(PriorityInfo, tag, message)
}

// W logs message under tag with priority PriorityWarning to w.
func W(w Writer, tag Tag, message string) error {
	return w.Write(PriorityWarn, tag, message)
}

// E logs message under tag with priority PriorityError to w.
func E(w Writer, tag Tag, message string) error {
	return w.Write(PriorityError, tag, message)
}

// F logs message under tag with priority PriorityFatal to w.
func F(w Writer, tag Tag, message string) error {
	return w.Write(PriorityFatal, tag, message)
}
