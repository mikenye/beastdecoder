package df

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeDF11(t *testing.T) {
	// define test data
	// test data was captured by running `viewadsb --no-interactive | grep -B 4 -A 30 "DF:11" | grep -B 30 -m 1 '^[[:space:]]*$'`
	var testTable = []struct {
		data         []byte
		expectedDf   DownlinkFormat
		expectedAddr int
		expectedCa   int
	}{
		{
			data:         []byte{0x5d, 0x7c, 0x0a, 0x2b, 0xbd, 0xfa, 0xbb},
			expectedAddr: 0x7C0A2B,
			expectedCa:   5,
		},
		{
			data:         []byte{0x5e, 0x7c, 0x19, 0xf2, 0xc8, 0xd2, 0xd3},
			expectedAddr: 0x7C19F2,
			expectedCa:   6,
		},
		{
			data:         []byte{0x5f, 0x7c, 0x7f, 0x38, 0x81, 0x27, 0x4c},
			expectedAddr: 0x7C7F38,
			expectedCa:   7,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		msg := DecodeDF11(testData.data)
		assert.Equal(testData.expectedAddr, msg.ICAO, testMsg+fmt.Sprintf("%06x", msg.ICAO))
		assert.Equal(testData.expectedCa, msg.ca, testMsg+"ca")
	}

}
