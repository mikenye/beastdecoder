package df

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeDF4(t *testing.T) {
	// define test data
	// test data was captured by running `viewadsb --no-interactive | grep -B 4 -A 30 "DF:4" | grep -B 30 -m 1 '^[[:space:]]*$'`
	var testTable = []struct {
		data            []byte
		expectedAddr    int
		expectedFs      int
		expectedDr      int
		expectedUm      int
		expectedAc      int
		expectedBaroAlt float64
	}{
		{
			data:            []byte{0x20, 0x00, 0x02, 0x94, 0xe7, 0xdc, 0x54},
			expectedAddr:    0x7C7F25,
			expectedFs:      0,
			expectedDr:      0,
			expectedUm:      0,
			expectedAc:      660,
			expectedBaroAlt: 3100,
		},
		{
			data:            []byte{0x20, 0x00, 0x01, 0x13, 0x0d, 0x19, 0x90},
			expectedAddr:    0x7C7A85,
			expectedFs:      0,
			expectedDr:      0,
			expectedUm:      0,
			expectedAc:      275,
			expectedBaroAlt: 675,
		},
		{
			data:            []byte{0x20, 0x00, 0x01, 0x31, 0x0c, 0xb2, 0x16},
			expectedAddr:    0x7C0CA8,
			expectedFs:      0,
			expectedDr:      0,
			expectedUm:      0,
			expectedAc:      305,
			expectedBaroAlt: 1025,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		msg, err := DecodeDF4(testData.data)
		assert.NoError(err, testMsg+"DecodeDF4 error")
		assert.Equal(testData.expectedAddr, msg.ICAO, testMsg+fmt.Sprintf("%06x", msg.ICAO))
		assert.Equal(testData.expectedFs, msg.fs, testMsg+"fs")
		assert.Equal(testData.expectedDr, msg.dr, testMsg+"dr")
		assert.Equal(testData.expectedUm, msg.um, testMsg+"um")
		assert.Equal(testData.expectedAc, msg.ac, testMsg+"ac")
		altFt, err := altitudeFromAltitudeCode13bit(msg.ac)
		assert.NoError(err, testMsg+"altitudeFromAltitudeCode error")
		assert.Equal(testData.expectedBaroAlt, altFt, testMsg+"altFt")
	}

}
