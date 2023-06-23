package bds

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS50(t *testing.T) {

	// define test data
	var testTable = []struct {
		data                            []byte
		expectedRollAngleValid          bool
		expectedRollAngle               float64
		expectedTrueTrackAngleValid     bool
		expectedTrueTrackAngle          float64
		expectedGroundSpeedValid        bool
		expectedGroundSpeed             float64
		expectedTrueTrackAngleRateValid bool
		expectedTrueTrackAngleRate      float64
		expectedTrueAirspeedValid       bool
		expectedTrueAirspeed            float64
	}{
		{
			data:                            []byte{0x80, 0x18, 0xD5, 0x2E, 0xBF, 0xFC, 0xBF},
			expectedRollAngleValid:          true,
			expectedRollAngle:               0.0,
			expectedTrueTrackAngleValid:     true,
			expectedTrueTrackAngle:          198.6,
			expectedGroundSpeedValid:        true,
			expectedGroundSpeed:             372.0,
			expectedTrueTrackAngleRateValid: true,
			expectedTrueTrackAngleRate:      -0.03,
			expectedTrueAirspeedValid:       true,
			expectedTrueAirspeed:            382,
		},
		{
			data:                            []byte{0xE8, 0xBA, 0x45, 0x1F, 0xBD, 0xD4, 0x7E},
			expectedRollAngleValid:          true,
			expectedRollAngle:               -32.9,
			expectedTrueTrackAngleValid:     true,
			expectedTrueTrackAngle:          231.0,
			expectedGroundSpeedValid:        true,
			expectedGroundSpeed:             252.0,
			expectedTrueTrackAngleRateValid: true,
			expectedTrueTrackAngleRate:      -2.19,
			expectedTrueAirspeedValid:       true,
			expectedTrueAirspeed:            252,
		},
		{
			data:                            []byte{0xE8, 0xDA, 0x3D, 0x1F, 0xBD, 0xEC, 0x7D},
			expectedRollAngleValid:          true,
			expectedRollAngle:               -32.7,
			expectedTrueTrackAngleValid:     true,
			expectedTrueTrackAngle:          230.3,
			expectedGroundSpeedValid:        true,
			expectedGroundSpeed:             252.0,
			expectedTrueTrackAngleRateValid: true,
			expectedTrueTrackAngleRate:      -2.09,
			expectedTrueAirspeedValid:       true,
			expectedTrueAirspeed:            250,
		},
		{
			data:                            []byte{0x97, 0x71, 0xBF, 0x18, 0xA3, 0x0C, 0x65},
			expectedRollAngleValid:          true,
			expectedRollAngle:               32.9,
			expectedTrueTrackAngleValid:     true,
			expectedTrueTrackAngle:          39.2,
			expectedGroundSpeedValid:        true,
			expectedGroundSpeed:             196.0,
			expectedTrueTrackAngleRateValid: true,
			expectedTrueTrackAngleRate:      3.03,
			expectedTrueAirspeedValid:       true,
			expectedTrueAirspeed:            202,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		frame, err := DecodeBDS50(testData.data)
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		assert.NoError(err, testMsg+"decodeBDS50 error")
		assert.Equal(testData.expectedRollAngleValid, frame.RollAngleValid, testMsg+"RollAngleValid")
		assert.Equal(testData.expectedRollAngle, math.Round(frame.RollAngle*10)/10, testMsg+"RollAngle")
		assert.Equal(testData.expectedTrueTrackAngleValid, frame.TrueTrackAngleValid, testMsg+"TrueTrackAngleValid")
		assert.Equal(testData.expectedTrueTrackAngle, math.Round(frame.TrueTrackAngle*10)/10, testMsg+"TrueTrackAngle")
		assert.Equal(testData.expectedGroundSpeedValid, frame.GroundSpeedValid, testMsg+"GroundSpeedValid")
		assert.Equal(testData.expectedGroundSpeed, frame.GroundSpeed, testMsg+"GroundSpeed")
		assert.Equal(testData.expectedTrueAirspeedValid, frame.TrueAirspeedValid, testMsg+"TrueAirspeedValid")
		assert.Equal(testData.expectedTrueAirspeed, frame.TrueAirspeed, testMsg+"TrueAirspeed")
	}

}

// func TestDecodeBDS50roll_neg(t *testing.T) {
// 	// a8001c10ff98a313bffc54b6366d
// 	data := []byte{0xa8, 0x00, 0x1c, 0x10, 0xff, 0x98, 0xa3, 0x13, 0xbf, 0xfc, 0x54, 0xb6, 0x36, 0x6d}
// 	msg := DecodeDF21(data)
// 	r, e := decodeBDS50roll(msg.mb)
// 	assert.NoError(t, e)
// 	assert.Equal(t, -0.7, math.Round(r*10)/10)
// }

// func TestDecodeBDS50roll_pos(t *testing.T) {
// 	// a8000099805021247ffc8d5f8eb4
// 	data := []byte{0xa8, 0x00, 0x00, 0x99, 0x80, 0x50, 0x21, 0x24, 0x7f, 0xfc, 0x8d, 0x5f, 0x8e, 0xb4}
// 	msg := decodeDF21(data)
// 	r, e := decodeBDS50roll(msg.mb)
// 	assert.NoError(t, e)
// 	assert.Equal(t, 0.4, math.Round(r*10)/10)
// }

// func TestDecodeBDS50groundSpeed(t *testing.T) {
// 	data := []byte{0xA8, 0x00, 0x06, 0xAC, 0xF9, 0x36, 0x3D, 0x3B, 0xBF, 0x9C, 0xE9, 0x8F, 0x1E, 0x1D}
// 	df := getDF(data)
// 	assert.Equal(t, DF21, df)
// 	msg := decodeDF21(data)
// 	gs, err := decodeBDS50groundSpeed(msg.mb)
// 	assert.NoError(t, err)
// 	assert.Equal(t, float64(476), gs)
// }

// func TestDecodeBDS50trueAirspeed(t *testing.T) {
// 	data := []byte{0xA8, 0x00, 0x06, 0xAC, 0xF9, 0x36, 0x3D, 0x3B, 0xBF, 0x9C, 0xE9, 0x8F, 0x1E, 0x1D}
// 	df := getDF(data)
// 	assert.Equal(t, DF21, df)
// 	msg := decodeDF21(data)
// 	tas, err := decodeBDS50trueAirspeed(msg.mb)
// 	assert.NoError(t, err)
// 	assert.Equal(t, float64(466), tas)
// }
