package alog

import "log"

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
