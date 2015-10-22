package alog

// #include <sys/ioctl.h>
// #define __LOGGERIO 0xAE
// #define LOGGER_GET_LOG_BUF_SIZE		_IO(__LOGGERIO, 1) /* size of log */
// #define LOGGER_GET_LOG_LEN		_IO(__LOGGERIO, 2) /* used log len */
// #define LOGGER_GET_NEXT_ENTRY_LEN	_IO(__LOGGERIO, 3) /* next entry len */
// #define LOGGER_FLUSH_LOG		_IO(__LOGGERIO, 4) /* flush log */
// #define LOGGER_GET_VERSION		_IO(__LOGGERIO, 5) /* abi version */
// #define LOGGER_SET_VERSION		_IO(__LOGGERIO, 6) /* abi version */
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/npat-efault/poller"
)

const (
	maxEntrySize = 5 * 1024 // Max size of a log entry when reading
	LoggerAbiV1  = 1        // Assume logger ABI version 1
	LoggerAbiV2  = 2        // Assume logger ABI version 2
)

var (
	loggerIoctlGetLogBufSize   int = C.LOGGER_GET_LOG_BUF_SIZE   // Query the size of the log
	loggerIoctlGetLogLen       int = C.LOGGER_GET_LOG_LEN        // Query the current length of the log
	loggerIoctlGetNextEntryLen int = C.LOGGER_GET_NEXT_ENTRY_LEN // Query the length of the next entry
	loggerIoctlFlushLog        int = C.LOGGER_FLUSH_LOG          // Flush the log
	loggerIoctlGetVersion      int = C.LOGGER_GET_VERSION        // Query the ABI version
	loggerIoctlSetVersion      int = C.LOGGER_SET_VERSION        // Set the ABI version
)

func ioctl(fd, cmd, ptr uintptr) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	if e != 0 {
		return e
	}
	return nil
}

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

	f, err := poller.Open(fn, poller.O_RO)
	if err != nil {
		return nil, err
	}

	if abiVersion == LoggerAbiV2 {
		sfd := f.Sysfd()
		fd := uintptr(unsafe.Pointer(&sfd))
		cmd := uintptr(unsafe.Pointer(&loggerIoctlSetVersion))
		ptr := uintptr(unsafe.Pointer(&abiVersion))

		if err = ioctl(fd, cmd, ptr); err != nil {
			return nil, err
		}
	}

	return &LoggerReader{abiVersion: abiVersion, f: f, buf: make([]byte, maxEntrySize, maxEntrySize)}, nil
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
