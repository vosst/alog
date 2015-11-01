package alog

import (
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockLoggerAbiExtension struct {
	mock.Mock
}

func (self *MockLoggerAbiExtension) Read(reader io.Reader) (map[string]interface{}, error) {
	args := self.Called(reader)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

type MockReader struct {
	mock.Mock
}

func (self *MockReader) Read(p []byte) (int, error) {
	args := self.Called(p)
	return args.Int(0), args.Error(1)
}

func TestChainedLoggerAbiExtensionsCallsIntoAllExtensions(t *testing.T) {
	mr := MockReader{}

	result := map[string]interface{}{}

	mle1 := &MockLoggerAbiExtension{}
	mle2 := &MockLoggerAbiExtension{}

	mle1.On("Read", mock.Anything).Return(result, nil)
	mle2.On("Read", mock.Anything).Return(result, nil)

	ch := ChainedLoggerAbiExtension{[]LoggerAbiExtension{mle1, mle2}}
	ch.Read(&mr)

	mle1.AssertExpectations(t)
	mle2.AssertExpectations(t)
}
