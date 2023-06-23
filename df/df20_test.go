package df

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test data was captured by running `viewadsb --no-interactive | grep -B 4 -A 30 "DF:20" | grep -B 30 -m 1 '^[[:space:]]*$'`

// fs   int    // Flight status
// dr   int    // Downlink request
// um   int    // Utility message
// ac   int    // Altitude code
// icao int    // Address announced: The address refers to the 24-bit transponder address (icao).
// mb   []byte // Message, Comm-B
// p    []byte // Parity

func TestDecodeDF20(t *testing.T) {
	// define test data
	var testTable = []struct {
		data         []byte
		expectedFs   int
		expectedDr   int
		expectedAc   int
		expectedAddr int
		expectedMb   []byte
		expectedPi   []byte
	}{
		{
			// Comm-B, Altitude Reply
			data:         []byte{0xa0, 0x00, 0x02, 0xbf, 0x10, 0x02, 0x0a, 0x80, 0xf0, 0x00, 0x00, 0x1b, 0x43, 0x5f},
			expectedFs:   0,
			expectedDr:   0,
			expectedAc:   703,
			expectedAddr: 0x7CF9DA,
			expectedMb:   []byte{0x10, 0x02, 0x0A, 0x80, 0xF0, 0x00, 0x00},
			// *a00002bf10020a80f000001b435f;
			// CRC: 7cf9da
			// RSSI: -11.6 dBFS
			// Time: 426397425102.25us
			// DF:20 addr:7CF9DA FS:0 DR:0 UM:0 AC:703 MB:10020A80F00000
			//  Comm-B, Altitude Reply
			//   Comm-B format: BDS1,0 Datalink capabilities
			//   ICAO Address:  7CF9DA (Mode S / ADS-B)
			//   Air/Ground:    airborne?
			//   Baro altitude: 3775 ft
		},
		{
			// Comm-B, Altitude Reply
			data:         []byte{0xa0, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xd5, 0x0a, 0x7d},
			expectedFs:   0,
			expectedDr:   0,
			expectedAc:   0,
			expectedAddr: 0x7C175A,
			expectedMb:   []byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			// *a0000000ff000000000000d50a7d;
			// CRC: 7c175a
			// RSSI: -24.0 dBFS
			// Time: 426397582458.83us
			// DF:20 addr:7C175A FS:0 DR:0 UM:0 AC:0 MB:FF000000000000
			//  Comm-B, Altitude Reply
			//   Comm-B format: BDS1,7 Common usage GICB capabilities
			//   ICAO Address:  7C175A (Mode S / ADS-B)
			//   Air/Ground:    airborne?
		},
		{
			// Comm-B, Altitude Reply
			data:         []byte{0xa0, 0x00, 0x01, 0x28, 0x20, 0x0c, 0x14, 0xa0, 0x82, 0x08, 0x20, 0x71, 0x52, 0x31},
			expectedFs:   0,
			expectedDr:   0,
			expectedAc:   296,
			expectedAddr: 0x7C0A31,
			expectedMb:   []byte{0x20, 0x0C, 0x14, 0xA0, 0x82, 0x08, 0x20},
			// *a0000128200c14a0820820715231;
			// CRC: 7c0a31
			// RSSI: -23.0 dBFS
			// Time: 426397729899.00us
			// DF:20 addr:7C0A31 FS:0 DR:0 UM:0 AC:296 MB:200C14A0820820
			//  Comm-B, Altitude Reply
			//   Comm-B format: BDS2,0 Aircraft identification
			//   ICAO Address:  7C0A31 (Mode S / ADS-B)
			//   Air/Ground:    airborne?
			//   Baro altitude: 800 ft
			//   Ident:         CAR
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		msg, err := DecodeDF20(testData.data)
		assert.NoError(err, testMsg+"DecodeDF20 error")
		assert.Equal(testData.expectedFs, msg.fs, testMsg+"fs")
		assert.Equal(testData.expectedDr, msg.dr, testMsg+"dr")
		assert.Equal(testData.expectedAc, msg.ac, testMsg+"ac")
		assert.Equal(testData.expectedAddr, msg.ICAO, testMsg+fmt.Sprintf("%06x", msg.ICAO))
		assert.Equal(testData.expectedMb, msg.MB, testMsg+"mb")
	}
}
