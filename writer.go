package alog

import (
	"io"
)

// A Writer allows for logging to Android's logging facilities.
type Writer interface {
	// A Writer needs to be closed explicitly
	io.Closer

	// Write logs an message with priority prio and a tag.
	//
	// Returns an error if writing to the underlying Android logging facilities fails.
	Write(prio Priority, tag Tag, message string) error
}
