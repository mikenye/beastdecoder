package bds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS20(t *testing.T) {
	// define test data
	var testTable = []struct {
		data             []byte
		expectedCallsign string
	}{
		{
			data:             []byte{0x20, 0x0C, 0x14, 0xA0, 0x82, 0x08, 0x20},
			expectedCallsign: "CAR",
		},
		{
			data:             []byte{0x20, 0x39, 0x72, 0xF1, 0xE3, 0x0D, 0x60},
			expectedCallsign: "NWK1805",
		},
		{
			data:             []byte{0x20, 0x55, 0x46, 0x77, 0xE7, 0x78, 0x20},
			expectedCallsign: "UTY797",
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		frame, err := DecodeBDS20(testData.data)
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		assert.NoError(err, testMsg+"decodeBDS20 error")
		assert.Equal(testData.expectedCallsign, frame.Callsign, testMsg+"callsign")
	}
}
