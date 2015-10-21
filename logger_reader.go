package alog

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

// Max size of a log entry when reading
const maxEntrySize = 5 * 1024

// A LoggerReader connects to pre-Lollipop kernel logging facilities.
type LoggerReader struct {
	guard sync.Mutex // Guards f and buf
	f     *os.File   // The file we read entries from
	buf   []byte     // Buffer for reading raw bytes from f
}

func NewLoggerReader(id LogId) (*LoggerReader, error) {
	f, err := os.OpenFile(filepath.Join("/dev", "alog", id.String()), syscall.O_NONBLOCK|syscall.O_RDONLY, os.FileMode(0))
	if err != nil {
		return nil, err
	}

	return &LoggerReader{f: f, buf: make([]byte, maxEntrySize, maxEntrySize)}, nil
}

func (self LoggerReader) ReadNext() (*Entry, error) {
	self.guard.Lock()
	defer self.guard.Unlock()

	_, err := self.f.Read(self.buf)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(self.buf)

	type Wire struct {
		Length uint16
		_      uint16
		Pid    int32
		Tid    int32
		Sec    int32
		Nsec   int32
	}

	wire := Wire{}

	if err = binary.Read(reader, binary.LittleEndian, &wire); err != nil {
		return nil, err
	}

	if buf, err := ioutil.ReadAll(reader); err != nil {
		return nil, err
	} else if len(buf) > 3 { // We need at least a priority, and two \0.
		tagEnd := bytes.IndexAny(buf[1:], "\x00")

		return &Entry{
			Pid:      wire.Pid,
			Tid:      wire.Tid,
			When:     Timestamp{Seconds: wire.Sec, Nanoseconds: wire.Nsec},
			Priority: Priority(buf[0]),
			Tag:      Tag(buf[1 : tagEnd+1]),
			Message:  string(buf[tagEnd+1:]),
			Euid:     nil,
			Id:       nil,
		}, nil
	} else {
		return nil, errors.New("Invalid log entry")
	}
}
