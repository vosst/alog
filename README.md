# alog [![GoDoc](https://godoc.org/github.com/vosst/alog?status.svg)](https://godoc.org/github.com/vosst/alog)

Package alog enables applications to write to and read from Android's logging
facilities.

## Writing Log Messages

Writing to the Android logging facilities can be accomplished in multiple ways.
For applications leveraging Go's log package, alog.NewLogger is the way to go:
```Go
import "github.com/vosst/alog"

logger, err := alog.NewLogger(alog.LogIdMain)
if err != nil {
	panic(err)
}

logger.Print("Successfully connected to Android's logging facilities")
```
Convenience functions for all log levels are available and applications can
output their tagged messages to Android's well-known logs as in the following
example:
```Go
import "github.com/vosst/alog"

alog.D(alog.Main, "tag", "message")
```
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
import (
	"fmt"
	"time"

	"github.com/vosst/alog"
)

lr, err := alog.NewLoggerReader(alog.LogIdMain, nil)
if err != nil {
	panic(err)
}

defer lr.Close()
lr.SetDeadline(time.Now().Add(500 * time.Millisecond))

// Loop over all entries in the log, waiting at most 500ms per iteration
// for a new log entry to arrive. Otherwise, we time out, an error is returned
// and the loop is terminated.
for entry, err := lr.ReadNext(); err == nil; entry, err = lr.ReadNext() {
	fmt.Printf("%s/%s(%5d): %s\n", entry.Priority, entry.Tag, entry.Pid, entry.Message)
	lr.SetDeadline(time.Now().Add(500 * time.Millisecond))
}
```
### A Tale of >= 2 ABIs

Android's kernel logging facilities as available until Lollipop support two different ABIs (see https://android.googlesource.com/platform/system/core/+/android-4.4.4_r2.0.1/include/log/logger.h), with the main difference being an additional member `euid` per log entry. In addition, different SOCs have come up with all sorts of interesting variations of the version 2 ABI. Package alog supports all of them and is easily extensible to account for specific customizations. Applications can enable the v2 ABI by passing in a non-nil implementation of `alog.LoggerAbiExtension` to alog.NewLoggerReader as in:
```Go
import (
	"fmt"
	"time"

	"github.com/vosst/alog"
)

lr, err := alog.NewLoggerReader(alog.LogIdMain, alog.LoggerAbiV2Extension{})
if err != nil {
	panic(err)
}

defer lr.Close()
lr.SetDeadline(time.Now().Add(500 * time.Millisecond))

for entry, err := lr.ReadNext(); err == nil; entry, err = lr.ReadNext() {
	fmt.Printf("%s/%s(%5d)@%d: %s\n", entry.Priority, entry.Tag, entry.Pid, entry.Message, entry.Ext["euid"])
	lr.SetDeadline(time.Now().Add(500 * time.Millisecond))
}
```

SOC-specific quirks are supported by chaining together multiple implementations of
`alog.LoggerAbiExtesion` and passing the respective chain into the NewLoggerReader call:
```Go
import (
	"fmt"
	"time"

	"github.com/vosst/alog"
	"github.com/vosst/alog/quirk"
)

chain := alog.ChainedLoggerAbiExtension{
	Extensions: []alog.LoggerAbiExtesion{quirk.MeizuMx4LoggerAbiExtension{}, alog.LoggerAbiV2Extension{}},
}

lr, err := alog.NewLoggerReader(alog.LogIdMain, chain)
if err != nil {
	panic(err)
}

defer lr.Close()
lr.SetDeadline(time.Now().Add(500 * time.Millisecond))

for entry, err := lr.ReadNext(); err == nil; entry, err = lr.ReadNext() {
	fmt.Printf("%s/%s(%5d)@%d|%d: %s\n", entry.Priority, entry.Tag, entry.Pid, entry.Message, entry.Ext["euid"], entry.Ext["tz"])
	lr.SetDeadline(time.Now().Add(500 * time.Millisecond))
}

```
