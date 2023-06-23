package bds

import (
	"beastdecoder/common"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS06(t *testing.T) {

	// define test data
	var testTable = []struct {
		data        []byte
		tc          int               // Type Code
		groundSpeed string            // Ground speed (decoded)
		s           GroundTrackStatus // Status for ground track
		groundTrack string            // Ground track (decoded)
		f           common.CprFormat  // CPR Format
		latCpr      int               // Encoded latitude
		lonCpr      int               // Encoded longitude
	}{
		{
			data:        []byte{0x39, 0x4B, 0xF2, 0xD8, 0x94, 0xD9, 0x1E},
			tc:          7,
			groundSpeed: "10.1860 km/h (5.500 kt)",
			s:           GroundTrackStatusValid,
			groundTrack: "177.1875°",
			f:           common.CprFormatEvenFrame,
			latCpr:      93258,
			lonCpr:      55582,
		},
		{
			data:        []byte{0x39, 0x4B, 0xF4, 0x43, 0xE4, 0x45, 0x6A},
			tc:          7,
			groundSpeed: "10.1860 km/h (5.500 kt)",
			s:           GroundTrackStatusValid,
			groundTrack: "177.1875°",
			f:           common.CprFormatOddFrame,
			latCpr:      8690,
			lonCpr:      17770,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		frame, err := DecodeBDS06(testData.data)
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		assert.NoError(err, testMsg+"decodeBDS06 error")
		assert.Equal(testData.tc, frame.tc, testMsg+"tc")
		assert.Equal(testData.groundSpeed, frame.GroundSpeed, testMsg+"groundSpeed")
		assert.Equal(testData.s, frame.S, testMsg+"s")
		assert.Equal(testData.groundTrack, frame.GroundTrack, testMsg+"groundTrack")
		assert.Equal(testData.f, frame.F, testMsg+"f")
		assert.Equal(testData.latCpr, frame.LatCpr, testMsg+"latCpr")
		assert.Equal(testData.lonCpr, frame.LonCpr, testMsg+"lonCpr")

	}
}
