package alog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFromMainWorks(t *testing.T) {
	if _, err := os.Open("/dev/alog/main"); err != nil {
		t.Skipf("Android logging facilities are not accessible [%s]", err)
	}

	lr, err := NewLoggerReader(LogIdMain)
	assert.Nil(t, err)

	_, err = lr.ReadNext()
	assert.Nil(t, err)
}
