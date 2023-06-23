package bds

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS09(t *testing.T) {

	// define test data
	var testTable = []struct {
		data []byte
		tc   int                     // Type Code
		st   AirborneVelocitySubType // Sub-Type
		ic   bool                    // Intent change flag
		ifr  bool                    // IFR capability flag
		nuc  int                     // Navigation uncertainty category for velocity

		// Sub-type specific fields for ground speed
		groundSpeedFields     BDS09FrameGroundSpeed
		groundSpeed           float64
		groundTrack           float64
		groundSpeedTrackValid bool

		// Sub-type specific fields for airspeed
		airSpeedFields     BDS09FrameAirSpeed
		airSpeed           float64
		airTrack           float64
		airSpeedTrackValid bool

		vrSrc verticalRateSource // Source bit for vertical rate
		svr   verticalRateSign   // Sign bit for vertical rate
		vr    int                // vertical rate
		sDif  int                // Sign bit for GNSS and Baro altitudes difference
		dAlt  int                // Difference between GNSS and Baro altitudes
	}{
		{
			data:                  []byte{0x99, 0x44, 0xC2, 0x83, 0x68, 0x2C, 0x01},
			tc:                    19,
			st:                    AirborneVelocityGroundSpeedSubsonic,
			ic:                    false,
			ifr:                   true,
			nuc:                   0,
			vrSrc:                 verticalRateSourceGNSS,
			svr:                   verticalRateDescent,
			vr:                    -640,
			groundSpeedTrackValid: true,
			groundSpeed:           194.7,
			groundTrack:           262.3,
			airSpeedTrackValid:    false,
			airSpeed:              0,
			airTrack:              0,
		},
		{
			data:                  []byte{0x99, 0x44, 0xA5, 0x04, 0xA8, 0x3C, 0x03},
			tc:                    19,
			st:                    AirborneVelocityGroundSpeedSubsonic,
			ic:                    false,
			ifr:                   true,
			nuc:                   0,
			vrSrc:                 verticalRateSourceGNSS,
			svr:                   verticalRateDescent,
			vr:                    -896,
			sDif:                  0,
			dAlt:                  50,
			groundSpeedTrackValid: true,
			groundSpeed:           167.9,
			groundTrack:           282.4,
			airSpeedTrackValid:    false,
			airSpeed:              0,
			airTrack:              0,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		frame, err := DecodeBDS09(testData.data)
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		assert.NoError(err, testMsg+"decodeBDS09 error")
		assert.Equal(testData.tc, frame.tc, testMsg+"tc")
		assert.Equal(testData.st, frame.st, testMsg+"st")
		assert.Equal(testData.ic, frame.ic, testMsg+"ic")
		assert.Equal(testData.ifr, frame.ifr, testMsg+"ifr")
		assert.Equal(testData.nuc, frame.nuc, testMsg+"nuc")
		assert.Equal(testData.vrSrc, frame.vrSrc, testMsg+"vrSrc")
		assert.Equal(testData.svr, frame.svr, testMsg+"svr")
		assert.Equal(testData.vr, frame.vr, testMsg+"vr")
		assert.Equal(testData.sDif, frame.sDif, testMsg+"sDif")
		assert.Equal(testData.dAlt, frame.dAlt, testMsg+"dAlt")
		assert.Equal(testData.groundSpeedTrackValid, frame.groundSpeedTrackValid, testMsg+"groundSpeedTrackValid")
		assert.Equal(testData.groundSpeed, math.Round(frame.groundSpeed*10)/10, testMsg+"groundSpeed")
		assert.Equal(testData.groundTrack, math.Round(frame.groundTrack*10)/10, testMsg+"groundTrack")
		assert.Equal(testData.airSpeedTrackValid, frame.airSpeedTrackValid, testMsg+"airSpeedTrackValid")
		assert.Equal(testData.airSpeed, math.Round(frame.airSpeed*10)/10, testMsg+"airSpeed")
		assert.Equal(testData.airTrack, math.Round(frame.airTrack*10)/10, testMsg+"airTrack")
	}
}
