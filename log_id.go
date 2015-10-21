package alog

// A LogId uniquely names a log stream.
type LogId int

const (
	LogIdMain   LogId = 0
	LogIdRadio        = 1
	LogIdEvents       = 2
	LogIdSystem       = 3
	LogIdCrash        = 4
)

// String returns the name of a LogId
func (self LogId) String() string {
	switch self {
	case LogIdMain:
		return "main"
	case LogIdRadio:
		return "radio"
	case LogIdEvents:
		return "events"
	case LogIdSystem:
		return "system"
	case LogIdCrash:
		return "crash"
	default:
		return "main"
	}
}
