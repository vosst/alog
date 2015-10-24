package alog

import (
	"errors"
	"io"
	"time"
)

// ErrReadTimeout is returned if a call to ReadNext times out.
var ErrReadTimeout = errors.New("Reading the next entry from the log timed out")

// A Reader provides means to read Entries from Android's log facilities.
type Reader interface {
	// A Reader has to be closed explicitly.
	io.Closer
	// SetDeadline adjusts the deadline such that all subsequent calls to
	// ReadNext will fail if they exceed t.
	SetDeadline(t time.Time) error

	// ReadNext reads the next entry from a Reader
	//
	// Returns an error if reading the next entry fails.
	// Returns ErrReadTimeout in case of timeouts.
	ReadNext() (*Entry, error)
}
