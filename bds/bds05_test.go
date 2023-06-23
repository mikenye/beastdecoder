package bds

import (
	"beastdecoder/common"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS05(t *testing.T) {
	// define test data
	var testTable = []struct {
		data        []byte
		tc          int                // Type Code
		ss          SurveillanceStatus // Surveillance status
		alt         float64            // Altitude
		f           common.CprFormat   // CPR Format
		latCpr      int                // Encoded latitude
		lonCpr      int                // Encoded longitude
		errExpected bool
	}{
		{
			data:        []byte{0x58, 0x13, 0xA7, 0x18, 0x29, 0x97, 0xF8},
			tc:          11,
			ss:          SurveillanceStatusNoCondition,
			alt:         2850,
			f:           common.CprFormatOddFrame,
			latCpr:      101396,
			lonCpr:      104440,
			errExpected: false,
		},
		{
			data:        []byte{0xB0, 0x0B, 0xF6, 0xF5, 0xC1, 0x81, 0x7B},
			tc:          22,
			ss:          SurveillanceStatusNoCondition,
			alt:         0,
			f:           common.CprFormatOddFrame,
			latCpr:      96992,
			lonCpr:      98683,
			errExpected: true,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		frame, err := DecodeBDS05(testData.data)
		if testData.errExpected {
			assert.Error(err, testMsg+"decodeBDS05 no error")
		} else {
			assert.NoError(err, testMsg+"decodeBDS05 error")
		}
		assert.Equal(testData.tc, frame.Tc, testMsg+"tc")
		assert.Equal(testData.ss, frame.ss, testMsg+"ss")
		assert.Equal(testData.alt, frame.Altitude, testMsg+"alt")
		assert.Equal(testData.f, frame.F, testMsg+"f")
		assert.Equal(testData.latCpr, frame.LatCpr, testMsg+"latCpr")
		assert.Equal(testData.lonCpr, frame.LonCpr, testMsg+"lonCpr")
	}
}
