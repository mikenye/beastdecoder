package bds

import "errors"

// Selected vertical intention (BDS 4,0)

type BDS40Frame struct {
	McpFcuSelectedAltitudeValid bool // Status (for MCP/FCU selected altitude)
	McpFcuSelectedAltitude      int  // MCP/FCU selected altitude

	FmsSelectedAltitudeValid bool // Status (for FMS selected altitude)
	FmsSelectedAltitude      int  // FMS selected altitude

	BarometricPressureSettingValid bool    // Status (for barometric press setting)
	BarometricPressureSetting      float64 // Barometric pressure setting

	McpFcuModeValid bool // Status of MCP/FCU mode
	McpFcuMode      BDS40FrameMcpFcpMode

	TargetAltitudeSourceValid bool // Status of target altitude source
	TargetAltitudeSource      int  // Target altitude source
}

type BDS40FrameMcpFcpMode struct {
	VnavMode     bool // VNAV mode
	AltHoldMode  bool // Alt hold mode
	ApproachMode bool // Approach mode
}

func DecodeBDS40(mb []byte) (frame BDS40Frame, err error) {
	// decode Selected vertical intention (BDS 4,0)
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#selected-vertical-intention-bds-40

	// MCP/FCU selected altitude
	switch (int(mb[0]) & 0b10000000) >> 7 {
	case 0:
		frame.McpFcuSelectedAltitudeValid = false
	case 1:
		frame.McpFcuSelectedAltitudeValid = true
	default:
		err = errors.New("Invalid Status (for MCP/FCU selected altitude)")
		return
	}
	if frame.McpFcuSelectedAltitudeValid {
		frame.McpFcuSelectedAltitude = (((int(mb[0]) & 0b01111111) << 5) + ((int(mb[1]) & 0b11111000) >> 3)) * 16
	}

	// FMS selected altitude
	switch (int(mb[1]) & 0b00000100) >> 2 {
	case 0:
		frame.FmsSelectedAltitudeValid = false
	case 1:
		frame.FmsSelectedAltitudeValid = true
	default:
		err = errors.New("Invalid Status (for FMS selected altitude)")
		return
	}
	if frame.FmsSelectedAltitudeValid {
		frame.FmsSelectedAltitude = ((int(mb[1]) & 0b00000011) << 10) + ((int(mb[2])) << 2) + ((int(mb[3]) & 0b11000000) >> 6)
	}

	// Barometric pressure setting
	switch (int(mb[3]) & 0b00100000) >> 5 {
	case 0:
		frame.BarometricPressureSettingValid = false
	case 1:
		frame.BarometricPressureSettingValid = true
	default:
		err = errors.New("Invalid Status (for barometric pressure setting)")
		return
	}
	if frame.BarometricPressureSettingValid {
		frame.BarometricPressureSetting = (float64(((int(mb[3])&0b00011111)<<7)+((int(mb[4])&0b11111110)>>1)) * 0.1) + 800
	}

	// MCP/FCU mode
	switch int(mb[5]) & 0b00000001 {
	case 0:
		frame.McpFcuModeValid = false
	case 1:
		frame.McpFcuModeValid = true
	default:
		err = errors.New("Invalid Status (for MCP/FCU mode)")
		return
	}
	if frame.McpFcuModeValid {
		frame.McpFcuMode = BDS40FrameMcpFcpMode{}

		// VNAV mode
		switch (int(mb[6]) & 0b10000000) >> 7 {
		case 0:
			frame.McpFcuMode.VnavMode = false
		case 1:
			frame.McpFcuMode.VnavMode = true
		default:
			err = errors.New("Invalid data for VNAV mode")
			return
		}

		// Alt hold mode
		switch (int(mb[6]) & 0b01000000) >> 6 {
		case 0:
			frame.McpFcuMode.AltHoldMode = false
		case 1:
			frame.McpFcuMode.AltHoldMode = true
		default:
			err = errors.New("Invalid data for Alt hold mode")
			return
		}

		// Approach mode
		switch (int(mb[6]) & 0b00100000) >> 5 {
		case 0:
			frame.McpFcuMode.ApproachMode = false
		case 1:
			frame.McpFcuMode.ApproachMode = true
		default:
			err = errors.New("Invalid data for Approach mode")
			return
		}
	}

	// Target altitude source
	switch (int(mb[6]) & 0b00000100) >> 2 {
	case 0:
		frame.TargetAltitudeSourceValid = false
	case 1:
		frame.TargetAltitudeSourceValid = true
	default:
		err = errors.New("Invalid data for Target altitude source")
		return
	}
	if frame.TargetAltitudeSourceValid {
		frame.TargetAltitudeSource = int(mb[6]) & 0b00000011
	}

	return
}
