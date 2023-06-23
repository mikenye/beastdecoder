package bds

// Aircraft operation status
// Type Code (TC): 31
// https://mode-s.org/decode/content/ads-b/6-operation-status.html

import (
	"errors"
)

type BDS07Frame struct {
	tc int // Type code
	St int // Sub-type code

	om int // Operational mode codes

	// 1. ACAS Resolution Advisory (RA) active
	//     0 = ACAS II or ACAS RA not active
	//     1 = ACAS RA is active
	// 2. IDENT switch active
	//     0 = Ident switch not active
	//     1 = Ident switch active — retained for 18 ±1 seconds
	// 3. Receiving ATC services
	//     0 = Aircraft not receiving ATC services
	//     1 = Aircraft receiving ATC services

	Ver int // ADS-B version number

	Version0Data BDS07FrameVersion0 // Version 0 data
	Version1Data BDS07FrameVersion1 // Version 1 data
	Version2Data BDS07FrameVersion2 // Version 2 data
}

type BDS07FrameVersion2 struct {
	// Aircraft operational status (Version 2)

	cc                   int                                    // Capability class codes
	airborneCapabilities BDS07FrameVersion2AirborneCapabilities // Capability class codes decoded to bools
	surfaceCapabilities  BDS07FrameVersion2SurfaceCapabilities  // Capability class codes decoded to bools

	om                       int                                        // Operational mode codes
	airborneOperationalModes BDS07FrameVersion2AirborneOperationalModes // Operational mode codes decoded to bools
	surfaceOperationalModes  BDS07FrameVersion2SurfaceOperationalModes  // Operational mode codes decoded to bools

	nica int // NIC supplement - A
	nacp int // Navigational accuracy category - position
	gva  int // Geometric vertical accuracy (depending on st bits)
	sil  int // Source integrity level
	bai  int // Barometric altitude integrity (depending on st bits) (NIC-baro)
	sils int // SIL supplement

	Airborne bool
	Hrd      int // Horizontal reference direction
	Trk      int // Track angle or heading (depending on st bits)
}

type BDS07FrameVersion2AirborneCapabilities struct {
	TcasAcasOperational                bool // TCAS/ACAS Operational
	ExtendedSquitter1090MhzReceive     bool // 1090 ES IN (1090 MHz Extended Squitter Receive capability)
	AirReferencedVelocityReport        bool // Capability of sending messages to support Air-Referenced Velocity (ARV) Reports
	TargetStateReport                  bool // Capability of sending messages to support Air-Referenced Velocity Reports
	SupportTargetChangePlus0ReportOnly bool // Capability of sending messages to support TC + 0 Report only
	SupportMultipleTargetChangeReports bool // Capability of sending information for multiple TC reports
	UniversalAccessTransceiver         bool // Aircraft has UAT Receive capability
}

type BDS07FrameVersion2SurfaceCapabilities struct {
	ExtendedSquitter1090MhzReceive bool // 1090 ES IN (1090 MHz Extended Squitter Receive capability)
	ClassB2TransmitPower           bool // B2 Low (Class B2 Transmit Power Less Than 70 Watts)
	UniversalAccessTransceiver     bool // Aircraft has UAT Receive capability
	NACv                           int  // NACV (Navigation Accuracy Category for Velocity)
	NicSupplementC                 int  // NIC Supplement-C (NIC Supplement for use on the Surface)
}

type BDS07FrameVersion2AirborneOperationalModes struct {
	TcasAcasResolutionAdvisoryActive bool // TCAS/ACAS Resolution Advisory (RA) Active
	IdentSwitchActive                bool // IDENT Switch Active
	SingleAntennaFlag                bool // Single Antenna Flag (SAF)
	SystemDesignAssurance            int  // System Design Assurance (SDA)
}

