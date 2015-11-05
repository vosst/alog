# alog
--
    import "github.com/vosst/alog"

Package alog enables applications to write to and read from Android's logging
facilities.


### Writing Log Messages

Writing to the Android logging facilities can be accomplished in multiple ways.
For applications leveraging Go's log package, alog.NewLogger is the way to go:

    logger, err := alog.NewLogger(alog.LogIdMain)
    if err != nil {
    	panic(err)
    }
    logger.Print("Successfully connected to Android's logging facilities")

Convenience functions for all log levels are available and applications can
output their tagged messages to Android's well-known logs as in the following
example:

    alog.D(alog.Main, "tag", "message")

Finally, applications can leverage the interface Writer and its implementation
LoggerWriter to write to the Android logging facilities. The respective types
and functions are meant to be used for integration purposes with other logging
frameworks.


### Reading Log Entries

Reading from the Android logging facilities is abstracted by the interface
Reader and its implementation LoggerReader. Applications can access Android's
well known logs and read individual entries as illustrated in the following
snippet:

    lr, err := alog.NewLoggerReader(alog.LogIdMain)
    lr.SetDeadline(time.Now().Add(500 * time.Millisecond))
    for entry, err := lr.ReadNext(); err == nil; entry, err = lr.ReadNext() {
    	lr.SetDeadline(time.Now().Add(500 * time.Millisecond))
    }

## Usage

```go
var (
	Main   = NewLoggerWriter(LogIdMain)   // Global Writer for accessing log Main.
	Radio  = NewLoggerWriter(LogIdRadio)  // Global Writer for accessing log Radio.
	Events = NewLoggerWriter(LogIdEvents) // Global Writer for accessing log Events.
	System = NewLoggerWriter(LogIdSystem) // Global Writer for accessing log System.
)
```

```go
var ErrReadTimeout = errors.New("Reading the next entry from the log timed out")
```
ErrReadTimeout is returned if a call to ReadNext times out.

#### func  D

```go
func D(w Writer, tag Tag, message string) error
```
D logs message under tag with priority PriorityDebug to w.

#### func  E

```go
func E(w Writer, tag Tag, message string) error
```
E logs message under tag with priority PriorityError to w.

#### func  F

```go
func F(w Writer, tag Tag, message string) error
```
F logs message under tag with priority PriorityFatal to w.

#### func  I

```go
func I(w Writer, tag Tag, message string) error
```
I logs message under tag with priority PriorityInfo to w.

#### func  NewLogger

```go
func NewLogger(logId LogId, prio Priority, tag Tag, logFlags int) (*log.Logger, error)
```
NewLogger returns a log.Logger instance, sending entries to the Android log
specificed by logId, with priority prio and tag. logFlags are passed on to the
log.Logger created by this function, but we disable all time-related flags.

Returns an error if accessing the Android logging facilities fails.

#### func  V

```go
func V(w Writer, tag Tag, message string) error
```
V logs message under tag with priority PriorityVerbose to w.

#### func  W

```go
func W(w Writer, tag Tag, message string) error
```
W logs message under tag with priority PriorityWarning to w.

#### type ChainedLoggerAbiExtension

```go
type ChainedLoggerAbiExtension struct {
	Extensions []LoggerAbiExtension // All extensions managed by a ChainedLoggerAbiExtension
}
```

A ChainedLoggerAbiExtension is a slice of LoggerAbiExtensions, forwarding calls
to Prepare and Read to the individual extensions.

#### func (ChainedLoggerAbiExtension) Read

```go
func (self ChainedLoggerAbiExtension) Read(reader io.Reader) (map[string]interface{}, error)
```
Read forwards the call to all extensions known to self, merging all results into
a single extension map.

Returns an error if any of the extensions known to self errors out.

#### type Entry

```go
type Entry struct {
	Pid      int32                  // Generating process's ID
	Tid      int32                  // Generating thread's ID
	When     Timestamp              // When the entry was logged
	Priority Priority               // Priority of the message
	Tag      Tag                    // Tag describing the origin of the message
	Message  string                 // The actual message of the Entry
	Ext      map[string]interface{} // Vendor specific extensions to individual entries
}
```

An Entry models an individual log message.

#### type LogId

```go
type LogId int
```

A LogId uniquely names a log stream.

```go
const (
	LogIdMain   LogId = 0
	LogIdRadio  LogId = 1
	LogIdEvents LogId = 2
	LogIdSystem LogId = 3
	LogIdCrash  LogId = 4
)
```

#### func (LogId) String

```go
func (self LogId) String() string
```
String returns the name of a LogId

#### type LoggerAbiExtension

```go
type LoggerAbiExtension interface {
	// Read allows implementations to unmarshal addition fields from a reader
	// into the buffer returned by a single read call to Android's logger facilities.
	Read(reader io.Reader) (map[string]interface{}, error)
}
```

A LoggerAbiExtension models additions to the logger v1 wire format defined by
AOSP.

#### type LoggerAbiV2Extension

```go
type LoggerAbiV2Extension struct {
}
```

