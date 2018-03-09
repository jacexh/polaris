package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	assert := assert.New(t)

	task, err := NewSniffTask("127.0.0.1", 8000, 10*time.Second, ConsolePrinter{}.Handle)
	assert.Nil(err)
	start := time.Now()
	task.Sniff()
	duration := time.Since(start)
	assert.True(duration >= 10*time.Second)
}

func TestManualStop(t *testing.T) {
	assert := assert.New(t)

	task, err := NewSniffTask("127.0.0.1", 8000, 10*time.Second, ConsolePrinter{}.Handle)
	assert.Nil(err)
	start := time.Now()
	go func() {
		time.Sleep(5 * time.Second)
		task.Stop()
	}()
	err = task.Sniff()
	assert.Nil(err)
	duration := time.Since(start)
	assert.True(duration < 10*time.Second)
	assert.True(duration >= 5*time.Second)
}
