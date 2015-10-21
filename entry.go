package alog

// A Timestamp marks the time when an entry was put to a log.
type Timestamp struct {
	Seconds     int32 // Seconds since the epoch
	Nanoseconds int32 // Nanoseconds since the epoch
}

// A Tag describes the origin of an Entry.
type Tag string

// A Priority models the log priority of a single Entry.
type Priority int

// String returns a Priority as a string (short code).
func (self Priority) String() string {
	switch self {
	default:
		return "U"
	case PriorityUnknown:
		return "U"
	case PriorityDefault:
		return "D"
	case PriorityVerbose:
		return "V"
	case PriorityDebug:
		return "D"
	case PriorityInfo:
		return "I"
	case PriorityWarn:
		return "W"
	case PriorityError:
		return "E"
	case PriorityFatal:
		return "F"
	case PrioritySilent:
		return "S"
	}
}

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

// An Entry models an individual log message.
type Entry struct {
	Pid      int32     // Generating process's ID
	Tid      int32     // Generating thread's ID
	When     Timestamp // When the entry was logged
	Priority Priority  // Priority of the message
	Tag      Tag       // Tag describing the origin of the message
	Message  string    // The actual message of the Entry
	Euid     *uint32   // Effective user ID of the logger, may be nil
	Id       *LogId    // Id of the log that the entry comes from, may be nil
}