type BDS07FrameVersion2SurfaceOperationalModes struct {
	TcasAcasResolutionAdvisoryActive bool // TCAS/ACAS Resolution Advisory (RA) Active
	IdentSwitchActive                bool // IDENT Switch Active
	SingleAntennaFlag                bool // Single Antenna Flag (SAF)
	SystemDesignAssurance            int  // System Design Assurance (SDA)
	GPSAntennaOffset                 int  // GPS Antenna Offset
}

type BDS07FrameVersion1 struct {
	// Aircraft operational status (Version 1)

	cc                   int                                    // Capability class codes
	airborneCapabilities BDS07FrameVersion1AirborneCapabilities // Capability class codes decoded to bools
	surfaceCapabilities  BDS07FrameVersion1SurfaceCapabilities  // Capability class codes decoded to bools

	om               int                                // Operational mode codes
	operationalModes BDS07FrameVersion1OperationalModes // Operational mode codes decoded to bools

	nics int // NIC supplement
	nacp int // Navigational accuracy category - position
	sil  int // Surveillance integrity level
	baq  int // Barometric altitude quality (airborne only)
	bai  int // Barometric altitude integrity (airborne only)

	Airborne bool
	Hrd      int // Horizontal reference direction
	Trk      int // Track angle or heading (surface only)
}

type BDS07FrameVersion1AirborneCapabilities struct {
	AcasOperationalOrUnknown           bool // Not-ACAS (Airborne Collision Avoidance System Status)
	CockpitDisplayOfTrafficInformation bool // CDTI (Cockpit Display of Traffic Information Status)
	AirReferencedVelocityReport        bool // ARV (Air-Referenced Velocity Report Capability)
	TargetStateReport                  bool // TS (Target State Report Capability)
	SupportTargetChangePlus0ReportOnly bool // Capability of sending messages to support TC + 0 Report only
	SupportMultipleTargetChangeReports bool // Capability of sending information for multiple TC reports
}

type BDS07FrameVersion1SurfaceCapabilities struct {
	CockpitDisplayOfTrafficInformation bool // CDTI (Cockpit Display of Traffic Information Status)
	PositionOffsetApplied              bool // POA (Position Offset Applied)
	ClassB2TransmitPower               bool // B2 Low (Class B2 Transmit Power Less Than 70 Watts)
}

type BDS07FrameVersion1OperationalModes struct {
	AcasResolutionAdvisoryActive bool // ACAS Resolution Advisory (RA) Active
	IdentSwitchActive            bool // IDENT Switch Active
	ReceivingATCServices         bool // Receiving ATC services
}

type BDS07FrameVersion0 struct {
	cc4 int
	cc3 int
	cc2 int
	cc1 int
	om4 int
	om3 int
	om2 int
	om1 int
}

func inferBDS07FrameVersion(mb []byte) (ver int, err error) {

	// get st & ver bits
	st := int(mb[0]) & 0b00000111
	ver = (int(mb[5]) & 0b11100000) >> 5

	// test for version 0
	// st should be 000
	// bits 41-56 reserved (0)
	if st == 0 && mb[5] == 0 && mb[6] == 0 {
		ver = 0
		return
	}

	// test for version 1
	if ver == 1 && (int(mb[6])&0b00000011) == 0 {
		if (int(mb[6]) & 0b00000011) != 0 {
			err = errors.New("reserved bits not zero")
			return
		}
		ver = 1
		return
	}

	// test for version 2
	if ver == 2 && (int(mb[6])&0b00000001) == 0 {
		if (int(mb[6]) & 0b00000001) != 0 {
			err = errors.New("reserved bits not zero")
			return
		}
		ver = 2
		return
	}

	err = errors.New("could not determine version")
	return
}

