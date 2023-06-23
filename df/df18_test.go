package df

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test data was captured by running `viewadsb --no-interactive | grep -B 4 -A 30 "DF:18" | grep -B 30 -m 1 '^[[:space:]]*$'`

// ca   int    // Capability bits
// icao int    // Address announced: The address refers to the 24-bit transponder address (icao).
// tc   int    // message type code
// me   []byte // Message, extended squitter
// pi   []byte // Parity/Interrogator ID

func TestDecodeDF18(t *testing.T) {
	// define test data
	var testTable = []struct {
		data         []byte
		expectedCf   int
		expectedAddr int
		expectedTc   int
		expectedMe   []byte
		expectedPi   []byte
	}{
		{
			// Extended Squitter (Non-Transponder) Aircraft identification and category (2)
			data:         []byte{0x90, 0x7c, 0xf7, 0xc6, 0x10, 0x40, 0x84, 0x98, 0xc8, 0x18, 0x20, 0x00, 0x67, 0x90},
			expectedCf:   0,
			expectedAddr: 0x7CF7C6,
			expectedTc:   2,
			expectedMe:   []byte{0x10, 0x40, 0x84, 0x98, 0xC8, 0x18, 0x20},
			// CRC: 000000
			// RSSI: -22.6 dBFS
			// Time: 427022932841.92us
			// DF:18 AA: CF:0 ME:10408498C81820
			// Extended Squitter (Non-Transponder) Aircraft identification and category (2)
			// ICAO Address:  7CF7C6 (ADS-B, non-transponder)
			// Ident:         PHRX2A
			// Category:      C0

		},
		{
			// Extended Squitter (Non-Transponder) Reserved for surface system status (24/0)
			data:         []byte{0x90, 0x7c, 0xf7, 0xc6, 0xc1, 0x04, 0x00, 0x00, 0x00, 0x20, 0x04, 0x6e, 0xfb, 0xac},
			expectedCf:   0,
			expectedAddr: 0x7CF7C6,
			expectedTc:   24,
			expectedMe:   []byte{0xC1, 0x04, 0x00, 0x00, 0x00, 0x20, 0x04},
			// *907cf7c6c10400000020046efbac;
			// CRC: 000000
			// RSSI: -23.5 dBFS
			// Time: 426398730147.08us
			// DF:18 AA:7CF7C6 CF:0 ME:C1040000002004
			// Extended Squitter (Non-Transponder) Reserved for surface system status (24/0)
			// ICAO Address:  7CF7C6 (ADS-B, non-transponder)
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		msg := DecodeDF18(testData.data)
		assert.Equal(testData.expectedCf, msg.cf, testMsg+"cf")
		assert.Equal(testData.expectedAddr, msg.ICAO, testMsg+fmt.Sprintf("%06x", msg.ICAO))
		assert.Equal(testData.expectedTc, msg.Tc, testMsg+"tc")
		assert.Equal(testData.expectedMe, msg.ME, testMsg+"me")
	}
}
