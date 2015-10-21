package alog

// ErrReadTimeout is returned if a call to ReadNext times out.
var ErrReadTimeout = errors.New("Reading the next entry from the log failed")

type Reader interface {
	// SetDeadline adjusts the deadline for all subsequent read operations.
	SetDeadline(deadline time.Time)

	// ReadNext reads the next entry from a Reader.
	//
	// Returns an error if reading the next entry fails.
	// Returns ErrReadTimeout in case of timeouts.
	ReadNext() (*Entry, error)
}