func decodeBDS07Version0(mb []byte) (frame BDS07Frame) {
	frame.tc = (int(mb[0]) & 0b11111000) >> 3
	frame.St = int(mb[0]) & 0b00000111
	frame.Version0Data.cc4 = (int(mb[1]) & 0b11110000) >> 4
	frame.Version0Data.cc3 = int(mb[1]) & 0b00001111
	frame.Version0Data.cc2 = (int(mb[2]) & 0b11110000) >> 4
	frame.Version0Data.cc1 = int(mb[2]) & 0b00001111
	frame.Version0Data.om4 = (int(mb[3]) & 0b11110000) >> 4
	frame.Version0Data.om3 = int(mb[3]) & 0b00001111
	frame.Version0Data.om2 = (int(mb[4]) & 0b11110000) >> 4
	frame.Version0Data.om1 = int(mb[4]) & 0b00001111
	return
}

func decodeBDS07Version1(mb []byte) (frame BDS07Frame) {
	frame.tc = (int(mb[0]) & 0b11111000) >> 3
	frame.St = int(mb[0]) & 0b00000111

	frame.Version1Data.cc = (int(mb[1]) << 8) + int(mb[2])

	// See B.2.3.10.3 "Capability Class (CC) Codes"
	// http://www.aviationchief.com/uploads/9/2/0/9/92098238/icao_doc_9871_-_technical_provisions_for_mode_s_-_advanced_edition_1.pdf

	switch frame.St {
	case 0: // airborne

		frame.Version1Data.Airborne = true

		// 1. Not-ACAS (Airborne Collision Avoidance System Status)
		//     0 = ACAS operational or unknown
		//     1 = ACAS not installed or not operational
		if (int(mb[1])&0b00100000)>>5 == 0 {
			frame.Version1Data.airborneCapabilities.AcasOperationalOrUnknown = true
		} else {
			frame.Version1Data.airborneCapabilities.AcasOperationalOrUnknown = false
		}

		// 2. CDTI (Cockpit Display of Traffic Information Status)
		//     0 = Traffic display not operational
		//     1 = Traffic display operational
		if (int(mb[1])&0b00010000)>>4 == 1 {
			frame.Version1Data.airborneCapabilities.CockpitDisplayOfTrafficInformation = true
		}

		// 3. ARV (Air-Referenced Velocity Report Capability)
		//     0 = No capability for sending messages to support Air-Referenced Velocity Reports
		//     1 = Capability of sending messages to support Air-Referenced Velocity Reports
		if (int(mb[1])&0b00000010)>>1 == 1 {
			frame.Version1Data.airborneCapabilities.AirReferencedVelocityReport = true
		}

		// 4. TS (Target State Report Capability)
		//     0 = No capability for sending messages to support Target State Reports
		//     1 = Capability of sending messages to support Target State Reports
		if (int(mb[1]) & 0b00000001) == 1 {
			frame.Version1Data.airborneCapabilities.TargetStateReport = true
		}

		// 5. TC (Target Change Report Capability)
		//     0 = No capability for sending messages to support Trajectory Change Reports
		//     1 = Capability of sending messages to support TC+0 Report only
		//     2 = Capability of sending information for multiple TC Reports
		//     3 = Reserved
		switch (int(mb[2]) & 0b11000000) >> 6 {
		case 1:
			frame.Version1Data.airborneCapabilities.SupportTargetChangePlus0ReportOnly = true
		case 2:
			frame.Version1Data.airborneCapabilities.SupportMultipleTargetChangeReports = true
		}

	case 1: // surface

		frame.Version1Data.Airborne = false

		// 1. CDTI (Cockpit Display of Traffic Information Status)
		//     0 = Traffic display not operational
		//     1 = Traffic display operational
		if (int(mb[1])&0b00010000)>>4 == 1 {
			frame.Version1Data.surfaceCapabilities.CockpitDisplayOfTrafficInformation = true
		}

		// 2. POA (Position Offset Applied)
		//     0 = Position transmitted is not the ADS-B position reference point
		//     1 = Position transmitted is the ADS-B position reference point
		if (int(mb[1])&0b00100000)>>5 == 1 {
			frame.Version1Data.surfaceCapabilities.PositionOffsetApplied = true
		}

		// 3. B2 Low (Class B2 transmit power less than 70 Watts)
		//     0 = Greater than or equal to 70 Watts transmit power
		//     1 = Less than 70 Watts transmit power
		if (int(mb[1])&0b00000010)>>1 == 1 {
			frame.Version1Data.surfaceCapabilities.ClassB2TransmitPower = true
		}
	}

	frame.Version1Data.om = (int(mb[3]) << 8) + int(mb[4])

	// 1. ACAS Resolution Advisory (RA) active
	//     0 = ACAS II or ACAS RA not active
	//     1 = ACAS RA is active
	if (int(mb[3])&0b00100000)>>5 == 1 {
		frame.Version1Data.operationalModes.AcasResolutionAdvisoryActive = true
	}

	// 2. IDENT switch active
	//     0 = Ident switch not active
	//     1 = Ident switch active — retained for 18 ±1 seconds
	if (int(mb[3])&0b00010000)>>4 == 1 {
		frame.Version1Data.operationalModes.IdentSwitchActive = true
	}

	// 3. Receiving ATC services
	//     0 = Aircraft not receiving ATC services
	//     1 = Aircraft receiving ATC services
	if (int(mb[3])&0b00001000)>>3 == 1 {
		frame.Version1Data.operationalModes.ReceivingATCServices = true
	}

	frame.Version1Data.nics = (int(mb[5]) & 0b00010000) >> 4
	frame.Version1Data.nacp = (int(mb[5]) & 0b00001111)
	frame.Version1Data.sil = ((int(mb[6]) & 0b00110000) >> 4)
	frame.Version1Data.Hrd = ((int(mb[6]) & 0b00000100) >> 2)

	if frame.St == 0 {
		frame.Version1Data.baq = ((int(mb[6]) & 0b11000000) >> 6)
		frame.Version1Data.bai = ((int(mb[6]) & 0b00001000) >> 3)
	} else {
		frame.Version1Data.Trk = ((int(mb[6]) & 0b00001000) >> 3)
	}
	return frame
}

