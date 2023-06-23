package df

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeDF16(t *testing.T) {
	// define test data
	// test data was captured by running `viewadsb --no-interactive | grep -B 4 -A 30 "DF:16" | grep -B 30 -m 1 '^[[:space:]]*$'`
	var testTable = []struct {
		data            []byte
		expectedDf      DownlinkFormat
		expectedAddr    int
		expectedVs      int
		expectedSl      int
		expectedRi      int
		expectedAc      int
		expectedMv      []byte
		expectedBaroAlt float64
	}{
		{
			data:            []byte{0x80, 0x00, 0x14, 0xa0, 0x58, 0xa6, 0x02, 0x9b, 0xa0, 0x44, 0xad, 0x5a, 0x76, 0x52},
			expectedAddr:    0x7C5E8F,
			expectedVs:      0,
			expectedSl:      0,
			expectedRi:      0,
			expectedAc:      5280,
			expectedMv:      []byte{0x58, 0xA6, 0x02, 0x9B, 0xA0, 0x44, 0xAD},
			expectedBaroAlt: 3100,
		},
		{
			data:            []byte{0x80, 0xc1, 0x87, 0xa0, 0x58, 0x3e, 0x02, 0xbd, 0xf6, 0x3f, 0x4a, 0x6f, 0x6d, 0x0e},
			expectedAddr:    0x7C5343,
			expectedVs:      0,
			expectedSl:      6,
			expectedRi:      3,
			expectedAc:      1952,
			expectedMv:      []byte{0x58, 0x3E, 0x02, 0xBD, 0xF6, 0x3F, 0x4A},
			expectedBaroAlt: 10600,
		},
		{
			data:            []byte{0x80, 0xa1, 0x86, 0x80, 0x58, 0x34, 0x07, 0x01, 0x57, 0xa7, 0x4e, 0x38, 0xc7, 0x48},
			expectedAddr:    0x7C5343,
			expectedVs:      0,
			expectedSl:      5,
			expectedRi:      3,
			expectedAc:      1664,
			expectedMv:      []byte{0x58, 0x34, 0x07, 0x01, 0x57, 0xA7, 0x4E},
			expectedBaroAlt: 7000,
		},
		{
			data:            []byte{0x80, 0x00, 0x01, 0xb5, 0x58, 0x0f, 0x52, 0x9e, 0x00, 0x30, 0x92, 0x2f, 0x7b, 0xfd},
			expectedAddr:    0x7C49FB,
			expectedVs:      0,
			expectedSl:      0,
			expectedRi:      0,
			expectedAc:      437,
			expectedMv:      []byte{0x58, 0x0F, 0x52, 0x9E, 0x00, 0x30, 0x92},
			expectedBaroAlt: 1925,
		},
		{
			data:            []byte{0x80, 0xc1, 0x8c, 0x90, 0x58, 0x65, 0x07, 0x9d, 0xd8, 0x61, 0xb2, 0xe0, 0x8e, 0xe0},
			expectedAddr:    0x7C2BF6,
			expectedVs:      0,
			expectedSl:      6,
			expectedRi:      3,
			expectedAc:      3216,
			expectedMv:      []byte{0x58, 0x65, 0x07, 0x9D, 0xD8, 0x61, 0xB2},
			expectedBaroAlt: 19000,
		},
		{
			data:            []byte{0x80, 0xe1, 0x94, 0xbd, 0x58, 0xa7, 0xd3, 0xac, 0x87, 0x02, 0x2f, 0x1b, 0x02, 0xc1},
			expectedAddr:    0x7C6C35,
			expectedVs:      0,
			expectedSl:      7,
			expectedRi:      3,
			expectedAc:      5309,
			expectedMv:      []byte{0x58, 0xA7, 0xD3, 0xAC, 0x87, 0x02, 0x2F},
			expectedBaroAlt: 32525,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		msg, err := DecodeDF16(testData.data)
		assert.NoError(err, testMsg+"DeodeDF16 error")
		assert.Equal(testData.expectedAddr, msg.ICAO, testMsg+fmt.Sprintf("%06x", msg.ICAO))
		assert.Equal(testData.expectedVs, msg.vs, testMsg+"vs")
		assert.Equal(testData.expectedSl, msg.sl, testMsg+"sl")
		assert.Equal(testData.expectedRi, msg.ri, testMsg+"ri")
		assert.Equal(testData.expectedAc, msg.ac, testMsg+"ac")
		assert.Equal(testData.expectedMv, msg.mv, testMsg+"mv")
		altFt, err := altitudeFromAltitudeCode13bit(msg.ac)
		assert.NoError(err, testMsg+"altitudeFromAltitudeCode error")
		assert.Equal(testData.expectedBaroAlt, altFt, testMsg+"altFt")
	}
}