A LoggerAbiV2Extension implements LoggerAbiExtension, reading the additional
euid field.

#### func (LoggerAbiV2Extension) Read

```go
func (self LoggerAbiV2Extension) Read(reader io.Reader) (map[string]interface{}, error)
```
Read unmarshals the additional euid field from reader, and return it in the
extension map under key 'euid'.

Returns an error if unmarshaling the euid field from reader fails.

#### type LoggerReader

```go
type LoggerReader struct {
}
```

A LoggerReader connects to pre-Lollipop kernel logging facilities.

#### func  NewLoggerReader

```go
func NewLoggerReader(id LogId, abiExtension LoggerAbiExtension) (*LoggerReader, error)
```
NewLoggerReader returns a new LoggerReader reading from the log stream
identified by id. If abiExtension is not nil it is used to parse additional
fields from a buffer read from Android's logging facilities.

Returns an error if accessing the underlying Android log facilities fails.

#### func (*LoggerReader) Close

```go
func (self *LoggerReader) Close() error
```
Close() closes the underlying connection to the Android logger facilities.

Outstanding ReadNext operations are cancelled and return an error.

#### func (*LoggerReader) ReadNext

```go
func (self *LoggerReader) ReadNext() (*Entry, error)
```
ReadNext reads the next entry from a LaggerReader. Extension fields (if any) are
placed into the Ext field of Entry.

Returns an error if reading from the underlying Android facilities fails,
specifically if the read operation times out.

#### func (*LoggerReader) SetDeadline

```go
func (self *LoggerReader) SetDeadline(t time.Time) error
```
SetDeadline adjusts the deadline for reading for a LoggerReader.

Returns an error if an issue arises in talking to the underlying Android log
facilities.

#### type LoggerWriter

```go
type LoggerWriter struct {
}
```

A LoggerWriter implements Writer, sending log entries to Android's kernel
logger.

#### func  NewLoggerWriter

```go
func NewLoggerWriter(id LogId) (*LoggerWriter, error)
```
NewLoggerWriter opens a connection to Android's kernel logger for id, returning
a LoggerWriter if the operation finishes sucessfully.

Returns an error if connection to the Android kernel logger with id fails.

#### func (*LoggerWriter) Close

```go
func (self *LoggerWriter) Close() error
```
Close shuts down the connection to Android's kernel logger.

#### func (*LoggerWriter) D

```go
func (self *LoggerWriter) D(tag Tag, message string) error
```

#### func (*LoggerWriter) E

```go
func (self *LoggerWriter) E(tag Tag, message string) error
```

#### func (*LoggerWriter) F

```go
func (self *LoggerWriter) F(tag Tag, message string) error
```

#### func (*LoggerWriter) I

```go
func (self *LoggerWriter) I(tag Tag, message string) error
```

#### func (*LoggerWriter) SetDeadline

```go
func (self *LoggerWriter) SetDeadline(t time.Time) error
```
SetDeadline is noop for LoggerWriter. The Android kernel logging facilities
always report the underlying file as writable. For that, polling would be
pointless and we can just hand out our write request.

#### func (*LoggerWriter) V

```go
func (self *LoggerWriter) V(tag Tag, message string) error
```

#### func (*LoggerWriter) W

```go
func (self *LoggerWriter) W(tag Tag, message string) error
```

#### func (*LoggerWriter) Write

```go
func (self *LoggerWriter) Write(prio Priority, tag Tag, message string) error
```
Write sends a log with prio, tag and message to Android's kernel logger.

Returns an error if writing to the kernel logger fails.

#### type Priority

```go
type Priority int
```

A Priority models the log priority of a single Entry.

```go
const (
	PriorityUnknown Priority = 0
	PriorityDefault Priority = 1
	PriorityVerbose Priority = 2
	PriorityDebug   Priority = 3
	PriorityInfo    Priority = 4
	PriorityWarn    Priority = 5
	PriorityError   Priority = 6
	PriorityFatal   Priority = 7
	PrioritySilent  Priority = 8
)
```

#### func (Priority) String

```go
func (self Priority) String() string
```
String returns a Priority as a string (short code).

#### type Reader

```go
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
```

A Reader provides means to read Entries from Android's log facilities.

#### type Tag

```go
type Tag string
```

A Tag describes the origin of an Entry.

#### type Timestamp

```go
type Timestamp struct {
	Seconds     int32 // Seconds since the epoch
	Nanoseconds int32 // Nanoseconds since the epoch
}
```

A Timestamp marks the time when an entry was put to a log.

#### type Writer

```go
type Writer interface {
	// A Writer needs to be closed explicitly
	io.Closer

	// SetDeadline adjusts the deadline such that all subsequent calls to
	// ReadNext will fail if they exceed t.
	SetDeadline(t time.Time) error

	// Write logs an message with priority prio and a tag.
	//
	// Returns an error if writing to the underlying Android logging facilities fails.
	Write(prio Priority, tag Tag, message string) error
}
```

A Writer allows for logging to Android's logging facilities.
