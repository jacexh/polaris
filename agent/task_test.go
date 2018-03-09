package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	assert := assert.New(t)

	task, err := NewSniffTask("10.0.1.36", 8000, 30*time.Second, ConsolePrinter{}.Handle)
	assert.Nil(err)
	start := time.Now()
	task.Sniff()
	duration := time.Since(start)
	assert.True(duration >= 30*time.Second)
	assert.True(false)
}
