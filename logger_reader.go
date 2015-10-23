package alog

// #include <sys/ioctl.h>
// #define __LOGGERIO 0xAE
// #define LOGGER_SET_VERSION		_IO(__LOGGERIO, 6) /* abi version */
//
// int RequestLoggerAbiV2(int fd)
// {
//     static int version = 2;
//     return ioctl(fd, LOGGER_SET_VERSION, &version);
// }
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/npat-efault/poller"
)

const (
	maxEntrySize = 5 * 1024 // Max size of a log entry when reading
	LoggerAbiV1  = 1        // Assume logger ABI version 1
	LoggerAbiV2  = 2        // Assume logger ABI version 2
)

// wire bundles all the fields available on the wire.
// Used for parsing a single entry received from the Android logging
// facilities.
type wire struct {
	Len  uint16
	_    uint16
	Pid  int32
	Tid  int32
	Sec  int32
	Nsec int32
}

// A LoggerReader connects to pre-Lollipop kernel logging facilities.
type LoggerReader struct {
	abiVersion int        // ABI version requested from the logging facilities.
	f          *poller.FD // The file we read entries from
	buf        []byte     // Buffer for reading raw bytes from f
}

// NewLoggerReader returns a new LoggerReader reading from the log stream
// identified by id, with the given abiVersion.
//
// Returns an error if accessing the underlying Android log facilities fails.
func NewLoggerReader(id LogId, abiVersion int) (*LoggerReader, error) {
	fn := filepath.Join("/dev", "alog", id.String())

	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	p, err := poller.NewFD(int(f.Fd()))
	if err != nil {
		return nil, err
	}

	if abiVersion == LoggerAbiV2 {
		if res := C.RequestLoggerAbiV2(C.int(f.Fd())); res != 0 {
			return nil, syscall.Errno(res)
		}
	}

	return &LoggerReader{abiVersion: abiVersion, f: p, buf: make([]byte, maxEntrySize, maxEntrySize)}, nil
}

// SetDeadline adjusts the deadline for reading for a LoggerReader.
//
// Returns an error if an issue arises in talking to the underlying
// Android log facilities.
func (self *LoggerReader) SetDeadline(t time.Time) error {
	return self.f.SetReadDeadline(t)
}

// ReadNext reads the next entry from a LaggerReader.
//
// Returns an error if reading from the underlying Android facilities fails,
// specifically if the read operation times out.
func (self *LoggerReader) ReadNext() (*Entry, error) {
	n, err := self.f.Read(self.buf)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(self.buf[:n])
	w := wire{}
	var euid *uint32

	if err = binary.Read(reader, binary.LittleEndian, &w); err != nil {
		return nil, err
	}

	if self.abiVersion == LoggerAbiV2 {
		tmp := uint32(0)
		euid = &tmp
		if err = binary.Read(reader, binary.LittleEndian, euid); err != nil {
			return nil, err
		}
	}

	if buf, err := ioutil.ReadAll(reader); err != nil {
		return nil, err
	} else if len(buf) > 3 { // We need at least a priority, and two \0.
		tagEnd := bytes.IndexAny(buf[1:], "\x00")

		return &Entry{
			Pid:      w.Pid,
			Tid:      w.Tid,
			When:     Timestamp{Seconds: w.Sec, Nanoseconds: w.Nsec},
			Priority: Priority(buf[0]),
			Tag:      Tag(buf[1 : tagEnd+1]),
			Message:  strings.TrimSpace(string(buf[tagEnd+1 : w.Len-1])),
			Euid:     euid,
			Id:       nil,
		}, nil
	} else {
		return nil, errors.New("Invalid log entry")
	}
}
