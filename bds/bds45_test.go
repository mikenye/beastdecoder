package bds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS45turbulence(t *testing.T) {
	testMsg := []byte{0b11100000}
	turb, err := decodeBDS45turbulence(testMsg)
	assert.NoError(t, err)
	assert.Equal(t, 3, turb)
}

func TestDecodeBDS45windShear(t *testing.T) {
	testMsg := []byte{0b00011100}
	ws, err := decodeBDS45windShear(testMsg)
	assert.NoError(t, err)
	assert.Equal(t, 3, ws)
}

func TestDecodeBDS45staticAirTemperature(t *testing.T) {
	mb := []byte{0x00, 0x01, 0xCF, 0x40, 0x00, 0x00, 0x00}
	sat, err := decodeBDS45staticAirTemperature(mb)
	assert.NoError(t, err)
	assert.Equal(t, -48.75, sat)
}
