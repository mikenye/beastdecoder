package bds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS08(t *testing.T) {
	// define test data
	var testTable = []struct {
		data       []byte
		tc         int
		ca         int
		callsign   string
		capability string
	}{
		{
			data:       []byte{0x20, 0x39, 0x72, 0xF2, 0xE7, 0x3C, 0x60},
			callsign:   "NWK2931",
			capability: "A0",
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		frame, err := DecodeBDS08(testData.data)
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		assert.NoError(err, testMsg+"decodeBDS08 error")
		assert.Equal(testData.callsign, frame.Callsign, testMsg+"callsign")
		assert.Equal(testData.capability, frame.capability, testMsg+"capability")
	}
}
