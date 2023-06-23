package df

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test data was captured by running `viewadsb --no-interactive | grep -B 4 -A 30 "DF:17" | grep -B 30 -m 1 '^[[:space:]]*$'`

// ca   int    // Capability
// icao int    // Address announced: The address refers to the 24-bit transponder address (icao).
// tc   int    // message type code
// me   []byte // Message, extended squitter
// pi   []byte // Parity/Interrogator ID

func TestDecodeDF17(t *testing.T) {
	// define test data
	var testTable = []struct {
		data         []byte
		expectedCa   int
		expectedAddr int
		expectedTc   int
		expectedMe   []byte
		expectedPi   []byte // viewadsb doesn't seem to calculate/provide CRC for DF17 frames...?
	}{
		{
			// Extended Squitter Aircraft identification and category (4)
			data:         []byte{0x8d, 0x7c, 0xf9, 0xd9, 0x21, 0x58, 0x94, 0x12, 0xd3, 0x18, 0x20, 0xf9, 0x88, 0x69},
			expectedCa:   5,
			expectedAddr: 0x7CF9D9,
			expectedTc:   4,
			expectedMe:   []byte{0x21, 0x58, 0x94, 0x12, 0xD3, 0x18, 0x20},
		},
		{
			// Extended Squitter Surface position (7)
			data:         []byte{0x8c, 0x7c, 0x6d, 0x26, 0x38, 0xee, 0x54, 0x3d, 0x9e, 0x47, 0x33, 0x48, 0xd0, 0x76},
			expectedCa:   4,
			expectedAddr: 0x7C6D26,
			expectedTc:   7,
			expectedMe:   []byte{0x38, 0xEE, 0x54, 0x3D, 0x9E, 0x47, 0x33},
		},
		{
			// Extended Squitter Airborne position (barometric altitude) (11)
			data:         []byte{0x8d, 0x7c, 0x42, 0xe4, 0x58, 0x4b, 0x32, 0xf4, 0x5a, 0xec, 0x29, 0x4b, 0x03, 0x56},
			expectedAddr: 0x7C42E4,
			expectedTc:   11,
			expectedCa:   5,
			expectedMe:   []byte{0x58, 0x4B, 0x32, 0xF4, 0x5A, 0xEC, 0x29},
		},
		{
			// Extended Squitter Airborne velocity over ground, subsonic (19/1)
			data:         []byte{0x8e, 0x7c, 0x44, 0x23, 0x99, 0x04, 0x92, 0x96, 0x88, 0x04, 0x04, 0xa5, 0x20, 0x92},
			expectedAddr: 0x7C4423,
			expectedTc:   19,
			expectedCa:   6,
			expectedMe:   []byte{0x99, 0x04, 0x92, 0x96, 0x88, 0x04, 0x04},
		},
		{
			// Extended Squitter Airborne position (geometric altitude) (22)
			data:         []byte{0x8e, 0x7c, 0x17, 0x5a, 0xb0, 0x0b, 0xe7, 0x0e, 0xd1, 0x8b, 0x14, 0xba, 0x18, 0xb3},
			expectedAddr: 0x7C175A,
			expectedTc:   22,
			expectedCa:   6,
			expectedMe:   []byte{0xB0, 0x0B, 0xE7, 0x0E, 0xD1, 0x8B, 0x14},
		},
		{
			// Extended Squitter National use / 1090-WP-15-20 Mode A squawk (23/7)
			data:         []byte{0x8d, 0x7c, 0x80, 0x33, 0xbf, 0xf0, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x7d, 0x68, 0x8e},
			expectedCa:   5,
			expectedAddr: 0x7C8033,
			expectedTc:   23,
			expectedMe:   []byte{0xBF, 0xF0, 0xF0, 0x00, 0x00, 0x00, 0x00},
		},
		{
			// Extended Squitter Emergency/priority status (28/1)
			data:         []byte{0x8d, 0x7c, 0xf9, 0xd9, 0xe1, 0x12, 0x04, 0x00, 0x00, 0x00, 0x00, 0x4a, 0x55, 0x2e},
			expectedCa:   5,
			expectedAddr: 0x7CF9D9,
			expectedTc:   28,
			expectedMe:   []byte{0xE1, 0x12, 0x04, 0x00, 0x00, 0x00, 0x00},
		},
		{
			// Extended Squitter Target state and status (V2) (29/1)
			data:         []byte{0x8d, 0x7c, 0x68, 0x12, 0xea, 0x06, 0xe8, 0x94, 0xab, 0x3e, 0x00, 0xa7, 0xa7, 0x0a},
			expectedAddr: 0x7C6812,
			expectedTc:   29,
			expectedCa:   5,
			expectedMe:   []byte{0xEA, 0x06, 0xE8, 0x94, 0xAB, 0x3E, 0x00},
		},
		{
			// Extended Squitter Aircraft operational status (airborne) (31/0)
			data:         []byte{0x8d, 0x7c, 0xf9, 0xe8, 0xf8, 0x02, 0x00, 0x02, 0x00, 0x49, 0xb8, 0xf7, 0x2d, 0xe1},
			expectedAddr: 0x7CF9E8,
			expectedTc:   31,
			expectedCa:   5,
			expectedMe:   []byte{0xF8, 0x02, 0x00, 0x02, 0x00, 0x49, 0xB8},
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		msg := DecodeDF17(testData.data)
		assert.Equal(testData.expectedCa, msg.ca, testMsg+"ca")
		assert.Equal(testData.expectedAddr, msg.ICAO, testMsg+fmt.Sprintf("%06x", msg.ICAO))
		assert.Equal(testData.expectedTc, msg.Tc, testMsg+"tc")
		assert.Equal(testData.expectedMe, msg.ME, testMsg+"me")

	}
}
