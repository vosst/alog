package quirk

import (
	"encoding/binary"
	"io"
)

// A MeizuMx4LoggerAbiExtension reads the additional timezone field
// as defined in the kernel source code at:
//   https://github.com/meizuosc/m75/blob/master/kernel/drivers/staging/android/logger.h.
type MeizuMx4LoggerAbiExtension struct {
}

// Read reads the tz field from reader, returning it under key "tz".
//
// Returns an error if reading from reader fails.
func (self MeizuMx4LoggerAbiExtension) Read(reader io.Reader) (map[string]interface{}, error) {
	tz := int32(-1)
	if err := binary.Read(reader, binary.LittleEndian, &tz); err != nil {
		return nil, err
	}

	return map[string]interface{}{"tz": tz}, nil
}
