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
	"io"
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

// requestExtendedLoggerAbi issues an ioctl on fd to request AOSP Logger wire format v2.
//
// Returns an error if the ioctl on fd fails.
func requestExtendedLoggerAbi(fd int) error {
	if rc := C.RequestLoggerAbiV2(C.int(fd)); rc != 0 {
		return syscall.Errno(rc)
	}

	return nil
}

// A LoggerAbiExtension models additions to the logger v1 wire format defined by
// AOSP.
type LoggerAbiExtension interface {
	// Read allows implementations to unmarshal addition fields from a reader
	// into the buffer returned by a single read call to Android's logger facilities.
	Read(reader io.Reader) (map[string]interface{}, error)
}

// A LoggerAbiV2Extension implements LoggerAbiExtension, reading the
// additional euid field.
type LoggerAbiV2Extension struct {
}

// Read unmarshals the additional euid field from reader, and return it in the extension map under key 'euid'.
//
// Returns an error if unmarshaling the euid field from reader fails.
func (self LoggerAbiV2Extension) Read(reader io.Reader) (map[string]interface{}, error) {
	euid := uint32(0)
	if err := binary.Read(reader, binary.LittleEndian, &euid); err != nil {
		return nil, err
	}
	return map[string]interface{}{"euid": euid}, nil
}

// A LoggerReader connects to pre-Lollipop kernel logging facilities.
type LoggerReader struct {
	abiExtension LoggerAbiExtension // ABI extension handler
	f            *poller.FD         // The file we read entries from
	buf          []byte             // Buffer for reading raw bytes from f
}

// NewLoggerReader returns a new LoggerReader reading from the log stream
// identified by id. If abiExtension is not nil it is used to parse additional fields from
// a buffer read from Android's logging facilities.
//
// Returns an error if accessing the underlying Android log facilities fails.
func NewLoggerReader(id LogId, abiExtension LoggerAbiExtension) (*LoggerReader, error) {
	fn := filepath.Join("/dev", "alog", id.String())

	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	if abiExtension != nil {
		if err = requestExtendedLoggerAbi(int(f.Fd())); err != nil {
			return nil, err
		}
	}

	p, err := poller.NewFD(int(f.Fd()))
	if err != nil {
		return nil, err
	}

	return &LoggerReader{abiExtension: abiExtension, f: p, buf: make([]byte, maxEntrySize, maxEntrySize)}, nil
}

// Close() closes the underlying connection to the Android logger facilities.
//
// Outstanding ReadNext operations are cancelled and return an error.
func (self *LoggerReader) Close() error {
	return self.f.Close()
}

// SetDeadline adjusts the deadline for reading for a LoggerReader.
//
// Returns an error if an issue arises in talking to the underlying
// Android log facilities.
func (self *LoggerReader) SetDeadline(t time.Time) error {
	return self.f.SetReadDeadline(t)
}

// ReadNext reads the next entry from a LaggerReader. Extension fields (if any)
// are placed into the Ext field of Entry.
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

	if err = binary.Read(reader, binary.LittleEndian, &w); err != nil {
		return nil, err
	}

	var ext map[string]interface{}

	if self.abiExtension != nil {
		ext, err = self.abiExtension.Read(reader)
		if err != nil {
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
			Ext:      ext,
		}, nil
	} else {
		return nil, errors.New("Invalid log entry")
	}
}
