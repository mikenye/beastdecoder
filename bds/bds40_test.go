package bds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBD40(t *testing.T) {
	// define test data
	var testTable = []struct {
		data                                   []byte
		expectedMcpFcuSelectedAltitudeValid    bool
		expectedMcpFcuSelectedAltitude         int
		expectedFmsSelectedAltitudeValid       bool
		expectedFmsSelectedAltitude            int
		expectedBarometricPressureSettingValid bool
		expectedBarometricPressureSetting      float64
		expectedMcpFcuModeValid                bool
		expectedMcpFcuMode                     BDS40FrameMcpFcpMode
		expectedTargetAltitudeSourceValid      bool
		expectedTargetAltitudeSource           int
	}{
		{
			data:                                   []byte{0x83, 0xE8, 0x00, 0x31, 0x08, 0x00, 0x00},
			expectedMcpFcuSelectedAltitudeValid:    true,
			expectedMcpFcuSelectedAltitude:         2000,
			expectedFmsSelectedAltitudeValid:       false,
			expectedFmsSelectedAltitude:            0,
			expectedBarometricPressureSettingValid: true,
			expectedBarometricPressureSetting:      1018.0,
			expectedMcpFcuModeValid:                false,
			expectedMcpFcuMode:                     BDS40FrameMcpFcpMode{},
			expectedTargetAltitudeSourceValid:      false,
			expectedTargetAltitudeSource:           0,
		},
		{
			data:                                   []byte{0x8D, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectedMcpFcuSelectedAltitudeValid:    true,
			expectedMcpFcuSelectedAltitude:         7136,
			expectedFmsSelectedAltitudeValid:       false,
			expectedFmsSelectedAltitude:            0,
			expectedBarometricPressureSettingValid: false,
			expectedBarometricPressureSetting:      0,
			expectedMcpFcuModeValid:                false,
			expectedMcpFcuMode:                     BDS40FrameMcpFcpMode{},
			expectedTargetAltitudeSourceValid:      true,
			expectedTargetAltitudeSource:           550,
		},
		{
			data:                                   []byte{0x89, 0xC8, 0x00, 0x31, 0x08, 0x01, 0x80},
			expectedMcpFcuSelectedAltitudeValid:    true,
			expectedMcpFcuSelectedAltitude:         5008,
			expectedFmsSelectedAltitudeValid:       false,
			expectedFmsSelectedAltitude:            0,
			expectedBarometricPressureSettingValid: true,
			expectedBarometricPressureSetting:      1018.0,
			expectedMcpFcuModeValid:                true,
			expectedMcpFcuMode: BDS40FrameMcpFcpMode{
				VnavMode: true,
			},
			expectedTargetAltitudeSourceValid: false,
			expectedTargetAltitudeSource:      0,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		frame, err := DecodeBDS40(testData.data)
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		assert.NoError(err, testMsg+"decodeBDS40 error")
		assert.Equal(testData.expectedMcpFcuSelectedAltitudeValid, frame.McpFcuSelectedAltitudeValid, testMsg+"McpFcuSelectedAltitudeValid")
		assert.Equal(testData.expectedMcpFcuSelectedAltitude, frame.McpFcuSelectedAltitude, testMsg+"McpFcuSelectedAltitude")
		assert.Equal(testData.expectedFmsSelectedAltitudeValid, frame.FmsSelectedAltitudeValid, testMsg+"FmsSelectedAltitudeValid")
		assert.Equal(testData.expectedFmsSelectedAltitude, frame.FmsSelectedAltitude, testMsg+"FmsSelectedAltitude")
		assert.Equal(testData.expectedBarometricPressureSettingValid, frame.BarometricPressureSettingValid, testMsg+"BarometricPressureSettingValid")
		assert.Equal(testData.expectedBarometricPressureSetting, frame.BarometricPressureSetting, testMsg+"BarometricPressureSetting")
		assert.Equal(testData.expectedMcpFcuModeValid, frame.McpFcuModeValid, testMsg+"McpFcuModeValid")
		assert.Equal(testData.expectedMcpFcuMode, frame.McpFcuMode, testMsg+"McpFcuMode")
	}
}
