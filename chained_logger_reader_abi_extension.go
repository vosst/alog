package alog

import "io"

// A ChainedLoggerAbiExtension is a slice of LoggerAbiExtensions,
// forwarding calls to Prepare and Read to the individual extensions.
type ChainedLoggerAbiExtension struct {
	Extensions []LoggerAbiExtension // All extensions managed by a ChainedLoggerAbiExtension
}

// Read forwards the call to all extensions known to self, merging all results
// into a single extension map.
//
// Returns an error if any of the extensions known to self errors out.
func (self ChainedLoggerAbiExtension) Read(reader io.Reader) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, ext := range self.Extensions {
		if e, err := ext.Read(reader); err != nil {
			return nil, err
		} else {
			for k, v := range e {
				result[k] = v
			}
		}
	}

	return result, nil
}
