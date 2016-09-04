package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeData(t *testing.T) {
	s := "0001000200030004000500060007"
	_, gpio, voltage, power, temperature := DecodeData(s)

	assert.Equal(t, gpio, 1)
	assert.Equal(t, voltage, 2)
	assert.Equal(t, power, 3)
	assert.Equal(t, temperature, 4)
}
