package df

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeDF0(t *testing.T) {
	// define test data
	// test data was captured by running `viewadsb --no-interactive | grep -B 4 -A 30 "DF:0" | grep -B 30 -m 1 '^[[:space:]]*$'`
	var testTable = []struct {
		data            []byte
		expectedDf      DownlinkFormat
		expectedAddr    int
		expectedVs      int
		expectedCc      int
		expectedSl      int
		expectedRi      int
		expectedAc      int
		expectedBaroAlt float64
	}{
		{
			data:            []byte{0x02, 0x00, 0x08, 0x1c, 0x5b, 0xbe, 0x71},
			expectedDf:      DF0,
			expectedAddr:    0x7C8055,
			expectedVs:      0,
			expectedCc:      1,
			expectedSl:      0,
			expectedRi:      0,
			expectedAc:      2076,
			expectedBaroAlt: 12100,
		},
		{
			data:            []byte{0x02, 0x00, 0x01, 0x34, 0x24, 0x95, 0xf8},
			expectedDf:      DF0,
			expectedAddr:    0x7C7B80,
			expectedVs:      0,
			expectedCc:      1,
			expectedSl:      0,
			expectedRi:      0,
			expectedAc:      308,
			expectedBaroAlt: 1100,
		},
		{
			data:            []byte{0x00, 0x00, 0x02, 0x95, 0x98, 0x4e, 0x02},
			expectedAddr:    0x7C7F25,
			expectedVs:      0,
			expectedCc:      0,
			expectedSl:      0,
			expectedRi:      0,
			expectedAc:      661,
			expectedBaroAlt: 3125,
		},
		{
			data:            []byte{0x02, 0x00, 0x05, 0x93, 0xe5, 0xf5, 0xfa},
			expectedAddr:    0x7C3AD6,
			expectedVs:      0,
			expectedCc:      1,
			expectedSl:      0,
			expectedRi:      0,
			expectedAc:      1427,
			expectedBaroAlt: 7875,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		msg, err := DecodeDF0(testData.data)
		assert.NoError(err, testMsg+"DecodeDF0 error")
		assert.Equal(testData.expectedAddr, msg.ICAO, testMsg+fmt.Sprintf("%06x", msg.ICAO))
		assert.Equal(testData.expectedVs, msg.vs, testMsg+"vs")
		assert.Equal(testData.expectedCc, msg.cc, testMsg+"cc")
		assert.Equal(testData.expectedSl, msg.sl, testMsg+"sl")
		assert.Equal(testData.expectedRi, msg.ri, testMsg+"ri")
		assert.Equal(testData.expectedAc, msg.ac, testMsg+"ac")
		altFt, err := altitudeFromAltitudeCode13bit(msg.ac)
		assert.NoError(err, testMsg+"altitudeFromAltitudeCode error")
		assert.Equal(testData.expectedBaroAlt, altFt, testMsg+"altFt")
	}

}
