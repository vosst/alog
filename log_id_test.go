package alog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogIdConstantValues(t *testing.T) {
	assert.EqualValues(t, 0, LogIdMain)
	assert.EqualValues(t, 1, LogIdRadio)
	assert.EqualValues(t, 2, LogIdEvents)
	assert.EqualValues(t, 3, LogIdSystem)
	assert.EqualValues(t, 4, LogIdCrash)
}

func TestLogIdStringReturnsCorrectValues(t *testing.T) {
	assert.Equal(t, "main", LogIdMain.String())
	assert.Equal(t, "radio", LogIdRadio.String())
	assert.Equal(t, "events", LogIdEvents.String())
	assert.Equal(t, "system", LogIdSystem.String())
	assert.Equal(t, "crash", LogIdCrash.String())

	assert.Equal(t, "main", LogId(42).String())
}
