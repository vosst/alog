package alog

import (
	"errors"
)

// ErrReadTimeout is returned if a call to ReadNext times out.
var ErrReadTimeout = errors.New("Reading the next entry from the log timed out")

type Reader interface {
	// ReadNext reads the next entry from a Reader.
	//
	// Returns an error if reading the next entry fails.
	// Returns ErrReadTimeout in case of timeouts.
	ReadNext() (*Entry, error)
}
