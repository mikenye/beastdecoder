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
// id   int    // Identity code
// icao int    // Address announced: The address refers to the 24-bit transponder address (icao).
// mb   []byte // Message, Comm-B
// p    []byte // Parity

func TestDecodeDF21(t *testing.T) {
	// define test data
	var testTable = []struct {
		data         []byte
		expectedFs   int
		expectedDr   int
		expectedId   int
		expectedAddr int
		expectedMb   []byte
	}{
		{
			data:         []byte{0xa8, 0x00, 0x00, 0xbd, 0x00, 0x12, 0xb9, 0x13, 0x7f, 0xf4, 0x4f, 0x28, 0xec, 0x30},
			expectedFs:   0,
			expectedDr:   0,
			expectedId:   189,
			expectedAddr: 0x7C534D,
			expectedMb:   []byte{0x00, 0x12, 0xB9, 0x13, 0x7F, 0xF4, 0x4F},
			// *a80000bd0012b9137ff44f28ec30;
			// CRC: 7c534d
			// RSSI: -14.0 dBFS
			// Time: 426397432615.17us
			// DF:21 addr:7c534d FS:0 DR:0 UM:0 ID:189 MB:0012B9137FF44F
			//  Comm-B, Identity Reply
			//   Comm-B format: unknown format
			//   ICAO Address:  7C534D (Mode S / ADS-B)
			//   Air/Ground:    airborne?
			//   Squawk:        4307
		},
		{
			data:         []byte{0xa8, 0x00, 0x00, 0x84, 0x80, 0x18, 0xd5, 0x2e, 0xbf, 0xfc, 0xbf, 0xac, 0x89, 0x3b},
			expectedFs:   0,
			expectedDr:   0,
			expectedId:   132,
			expectedAddr: 0x7C6C27,
			expectedMb:   []byte{0x80, 0x18, 0xD5, 0x2E, 0xBF, 0xFC, 0xBF},
			// *a80000848018d52ebffcbfac893b;
			// CRC: 7c6c27
			// RSSI: -18.9 dBFS
			// Time: 426397574263.42us
			// DF:21 addr:7c6c27 FS:0 DR:0 UM:0 ID:132 MB:8018D52EBFFCBF
			//  Comm-B, Identity Reply
			//   Comm-B format: BDS5,0 Track and turn report
			//   ICAO Address:  7C6C27 (Mode S / ADS-B)
			//   Air/Ground:    airborne?
			//   Ground track   198.6
			//   Track rate:    -0.03 deg/sec left
			//   Roll:          0.0 degrees
			//   Groundspeed:   372.0 kt
			//   TAS:           382 kt
			//   Squawk:        4002
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		msg, err := DecodeDF21(testData.data)
		assert.NoError(err, testMsg+"DecodeDF21 error")
		assert.Equal(testData.expectedFs, msg.fs, testMsg+"fs")
		assert.Equal(testData.expectedDr, msg.dr, testMsg+"dr")
		assert.Equal(testData.expectedId, msg.id, testMsg+"id")
		assert.Equal(testData.expectedAddr, msg.ICAO, testMsg+fmt.Sprintf("%06x", msg.ICAO))
		assert.Equal(testData.expectedMb, msg.MB, testMsg+"mb")
	}
}
