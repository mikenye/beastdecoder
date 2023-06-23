package bds

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MagneticHeadingValid bool
// MagneticHeading      float64

// IndicatedAirspeedValid bool
// IndicatedAirspeed      float64

// MachNumberValid bool
// MachNumber      float64

// BarometricAltitudeRateValid bool
// BarometricAltitudeRate      float64

// InertialVerticalVelocityValid bool
// InertialVerticalVelocity      float64

func TestDecodeBDS60(t *testing.T) {

	// define test data
	var testTable = []struct {
		data                                []byte
		expectedMagneticHeadingValid        bool
		expectedMagneticHeading             float64
		expectedIndicatedAirspeedValid      bool
		expectedIndicatedAirspeed           float64
		expectedMachNumberValid             bool
		expectedMachNumber                  float64
		expectedBarometricAltitudeRateValid bool
		expectedBarometricAltitudeRate      float64
		GNSSAltitudeRateValid               bool
		GNSSAltitudeRate                    float64
	}{
		{
			data:                                []byte{0xCE, 0x69, 0xE1, 0x18, 0x7D, 0xFF, 0xB9},
			expectedMagneticHeadingValid:        true,
			expectedMagneticHeading:             220.4,
			expectedIndicatedAirspeedValid:      true,
			expectedIndicatedAirspeed:           240,
			expectedMachNumberValid:             true,
			expectedMachNumber:                  0.388,
			expectedBarometricAltitudeRateValid: true,
			expectedBarometricAltitudeRate:      -2080,
			GNSSAltitudeRateValid:               true,
			GNSSAltitudeRate:                    -2272,
		},
		{
			data:                                []byte{0xE6, 0x59, 0x45, 0x0F, 0xFF, 0x0F, 0xE6},
			expectedMagneticHeadingValid:        true,
			expectedMagneticHeading:             287.8,
			expectedIndicatedAirspeedValid:      true,
			expectedIndicatedAirspeed:           162,
			expectedMachNumberValid:             true,
			expectedMachNumber:                  0.252,
			expectedBarometricAltitudeRateValid: true,
			expectedBarometricAltitudeRate:      -992,
			GNSSAltitudeRateValid:               true,
			GNSSAltitudeRate:                    -832,
		},
		{
			data:                                []byte{0xC9, 0x59, 0xDD, 0x20, 0x7E, 0xBF, 0xD1},
			expectedMagneticHeadingValid:        true,
			expectedMagneticHeading:             206.2,
			expectedIndicatedAirspeedValid:      true,
			expectedIndicatedAirspeed:           238,
			expectedMachNumberValid:             true,
			expectedMachNumber:                  0.516,
			expectedBarometricAltitudeRateValid: true,
			expectedBarometricAltitudeRate:      -1312,
			GNSSAltitudeRateValid:               true,
			GNSSAltitudeRate:                    -1504,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		frame, err := DecodeBDS60(testData.data)
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		assert.NoError(err, testMsg+"decodeBDS60 error")
		assert.Equal(testData.expectedMagneticHeadingValid, frame.MagneticHeadingValid, testMsg+"MagneticHeadingValid")
		assert.Equal(testData.expectedMagneticHeading, math.Round(frame.MagneticHeading*10)/10, testMsg+"MagneticHeading")
		assert.Equal(testData.expectedIndicatedAirspeedValid, frame.IndicatedAirspeedValid, testMsg+"IndicatedAirspeedValid")
		assert.Equal(testData.expectedIndicatedAirspeed, frame.IndicatedAirspeed, testMsg+"IndicatedAirspeed")
		assert.Equal(testData.expectedMachNumberValid, frame.MachNumberValid, testMsg+"MachNumberValid")
		assert.Equal(testData.expectedMachNumber, frame.MachNumber, testMsg+"MachNumber")
		assert.Equal(testData.expectedBarometricAltitudeRateValid, frame.BarometricAltitudeRateValid, testMsg+"BarometricAltitudeRateValid")
		assert.Equal(testData.expectedBarometricAltitudeRate, frame.BarometricAltitudeRate, testMsg+"BarometricAltitudeRate")
		assert.Equal(testData.GNSSAltitudeRateValid, frame.GNSSAltitudeRateValid, testMsg+"GNSSAltitudeRateValid")
		assert.Equal(testData.GNSSAltitudeRate, frame.GNSSAltitudeRate, testMsg+"GNSSAltitudeRate")
	}
}

// func TestDecodeBDS60indicatedAirspeed(t *testing.T) {
// 	data := []byte{0xA8, 0x00, 0x04, 0xAA, 0xA7, 0x4A, 0x07, 0x2B, 0xFD, 0xEF, 0xC1, 0xD5, 0xCB, 0x4F}
// 	df := getDF(data)
// 	assert.Equal(t, DF21, df)
// 	msg := decodeDF21(data)
// 	ias, err := decodeBDS60indicatedAirspeed(msg.mb)
// 	assert.NoError(t, err)
// 	assert.Equal(t, float64(259), ias)
// }

// func TestDecodeBDS60machNumber(t *testing.T) {
// 	data := []byte{0xA8, 0x00, 0x04, 0xAA, 0xA7, 0x4A, 0x07, 0x2B, 0xFD, 0xEF, 0xC1, 0xD5, 0xCB, 0x4F}
// 	df := getDF(data)
// 	assert.Equal(t, DF21, df)
// 	msg := decodeDF21(data)
// 	mach, err := decodeBDS60machNumber(msg.mb)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 0.7, math.Round(mach*1000)/1000)
// }

// func TestDecodeBDS60barometricAltitudeRate(t *testing.T) {
// 	data := []byte{0xA8, 0x00, 0x04, 0xAA, 0xA7, 0x4A, 0x07, 0x2B, 0xFD, 0xEF, 0xC1, 0xD5, 0xCB, 0x4F}
// 	df := getDF(data)
// 	assert.Equal(t, DF21, df)
// 	msg := decodeDF21(data)
// 	bar, err := decodeBDS60barometricAltitudeRate(msg.mb)
// 	assert.NoError(t, err)
// 	assert.Equal(t, float64(-2144), bar)
// }

// func TestDecodeBDS60GNSSAltitudeRate(t *testing.T) {
// 	data := []byte{0xA8, 0x00, 0x04, 0xAA, 0xA7, 0x4A, 0x07, 0x2B, 0xFD, 0xEF, 0xC1, 0xD5, 0xCB, 0x4F}
// 	df := getDF(data)
// 	assert.Equal(t, DF21, df)
// 	msg := decodeDF21(data)
// 	ivv, err := decodeBDS60GNSSAltitudeRate(msg.mb)
// 	assert.NoError(t, err)
// 	assert.Equal(t, float64(-2016), ivv)
// }
