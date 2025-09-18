package main

import (
	"bytes"
	"testing"

	"github.com/gcalvocr/go-testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Example service that uses ReadWriter
type DataProcessor struct {
	rw ReadWriter
}

func NewDataProcessor(rw ReadWriter) *DataProcessor {
	return &DataProcessor{rw: rw}
}

func (dp *DataProcessor) ProcessData(input []byte) ([]byte, error) {
	// Write input
	_, err := dp.rw.Write(input)
	if err != nil {
		return nil, err
	}

	// Read back (simplified example)
	output := make([]byte, len(input))
	_, err = dp.rw.Read(output)
	return output, err
}

func TestDataProcessor_ProcessData_StandardMock(t *testing.T) {
	mockRW := &mocks.MockReadWriter{}

	input := []byte("hello world")
	expectedOutput := []byte("HELLO WORLD")

	// Set up expectations
	mockRW.On("Write", mock.Anything).Return(len(input), nil)
	mockRW.On("Read", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
		buf := args.Get(0).([]byte)
		copy(buf, expectedOutput)
	}).Return(len(expectedOutput), nil)

	processor := NewDataProcessor(mockRW)

	result, err := processor.ProcessData(input)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockRW.AssertExpectations(t)
}

func TestDataProcessor_ProcessData_WithExpect(t *testing.T) {
	mockRW := &mocks.MockReadWriter{}

	input := []byte("hello world")
	expectedOutput := []byte("HELLO WORLD")

	// Set up expectations using standard mock (since ReadWriter uses standard mocks)
	mockRW.On("Write", input).Return(len(input), nil)
	mockRW.On("Read", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
		buf := args.Get(0).([]byte)
		copy(buf, expectedOutput)
	}).Return(len(expectedOutput), nil)

	processor := NewDataProcessor(mockRW)

	result, err := processor.ProcessData(input)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, result)
	mockRW.AssertExpectations(t)
}

func TestReadWriter_BufferedReadWrite(t *testing.T) {
	mockRW := &mocks.MockReadWriter{}

	var buf bytes.Buffer
	data := []byte("test data")

	// Simulate writing to buffer
	mockRW.On("Write", data).Run(func(args mock.Arguments) {
		buf.Write(args.Get(0).([]byte))
	}).Return(len(data), nil).Once()

	// Simulate reading from buffer
	mockRW.On("Read", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
		readBuf := args.Get(0).([]byte)
		n, _ := buf.Read(readBuf)
		mockRW.MethodCalled("Read", readBuf[:n])
	}).Return(len(data), nil).Once()

	processor := NewDataProcessor(mockRW)

	result, err := processor.ProcessData(data)

	assert.NoError(t, err)
	assert.Equal(t, data, result)
	mockRW.AssertExpectations(t)
}