func decodeBDS07Version2(mb []byte) (frame BDS07Frame) {
	frame.tc = (int(mb[0]) & 0b11111000) >> 3
	frame.St = int(mb[0]) & 0b00000111
	frame.Version2Data.cc = (int(mb[1]) << 8) + int(mb[2])

	switch frame.St {
	case 0: // airborne

		frame.Version2Data.Airborne = true

		// http://www.icscc.org.cn/upload/file/20190102/Doc.9871-EN%20Technical%20Provisions%20for%20Mode%20S%20Services%20and%20Extended%20Squitter.pdf
		// Airborne Capability Class (CC) Code for Version 2

		// 1. TCAS/ACAS Operational
		//     = 0: TCAS/ACAS is NOT Operational
		//     = 1: TCAS/ACAS IS Operational
		if (int(mb[1])&0b00100000)>>5 == 1 {
			frame.Version2Data.airborneCapabilities.TcasAcasOperational = true
		}

		// 2. 1090 ES IN (1090 MHz Extended Squitter)
		//     = 0: Aircraft has NO 1 090 ES Receive capability
		//     = 1: Aircraft has 1 090 ES Receive capability
		if (int(mb[1])&0b00010000)>>4 == 1 {
			frame.Version2Data.airborneCapabilities.ExtendedSquitter1090MhzReceive = true
		}

		// 3. ARV (Air-Referenced Velocity Report Capability)
		//     = 0: No capability for sending messages to support Air Referenced Velocity Reports
		//     = 1: Capability of sending messages to support Air-Referenced Velocity Reports
		if (int(mb[1])&0b00000010)>>1 == 1 {
			frame.Version2Data.airborneCapabilities.AirReferencedVelocityReport = true
		}

		// 4. TS (Target State Report Capability)
		//     = 0: No capability for sending messages to support Target State Reports
		//     = 1: Capability of sending messages to support Target State Reports
		if (int(mb[1]) & 0b00000001) == 1 {
			frame.Version2Data.airborneCapabilities.TargetStateReport = true
		}

		// 5. TC (Target Change Report Capability)
		//     = 0: No capability for sending messages to support Trajectory Change Reports
		//     = 1: Capability of sending messages to support TC + 0 Report only
		//     = 2: Capability of sending information for multiple TC reports
		//     = 3: Reserved
		switch (int(mb[2]) & 0b11000000) >> 6 {
		case 1:
			frame.Version2Data.airborneCapabilities.SupportTargetChangePlus0ReportOnly = true
		case 2:
			frame.Version2Data.airborneCapabilities.SupportMultipleTargetChangeReports = true
		}

		// 6. UAT IN (Universal Access Transceiver)
		//     = 0: Aircraft has No UAT Receive capability
		//     = 1: Aircraft has UAT Receive capability
		if (int(mb[2])&0b00100000)>>5 == 1 {
			frame.Version2Data.airborneCapabilities.UniversalAccessTransceiver = true
		}

	case 1: // surface

		frame.Version2Data.Airborne = false

		// http://www.icscc.org.cn/upload/file/20190102/Doc.9871-EN%20Technical%20Provisions%20for%20Mode%20S%20Services%20and%20Extended%20Squitter.pdf
		// Surface Capability Class (CC) Code for Version 2 Systems

		// 1. 1090 ES IN (1090 MHz Extended Squitter)
		//     = 0: Aircraft has NO 1 090 ES Receive capability
		//     = 1: Aircraft has 1 090 ES Receive capability
		if (int(mb[1])&0b00010000)>>4 == 1 {
			frame.Version2Data.surfaceCapabilities.ExtendedSquitter1090MhzReceive = true
		}

		// 2. B2 Low (Class B2 Transmit Power Less Than 70 Watts)
		//     = 0: Greater than or equal to 70 Watts Transmit Power
		//     = 1: Less than 70 Watts Transmit Power
		if (int(mb[1])&0b00000010)>>1 == 1 {
			frame.Version2Data.surfaceCapabilities.ClassB2TransmitPower = true
		}

		// 3. UAT IN (Universal Access Transceiver)
		//     = 0: Aircraft has NO UAT Receive capability
		//     = 1: Aircraft has UAT Receive capability
		if (int(mb[1]) & 0b00000001) == 1 {
			frame.Version2Data.surfaceCapabilities.UniversalAccessTransceiver = true
		}

		// 4. NACV (Navigation Accuracy Category for Velocity)
		frame.Version2Data.surfaceCapabilities.NACv = (int(mb[2]) & 0b11100000) >> 5

		// 5. NIC Supplement-C (NIC Supplement for use on the Surface)
		frame.Version2Data.surfaceCapabilities.NicSupplementC = (int(mb[2]) & 0b00010000) >> 4

	}

	frame.Version2Data.om = (int(mb[3]) << 8) + int(mb[4])

	switch frame.St {
	case 0: // airborne

		// C.2.3.10.3 in http://www.icscc.org.cn/upload/file/20190102/Doc.9871-EN%20Technical%20Provisions%20for%20Mode%20S%20Services%20and%20Extended%20Squitter.pdf
		// Airborne Operational Mode (OM) Subfield Format

		// 1. TCAS/ACAS Resolution Advisory (RA) Active
		//     = 0: TCAS II or ACAS RA not active
		//     = 1: TCAS/ACAS RA is active
		if (int(mb[3])&0b00100000)>>5 == 1 {
			frame.Version2Data.airborneOperationalModes.TcasAcasResolutionAdvisoryActive = true
		}

		// 2. IDENT Switch Active
		//     = 0: Ident switch not active
		//     = 1: Ident switch active – retained for 18 ±1 seconds
		if (int(mb[3])&0b00010000)>>4 == 1 {
			frame.Version2Data.airborneOperationalModes.IdentSwitchActive = true
		}

		// 3. Reserved for Receiving ATC Services
		//     = 0: Set to ZERO for this edition of this manual
		if (int(mb[3])&0b00000100)>>2 == 1 {
			frame.Version2Data.airborneOperationalModes.SingleAntennaFlag = true
		}

		// 4. Single Antenna Flag (SAF)
		//     = 0: Systems with two functioning antennas
		//     = 1: Systems that use only one antenna
		frame.Version2Data.airborneOperationalModes.SystemDesignAssurance = (int(mb[3]) & 0b00000011)

		// 5. System Design Assurance (SDA)
		// TODO

	case 1: // surface

		// C.2.3.10.3 in http://www.icscc.org.cn/upload/file/20190102/Doc.9871-EN%20Technical%20Provisions%20for%20Mode%20S%20Services%20and%20Extended%20Squitter.pdf
		// Surface Operational Mode (OM) Subfield Format

		// 1. TCAS/ACAS Resolution Advisory (RA) Active
		//     = 0: TCAS II or ACAS RA not active
		//     = 1: TCAS/ACAS RA is active
		if (int(mb[3])&0b00100000)>>5 == 1 {
			frame.Version2Data.surfaceOperationalModes.TcasAcasResolutionAdvisoryActive = true
		}

		// 2. IDENT Switch Active
		//    = 0: Ident switch not active
		//    = 1: Ident switch active – retained for 18 ±1 seconds
		if (int(mb[3])&0b00010000)>>4 == 1 {
			frame.Version2Data.surfaceOperationalModes.IdentSwitchActive = true
		}

		// 3. Reserved for Receiving ATC Services
		//    = 0: Set to ZERO for this edition of this manual
		if (int(mb[3])&0b00000100)>>2 == 1 {
			frame.Version2Data.surfaceOperationalModes.SingleAntennaFlag = true
		}

		// 4. Single Antenna Flag (SAF)
		//    = 0: Systems with two functioning antennas
		//    = 1: Systems that use only one antenna
		frame.Version2Data.surfaceOperationalModes.SystemDesignAssurance = (int(mb[3]) & 0b00000011)

		// 5. System Design Assurance (SDA)
		frame.Version2Data.surfaceOperationalModes.GPSAntennaOffset = int(mb[4])

		// 6. GPS Antenna Offset
		// TODO

	}

	frame.Version2Data.nica = (int(mb[5]) & 0b00010000) >> 4
	frame.Version2Data.nacp = (int(mb[5]) & 0b00001111)
	frame.Version2Data.sil = (int(mb[6]) & 0b00110000) >> 4
	frame.Version2Data.Hrd = (int(mb[6]) & 0b00000100) >> 2 // 0= True heading, 1= Magnetic heading?
	frame.Version2Data.sils = (int(mb[6]) & 0b00000010) >> 1

	if frame.St == 0 {
		frame.Version2Data.gva = (int(mb[6]) & 0b11000000) >> 6
		frame.Version2Data.bai = (int(mb[6]) & 0b00001000) >> 3
	} else {
		frame.Version2Data.Trk = (int(mb[6]) & 0b00001000) >> 3
	}
	return
}

func DecodeBDS07(mb []byte) (frame BDS07Frame, err error) {

	// check type code is 31
	tc := (int(mb[0]) & 0b11111000) >> 3
	if tc != 0b11111 {
		err = errors.New("type code not 31")
		return
	}

	ver, err := inferBDS07FrameVersion(mb)
	if err != nil {
		return
	}

	switch ver {
	case 0:
		frame = decodeBDS07Version0(mb)
	case 1:
		frame = decodeBDS07Version1(mb)
	case 2:
		frame = decodeBDS07Version2(mb)
	}

	frame.Ver = ver

	return
}
