# alog [![GoDoc](https://godoc.org/github.com/vosst/alog?status.svg)](https://godoc.org/github.com/vosst/alog)

Package alog enables applications to write to and read from Android's logging
facilities.

## Writing Log Messages

Writing to the Android logging facilities can be accomplished in multiple ways.
For applications leveraging Go's log package, alog.NewLogger is the way to go:

    import "github.com/vosst/alog"

    logger, err := alog.NewLogger(alog.LogIdMain)
    if err != nil {
    	panic(err)
    }
    logger.Print("Successfully connected to Android's logging facilities")

Convenience functions for all log levels are available and applications can
output their tagged messages to Android's well-known logs as in the following
example:

    import "github.com/vosst/alog"

    alog.D(alog.Main, "tag", "message")

Finally, applications can leverage the interface Writer and its implementation
LoggerWriter to write to the Android logging facilities. The respective types
and functions are meant to be used for integration purposes with other logging
frameworks.


## Reading Log Entries

Reading from the Android logging facilities is abstracted by the interface
Reader and its implementation LoggerReader. Applications can access Android's
well known logs and read individual entries as illustrated in the following
snippet:
```Go
package main

import (
	"fmt"
	"time"

	"github.com/vosst/alog"
)

func main() {
	lr, err := alog.NewLoggerReader(alog.LogIdMain)
	lr.SetDeadline(time.Now().Add(500 * time.Millisecond))
	for entry, err := lr.ReadNext(); err == nil; entry, err = lr.ReadNext() {
		fmt.Printf("%s/%s(%5d): %s\n", entry.Priority, entry.Tag, entry.Pid, entry.Message)
		lr.SetDeadline(time.Now().Add(500 * time.Millisecond))
	}
}
```
