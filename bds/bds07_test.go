package bds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS07(t *testing.T) {

	// define test data
	var testTable = []struct {
		data         []byte
		tc           int                // Type code
		st           int                // Sub-type code
		ver          int                // Version
		version0Data BDS07FrameVersion0 // Version 0 data
		version1Data BDS07FrameVersion1 // Version 1 data
		version2Data BDS07FrameVersion2 // Version 2 data
	}{
		{
			data: []byte{0xF8, 0x10, 0x20, 0x06, 0x00, 0x49, 0xB8},
			tc:   31,
			ver:  2,
			version2Data: BDS07FrameVersion2{
				airborneCapabilities: BDS07FrameVersion2AirborneCapabilities{
					ExtendedSquitter1090MhzReceive: true,
					UniversalAccessTransceiver:     true,
				},
				airborneOperationalModes: BDS07FrameVersion2AirborneOperationalModes{
					SingleAntennaFlag: true,
				},
				Hrd:  0, // true heading
				nica: 0,
				bai:  1,
				nacp: 9,
				gva:  2,
				sil:  3,
			},
		},
		{
			data: []byte{0xF8, 0x00, 0x00, 0x00, 0x00, 0x29, 0x28},
			tc:   31,
			ver:  1,
			version1Data: BDS07FrameVersion1{
				airborneCapabilities: BDS07FrameVersion1AirborneCapabilities{
					AcasOperationalOrUnknown: true,
				},
				Hrd:  0, // true heading
				nics: 0, // nica?
				bai:  1,
				nacp: 9,
				sil:  2,
			},
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		frame, err := DecodeBDS07(testData.data)
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		assert.NoError(err, testMsg+"decodeBDS07 error")
		assert.Equal(testData.tc, frame.tc, testMsg+"tc")
		assert.Equal(testData.ver, frame.Ver, testMsg+"ver")

		// version 1 frames
		assert.Equal(testData.version1Data.airborneCapabilities.AcasOperationalOrUnknown, frame.Version1Data.airborneCapabilities.AcasOperationalOrUnknown, testMsg+"version1Data.airborneCapabilities.AcasOperationalOrUnknown")
		assert.Equal(testData.version1Data.Hrd, frame.Version1Data.Hrd, testMsg+"version1Data.hrd")
		assert.Equal(testData.version1Data.nics, frame.Version1Data.nics, testMsg+"version1Data.nics")
		assert.Equal(testData.version1Data.bai, frame.Version1Data.bai, testMsg+"version1Data.bai")
		assert.Equal(testData.version1Data.nacp, frame.Version1Data.nacp, testMsg+"version1Data.nacp")
		assert.Equal(testData.version1Data.sil, frame.Version1Data.sil, testMsg+"version1Data.sil")

		// version 2 frames
		assert.Equal(testData.version2Data.airborneCapabilities.ExtendedSquitter1090MhzReceive, frame.Version2Data.airborneCapabilities.ExtendedSquitter1090MhzReceive, testMsg+"version2Data.airborneCapabilities.ExtendedSquitter1090MhzReceive")
		assert.Equal(testData.version2Data.airborneCapabilities.UniversalAccessTransceiver, frame.Version2Data.airborneCapabilities.UniversalAccessTransceiver, testMsg+"version2Data.airborneCapabilities.UniversalAccessTransceiver")
		assert.Equal(testData.version2Data.airborneOperationalModes.SingleAntennaFlag, frame.Version2Data.airborneOperationalModes.SingleAntennaFlag, testMsg+"version2Data.airborneOperationalModes.SingleAntennaFlag")
		assert.Equal(testData.version2Data.Hrd, frame.Version2Data.Hrd, testMsg+"version2Data.hrd")
		assert.Equal(testData.version2Data.nica, frame.Version2Data.nica, testMsg+"version2Data.nica")
		assert.Equal(testData.version2Data.bai, frame.Version2Data.bai, testMsg+"version2Data.bai")
		assert.Equal(testData.version2Data.nacp, frame.Version2Data.nacp, testMsg+"version2Data.nacp")
		assert.Equal(testData.version2Data.gva, frame.Version2Data.gva, testMsg+"version2Data.gva")
		assert.Equal(testData.version2Data.sil, frame.Version2Data.sil, testMsg+"version2Data.sil")
	}

}
