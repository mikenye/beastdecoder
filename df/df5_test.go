package df

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeDF5(t *testing.T) {
	// define test data
	// test data was captured by running `viewadsb --no-interactive | grep -B 4 -A 30 "DF:5" | grep -B 30 -m 1 '^[[:space:]]*$'`
	var testTable = []struct {
		data           []byte
		expectedAddr   int
		expectedFs     int
		expectedDr     int
		expectedUm     int
		expectedId     int
		expectedSquawk int
	}{
		{
			data:           []byte{0x28, 0x00, 0x0a, 0x00, 0x30, 0x8d, 0xe4},
			expectedAddr:   0x7C822A,
			expectedFs:     0,
			expectedDr:     0,
			expectedUm:     0,
			expectedId:     2560,
			expectedSquawk: 3000,
		},
		{
			data:           []byte{0x28, 0x00, 0x05, 0x01, 0x95, 0xdc, 0x9b},
			expectedAddr:   0x7C6BDC,
			expectedFs:     0,
			expectedDr:     0,
			expectedUm:     0,
			expectedId:     1281,
			expectedSquawk: 64,
		},
		{
			data:           []byte{0x28, 0x00, 0x01, 0x1f, 0xad, 0x39, 0x76},
			expectedAddr:   0x7C0CA8,
			expectedFs:     0,
			expectedDr:     0,
			expectedUm:     0,
			expectedId:     287,
			expectedSquawk: 647,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		msg, err := DecodeDF5(testData.data)
		assert.NoError(err, testMsg+"DecodeDF5 error")
		assert.Equal(testData.expectedAddr, msg.ICAO, testMsg+fmt.Sprintf("%06x", msg.ICAO))
		assert.Equal(testData.expectedFs, msg.fs, testMsg+"fs")
		assert.Equal(testData.expectedDr, msg.dr, testMsg+"dr")
		assert.Equal(testData.expectedUm, msg.um, testMsg+"um")
		assert.Equal(testData.expectedId, msg.id, testMsg+"id")
		squawk, err := squawkFromIdentityCode(msg.id)
		assert.NoError(err, testMsg+"squawk error")
		assert.Equal(testData.expectedSquawk, squawk, testMsg+"squawk")
	}

}
