package alog

// A Timestamp marks the time when an entry was put to a log.
type Timestamp struct {
	Seconds     int32 // Seconds since the epoch
	Nanoseconds int32 // Nanoseconds since the epoch
}

// An Entry models an individual log message.
type Entry struct {
	Pid     int32     // Generating process's ID
	Tid     int32     // Generating thread's ID
	When    Timestamp // When the entry was logged
	Message string    // The Entry's payload
	Euid    *uint32   // Effective user ID of the logger, may be nil.
	Id      *LogId    // Id of the log that the entry comes from, may be nil.
}
