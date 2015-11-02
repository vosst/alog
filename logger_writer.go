package alog

import (
	"os"
	"path/filepath"
	"time"

	"github.com/tedb/vectorio"
)

// A LoggerWriter implements Writer, sending log entries to Android's kernel logger.
type LoggerWriter struct {
	f *os.File // Open File representing our connection to Android's kernel logger.
}

// NewLoggerWriter opens a connection to Android's kernel logger for id,
// returning a LoggerWriter if the operation finishes sucessfully.
//
// Returns an error if connection to the Android kernel logger with id fails.
func NewLoggerWriter(id LogId) (*LoggerWriter, error) {
	fn := filepath.Join("/dev", "alog", id.String())

	f, err := os.OpenFile(fn, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return &LoggerWriter{f: f}, nil
}

// Close shuts down the connection to Android's kernel logger.
func (self *LoggerWriter) Close() error {
	return self.f.Close()
}

// SetDeadline is noop for LoggerWriter. The Android kernel logging
// facilities always report the underlying file as writable. For that,
// polling would be pointless and we can just hand out our write request.
func (self *LoggerWriter) SetDeadline(t time.Time) error {
	return nil
}

// Write sends a log with prio, tag and message to Android's kernel logger.
//
// Returns an error if writing to the kernel logger fails.
func (self *LoggerWriter) Write(prio Priority, tag Tag, message string) error {
	iov := [][]byte{
		[]byte{byte(prio)},
		[]byte(tag),
		[]byte(message),
	}

	// Both the tag and the message need to be null-terminated
	iov[1] = append(iov[1], '\x00')
	iov[2] = append(iov[2], '\x00')

	_, err := vectorio.Writev(self.f, iov)
	return err
}

func (self *LoggerWriter) V(tag Tag, message string) error {
	return self.Write(PriorityVerbose, tag, message)
}

func (self *LoggerWriter) D(tag Tag, message string) error {
	return self.Write(PriorityDebug, tag, message)
}

func (self *LoggerWriter) I(tag Tag, message string) error {
	return self.Write(PriorityInfo, tag, message)
}

func (self *LoggerWriter) W(tag Tag, message string) error {
	return self.Write(PriorityWarn, tag, message)
}

func (self *LoggerWriter) E(tag Tag, message string) error {
	return self.Write(PriorityError, tag, message)
}

func (self *LoggerWriter) F(tag Tag, message string) error {
	return self.Write(PriorityFatal, tag, message)
}
