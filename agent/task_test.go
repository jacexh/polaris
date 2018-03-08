package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	assert := assert.New(t)

	task, err := NewSniffTask("10.0.1.36", 8000, 30*time.Second, ConsolePrinter{}.Handle)
	go task.Sniff()
	assert.Nil(err)
	time.Sleep(5 * time.Second)
	task.S.Close()
	assert.Equal(123, 456)
}
