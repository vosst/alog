package alog

import (
	"fmt"
	"time"
)

func ExampleLoggerReader() {
	lr, err := NewLoggerReader(LogIdMain, nil)
	if err != nil {
		panic(err)
	}

	// Keep on reading for the next 5 seconds.
	lr.SetDeadline(time.Now().Add(5 * time.Second))

	// Read all available entries from log until an error occurs.
	// Timing out would return an error here.
	for le, err := lr.ReadNext(); err == nil; le, err = lr.ReadNext() {
		fmt.Printf("%s/%s(%d): %s\n", le.Priority, le.Tag, le.Pid, le.Message)
	}
}
