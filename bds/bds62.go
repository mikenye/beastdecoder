package bds

// TARGET STATE AND STATUS INFORMATION
// Type Code (TC): 29

import (
	"errors"
	"fmt"
)

// http://www.icscc.org.cn/upload/file/20190102/Doc.9871-EN%20Technical%20Provisions%20for%20Mode%20S%20Services%20and%20Extended%20Squitter.pdf
// Table: B-2-98
// Page: B-38

type BDS62Frame struct {

	// Vertical Data Available / Source Indicator:
	//
	// This subfield shall be used to identify whether aircraft vertical state information is available and present as well as the data source for the vertical data when present in the subsequent subfields.
	// Encoding shall be defined as specified below. Any message parameter associated with the vertical target state for which an update has not been received from an on-board data source within the past 5 seconds shall be considered invalid and so indicated in the Vertical Data Available/Source Indicator subfield.
	//
	//  VTSUnavailable = No valid Vertical Target State data is available
	//  VTSAutopilotSelectedValue = Autopilot control panel selected value, such as Mode Control Panel (MCP) or Flight Control Unit (FCU)
	//  VTSHoldingAltitude = Holding altitude
	//  VTSFmsRnavSystem = FMS/RNAV system
	vdaSi BDS62VerticalTargetState

	// Target Altitude Type
	//
	// This one bit (ME bit 10, Message bit 42) subfield shall be used to identify whether the altitude reported
	// in the “Target Altitude” subfield is referenced to mean sea level (MSL) or to a flight level (FL).
	//  TATReferencedToFL = Indicates target altitude referenced to pressure-altitude (FL).
	//  TATReferencedToMSL = Indicates a target altitude referenced to barometric corrected altitude (MSL).
	tat BDS62TargetAltitudeType

	// Target Altitude Compatibility
	//
	// This subfield shall be used to describe the
	// aircraft’s capabilities for providing the data reported in the target altitude subfield.
	// The target altitude capability subfield shall be encoded as specified below.
	//  TACReportingHoldingOnly = Capability for reporting holding altitude only
	//  TACReportingHoldingAutopilotSelected = Capability for reporting either holding altitude or autopilot control panel selected altitude
	//  TACReportingHoldingAutopilotSelectedFmsRnavLevelOff = Capability for reporting either holding altitude, autopilot control panel selected altitude, or any FMS/RNAV level-off altitude
	tac BDS62TargetAltitudeCapability

	// Vertical Mode Indicator
	//
	// This subfield shall be used to indicate whether the
	// target altitude is in the process of being acquired
	// (i.e. aircraft is climbing or descending toward the target altitude) or whether the
	// target altitude has been acquired/being held.
	// The Vertical Mode Indicator subfield shall be encoded as specified in the following table.
	//  VMIUnknown = Unknown mode or information unavailable
	//  VMIAcquiring = “Acquiring” Mode
	//  VMICapturingOrMaintaining = “Capturing” or “Maintaining” Mode
	vmi BDS62VerticalModeIndicator

	// Target Altitude
	//
	// This subfield shall be used to provide the aircraft’s
	// next intended level-off altitude if in a climb or descent, or the aircraft current intended altitude
	// if it is intending to hold its current altitude.
	// The reported target altitude shall be the operational altitude recognized by the aircraft’s guidance system.
	//
	// The target altitude subfield units is in "ft".
	ta int

	// Horizontal Data Available / Source Indicator
	//
	// This subfield shall be used to identify whether the aircraft
	// horizontal state information is available and present as well as the data source for the horizontal
	// target data when present in the subsequent subfields.
	//
	// The horizontal data available/source Indicator subfield shall be encoded as specified as follows.
	//
	// Any message parameter associated with the horizontal target state for which an update has not been received
	// from an on-board data source within the past 5 seconds shall be considered invalid and so indicated in the
	// horizontal data available/source indicator subfield.
	//  HTSUnavailable = No valid horizontal target state data is available
	//  HTSAutopilotSelectedValue = Autopilot control panel selected value, such as Mode Control Panel (MCP) or Flight Control Unit (FCU)
	//  HTSMaintainingHeading = Maintaining current heading or track angle (e.g. autopilot mode select)
	//  HTSFmsRnavSystem = FMS/RNAV system (indicates track angle specified by leg type)
	hdaSi BDS62HorizontalTargetState

	// Target Heading / Track Angle
	//
	// This subfield shall be used to provide the aircraft’s intended
	// (i.e. target or selected) heading or track.
	//
	// The target heading / track angle is in degrees.
	thta int

	// Target Heading / Track Indicator
	//
	// This subfield shall be used to indicate whether a Heading Angle or a
	// track angle is being reported in the target heading/track angle subfield.
	//  THTITargetHeadingAngle = Indicates the Target Heading Angle is being reported.
	//  THTITrackAngle = Indicates the Track Angle is being reported.
	thti BDS62TargetHeadingTrackIndicator

	// Horizontal Mode Indicator
	//
	// This subfield shall be used to indicate whether the target heading/track
	// is being acquired (i.e. lateral transition toward the target direction is in progress) or whether the target heading/track
	// has been acquired and is currently being maintained.
	// The horizontal mode Indicator subfield shall be encoded as specified as follows:
	//  HMIUnknown = Unknown mode or information unavailable
	//  HMIAcquiring = “Acquiring” mode
	//  HMICapturingOrMaintaining = “Capturing” or “Maintaining” mode
	hmi BDS62HorizontalModeIndicator

	// Navigation Accuracy Category for Position (NACp)
	//
	// This subfield shall be used to indicate the navigational accuracy category
	// of the navigation information used as the basis for the aircraft reported position.
	// The NACP subfield shall be encoded as specified as follows.
	// If an update has not been received from an on-board data source for NACP within the past 5 seconds,
	// then the NACP subfield shall be encoded as a value indicating unknown accuracy.
	//  NACpUnknown = EPU ≥ 18.52 km (10 NM) — Unknown accuracy
	//  NACp10NM = EPU < 18.52 km (10 NM) — RNP –10 accuracy
	//  NACp4NM = EPU < 7.408 km (4 NM) — RNP –4 accuracy
	//  NACp2NM = EPU < 3.704 km (2 NM) — RNP –2 accuracy
	//  NACp1NM = EPU < 1 852 m (1NM) — RNP –1 accuracy
	//  NACp0NM5 = EPU < 926 m (0.5 NM) — RNP –0.5 accuracy
	//  NACp0NM3 = EPU < 555.6 m ( 0.3 NM) — RNP –0.3 accuracy
	//  NACp0NM1 = EPU < 185.2 m (0.1 NM) — RNP –0.1 accuracy
	//  NACp0NM05 = EPU < 92.6 m (0.05 NM) — e.g. GPS (with SA)
	//  NACpGPS = EPU *and* VEPU < 45m — e.g. GPS (SA off)
	//  NACpWAAS = EPU *and* VEPU < 15m — e.g. WAAS
	//  NACpLAAS = EPU *and* VEPU < 4m — e.g.LAAS
	nacp BDS62NACp

	// Navigation Integrity Category for Baro (NICbaro)
	//
	// This subfield shall be used to indicate whether or not the barometric pressure-altitude
	// being reported in the airborne position message has been cross-checked against another source of pressure-altitude.
	// The NICbaro subfield shall be encoded as follows.
	// If an update has not been received from an on-board data source for NICbaro within the past 5 seconds,
	// then the NICbaro subfield shall be encoded as a value of NICBaroNotCrossChecked.
	//
	//  NICBaroNotCrossChecked = The barometric altitude that is being reported in the Airborne Position Message is based on a Gilham coded input that has not been cross-checked against another source of pressure-altitude.
	//
	//  NICBaroCrossChecked = The barometric altitude that is being reported in the Airborne Position Message is either based on a Gilham code input that has been cross-checked against another source of pressure-altitude and verified as being consistent or is based on a non-Gilham coded source.
	nicbaro BDS62NICBaro

	// Surveillance Integrity Level
	//
	// This subfield shall be used to define the probability of the integrity
	// containment region described by the NIC subfield being exceeded for the selected position source, including any external
	// signals used by the source.
	// If an update has not been received from an on-board data source for SIL within the past 5 seconds,
	// then the SIL subfield shall be encoded as a value indicating “Unknown.”
	//
	sil int

	// Capability / Mode Codes
	//
	// This subfield shall be used to indicate the current operational status of
	// TCAS/ACAS systems/functions. This subfield shall be encoded as specified as below.
	// If an update has not been received from an on-board data source for a Capability/Mode Code data element
	// within the past 2 seconds, then that data element shall be encoded with a value of TcasAcasOperationalOrUnknown.
	//  TcasAcasOperationalOrUnknown: TCAS/ACAS operational or unknown
	//  TcasAcasNotOperational: TCAS/ACAS not operational
	cmc BDS62CapabilityModeCodes

	// Capability / Mode Codes
	//
	// This subfield shall be used to indicate the current operational status of
	// TCAS/ACAS systems/functions. This subfield shall be encoded as specified as below.
	//  TcasAcasResolutionAdvisoryInactive: No TCAS/ACAS Resolution Advisory active
	//  TcasAcasResolutionAdvisoryActive: TCAS/ACAS Resolution Advisory active
	cmcRa BDS62CapabilityModeCodesRA

	// Emergency / Priority Status
	//
	// This subfield shall be used to provide additional information
	// regarding aircraft status. The Emergency/Priority Status subfield shall be encoded as specified below.
	// If an update has not been received from an on-board data source for the Emergency/Priority Status within the past 5 seconds,
	// then the emergency/priority status subfield shall be encoded with a value indicating no emergency.
	//  EmergencyNone = No emergency
	//  EmergencyGeneral = General emergency
	//  EmergencyLifeguardMedical = Lifeguard/medical emergency
	//  EmergencyMinimumFuel = Minimum fuel
	//  EmergencyNoCommunications = No communications
	//  EmergencyUnlawfulInterference = Unlawful interference
	//  EmergencyDownedAircraft = Downed aircraft
	eps EmergencyPriorityStatus
}

func (frame *BDS62Frame) Sprint() string {
	// Outputs a multi-line string representing the frame. Useful for debugging.
	var output string

	// Header
	output += "BDS 6,2: Target State and Status (29)\n"

	// VERTICAL DATA AVAILABLE/SOURCE INDICATOR
	output += "  Vertical Data Available/Source Indicator: "
	switch frame.vdaSi {
	case VTSUnavailable:
		output += "No valid Vertical Target State data is available\n"
	case VTSAutopilotSelectedValue:
		output += "Autopilot control panel selected value, such as Mode Control Panel (MCP) or Flight Control Unit (FCU)\n"
	case VTSHoldingAltitude:
		output += "Holding altitude\n"
	case VTSFmsRnavSystem:
		output += "FMS/RNAV system\n"
	}

	// VERTICAL MODE INDICATOR
	output += "  Vertical Mode Indicator: "
	switch frame.vmi {
	case VMIUnknown:
		output += "Unknown mode or information unavailable\n"
	case VMIAcquiring:
		output += "“Acquiring” Mode\n"
	case VMICapturingOrMaintaining:
		output += "“Capturing” or “Maintaining” Mode\n"
	}

	// TARGET ALTITUDE TYPE
	output += "  Target Altitude Type: "
	switch frame.tat {
	case TATReferencedToFL:
		output += "Target altitude referenced to pressure-altitude (FL)\n"
	case TATReferencedToMSL:
		output += "Target altitude referenced to barometric corrected altitude (MSL)\n"
	}

	// TARGET ALTITUDE CAPABILITY
	output += "  Target Altitude Capability: "
	switch frame.tac {
	case TACReportingHoldingOnly:
		output += "Capability for reporting holding altitude only\n"
	case TACReportingHoldingAutopilotSelected:
		output += "Capability for reporting either holding altitude or autopilot control panel selected altitude\n"
	case TACReportingHoldingAutopilotSelectedFmsRnavLevelOff:
		output += "Capability for reporting either holding altitude, autopilot control panel selected altitude, or any FMS/RNAV level-off altitude\n"
	}

	// TARGET ALTITUDE
	output += fmt.Sprintf("  Target Altitude: %d ft\n", frame.ta)

	// HORIZONTAL DATA AVAILABLE/SOURCE INDICATOR
	output += "  Horizontal Data Available / Source Indicator: "
	switch frame.hdaSi {
	case HTSUnavailable:
		output += "No valid horizontal target state data is available\n"
	case HTSAutopilotSelectedValue:
		output += "Autopilot control panel selected value, such as Mode Control Panel (MCP) or Flight Control Unit (FCU)\n"
	case HTSMaintainingHeading:
		output += "Maintaining current heading or track angle (e.g. autopilot mode select)\n"
	case HTSFmsRnavSystem:
		output += "FMS/RNAV system (indicates track angle specified by leg type)\n"
	}

	// HORIZONTAL MODE INDICATOR
	output += "  Horizontal Mode Indicator: "
	switch frame.hmi {
	case HMIUnknown:
		output += "Unknown mode or information unavailable\n"
	case HMIAcquiring:
		output += "“Acquiring” mode\n"
	case HMICapturingOrMaintaining:
		output += "“Capturing” or “Maintaining” mode\n"
	}

	// TARGET HEADING/TRACK ANGLE
	output += fmt.Sprintf("  Target Heading/Track Angle: %d°\n", frame.thta)

	// TARGET HEADING/TRACK INDICATOR
	output += "  Target Heading/Track Indicator: "
	switch frame.thti {
	case THTITargetHeadingAngle:
		output += "Target Heading Angle is being reported\n"
	case THTITrackAngle:
		output += "Track Angle is being reported\n"
	}

	// NAVIGATION ACCURACY CATEGORY FOR POSITION (NACP)
	output += "  Navigation Accuracy Category for Position (NACp): "
	switch frame.nacp {
	case NACpUnknown:
		output += "EPU ≥ 18.52 km (10 NM) — Unknown accuracy\n"
	case NACp10NM:
		output += "EPU < 18.52 km (10 NM) — RNP –10 accuracy\n"
	case NACp4NM:
		output += "EPU < 7.408 km (4 NM) — RNP –4 accuracy\n"
	case NACp2NM:
		output += "EPU < 3.704 km (2 NM) — RNP –2 accuracy\n"
	case NACp1NM:
		output += "EPU < 1 852 m (1NM) — RNP –1 accuracy\n"
	case NACp0NM5:
		output += "EPU < 926 m (0.5 NM) — RNP –0.5 accuracy\n"
	case NACp0NM3:
		output += "EPU < 555.6 m ( 0.3 NM) — RNP –0.3 accuracy\n"
	case NACp0NM1:
		output += "EPU < 185.2 m (0.1 NM) — RNP –0.1 accuracy\n"
	case NACp0NM05:
		output += "EPU < 92.6 m (0.05 NM) — e.g. GPS (with SA)\n"
	case NACpGPS:
		output += "EPU < 30m and VEPU < 45m — e.g. GPS (SA off)\n"
	case NACpWAAS:
		output += "EPU < 10m and VEPU < 15m — e.g. WAAS\n"
	case NACpLAAS:
		output += "EPU < 3m and VEPU < 4m — e.g. LAAS\n"
	}

	// NAVIGATION INTEGRITY CATEGORY FOR BARO (NICBARO)
	output += "  Navigation Integrity Category for Baro (NICbaro): "
	switch frame.nicbaro {
	case NICBaroNotCrossChecked:
		output += "Barometric altitude reported has not been cross-checked against another source\n"
	case NICBaroCrossChecked:
		output += "Barometric altitude reported has been cross-checked against another source\n"
	}

	// SURVEILLANCE INTEGRITY LEVEL (SIL)
	output += "  Surveillance Integrity Level (SIL), probability of exceeding Horizontal Containment Radius (Rc): "
	switch frame.sil {
	case 0:
		output += "unknown\n"
	case 1:
		output += "≤1 × 10^–3 per flight hour or per sample\n"
	case 2:
		output += "≤1 × 10^–5 per flight hour or per sample\n"
	case 3:
		output += "≤1 × 10^–7 per flight hour or per sample\n"
	}
	output += "  Surveillance Integrity Level (SIL), probability of exceeding Vertical Integrity Containment Region (VPL): "
	switch frame.sil {
	case 0:
		output += "unknown\n"
	case 1:
		output += "≤1 × 10^–3 per flight hour or per sample\n"
	case 2:
		output += "≤1 × 10^–5 per flight hour or per sample\n"
	case 3:
		output += "≤2 × 10^–7 per 150 seconds or per sample\n"
	}

	// CAPABILITY/MODE CODES
	output += "  Capability/Mode Codes: "
	switch frame.cmc {
	case TcasAcasOperationalOrUnknown:
		output += "TCAS/ACAS operational or unknown; and "
	case TcasAcasNotOperational:
		output += "TCAS/ACAS not operational; and "
	}
	switch frame.cmcRa {
	case TcasAcasResolutionAdvisoryInactive:
		output += "No TCAS/ACAS Resolution Advisory active\n"
	case TcasAcasResolutionAdvisoryActive:
		output += "TCAS/ACAS Resolution Advisory active\n"
	}

	// EMERGENCY/PRIORITY STATUS
	output += "  Emergency / Priority Status: "
	switch frame.eps {
	case EmergencyNone:
		output += "No emergency\n"
	case EmergencyGeneral:
		output += "General emergency\n"
	case EmergencyLifeguardMedical:
		output += "Lifeguard/medical emergency\n"
	case EmergencyMinimumFuel:
		output += "Minimum fuel\n"
	case EmergencyNoCommunications:
		output += "No communications\n"
	case EmergencyUnlawfulInterference:
		output += "Unlawful interference\n"
	case EmergencyDownedAircraft:
		output += "Downed aircraft\n"
	}

	output += "\n"

	return output
}

type BDS62VerticalTargetState uint8

const VTSUnavailable = BDS62VerticalTargetState(0)            // No valid Vertical Target State data is available
const VTSAutopilotSelectedValue = BDS62VerticalTargetState(1) // Autopilot control panel selected value, such as Mode Control Panel (MCP) or Flight Control Unit (FCU)
const VTSHoldingAltitude = BDS62VerticalTargetState(2)        // Holding altitude
const VTSFmsRnavSystem = BDS62VerticalTargetState(3)          // FMS/RNAV system

type BDS62TargetAltitudeType uint8

const TATReferencedToFL = BDS62TargetAltitudeType(0)  // Target altitude referenced to pressure-altitude (FL)
const TATReferencedToMSL = BDS62TargetAltitudeType(1) // Target altitude referenced to barometric corrected altitude (MSL)

type BDS62TargetAltitudeCapability uint8

const TACReportingHoldingOnly = BDS62TargetAltitudeCapability(0)                             // Capability for reporting holding altitude only
const TACReportingHoldingAutopilotSelected = BDS62TargetAltitudeCapability(1)                // Capability for reporting either holding altitude or autopilot control panel selected altitude
const TACReportingHoldingAutopilotSelectedFmsRnavLevelOff = BDS62TargetAltitudeCapability(3) // Capability for reporting either holding altitude, autopilot control panel selected altitude, or any FMS/RNAV level-off altitude

type BDS62VerticalModeIndicator uint8

const VMIUnknown = BDS62VerticalModeIndicator(0)                // Unknown mode or information unavailable
const VMIAcquiring = BDS62VerticalModeIndicator(1)              // “Acquiring” Mode
const VMICapturingOrMaintaining = BDS62VerticalModeIndicator(2) // “Capturing” or “Maintaining” Mode

type BDS62HorizontalTargetState uint8

const HTSUnavailable = BDS62HorizontalTargetState(0)            // No valid horizontal target state data is available
const HTSAutopilotSelectedValue = BDS62HorizontalTargetState(1) // Autopilot control panel selected value, such as Mode Control Panel (MCP) or Flight Control Unit (FCU)
const HTSMaintainingHeading = BDS62HorizontalTargetState(2)     // Maintaining current heading or track angle (e.g. autopilot mode select)
const HTSFmsRnavSystem = BDS62HorizontalTargetState(3)          // FMS/RNAV system (indicates track angle specified by leg type)

type BDS62TargetHeadingTrackIndicator uint8

const THTITargetHeadingAngle = BDS62TargetHeadingTrackIndicator(0) // Target Heading Angle is being reported
const THTITrackAngle = BDS62TargetHeadingTrackIndicator(1)         // Track Angle is being reported

type BDS62HorizontalModeIndicator uint8

const HMIUnknown = BDS62HorizontalModeIndicator(0)                // Unknown mode or information unavailable
const HMIAcquiring = BDS62HorizontalModeIndicator(1)              // “Acquiring” mode
const HMICapturingOrMaintaining = BDS62HorizontalModeIndicator(2) // “Capturing” or “Maintaining” mode

type BDS62NACp uint8

const NACpUnknown = BDS62NACp(0) // EPU ≥ 18.52 km (10 NM) — Unknown accuracy
const NACp10NM = BDS62NACp(1)    // EPU < 18.52 km (10 NM) — RNP –10 accuracy
const NACp4NM = BDS62NACp(2)     // EPU < 7.408 km (4 NM) — RNP –4 accuracy
const NACp2NM = BDS62NACp(3)     // EPU < 3.704 km (2 NM) — RNP –2 accuracy
const NACp1NM = BDS62NACp(4)     // EPU < 1 852 m (1NM) — RNP –1 accuracy
const NACp0NM5 = BDS62NACp(5)    // EPU < 926 m (0.5 NM) — RNP –0.5 accuracy
const NACp0NM3 = BDS62NACp(6)    // EPU < 555.6 m ( 0.3 NM) — RNP –0.3 accuracy
const NACp0NM1 = BDS62NACp(7)    // EPU < 185.2 m (0.1 NM) — RNP –0.1 accuracy
const NACp0NM05 = BDS62NACp(8)   // EPU < 92.6 m (0.05 NM) — e.g. GPS (with SA)
const NACpGPS = BDS62NACp(9)     // EPU < 30m *and* VEPU < 45m — e.g. GPS (SA off)
const NACpWAAS = BDS62NACp(10)   // EPU < 10m *and* VEPU < 15m — e.g. WAAS
const NACpLAAS = BDS62NACp(11)   // EPU < 3m *and* VEPU < 4m — e.g. LAAS

type BDS62NICBaro uint8

const NICBaroNotCrossChecked = BDS62NICBaro(0) // The barometric altitude that is being reported in the Airborne Position Message is based on a Gilham coded input that has not been cross-checked against another source of pressure-altitude.
const NICBaroCrossChecked = BDS62NICBaro(1)    // The barometric altitude that is being reported in the Airborne Position Message is either based on a Gilham code input that has been cross-checked against another source of pressure-altitude and verified as being consistent or is based on a non-Gilham coded source.

type BDS62CapabilityModeCodes uint8

const TcasAcasOperationalOrUnknown = BDS62CapabilityModeCodes(0) // TCAS/ACAS operational or unknown
const TcasAcasNotOperational = BDS62CapabilityModeCodes(1)       // TCAS/ACAS not operational

type BDS62CapabilityModeCodesRA uint8

const TcasAcasResolutionAdvisoryInactive = BDS62CapabilityModeCodesRA(2) // No TCAS/ACAS Resolution Advisory active
const TcasAcasResolutionAdvisoryActive = BDS62CapabilityModeCodesRA(3)   // TCAS/ACAS Resolution Advisory active

func DecodeBDS62(mb []byte) (frame BDS62Frame, err error) {
	// Decode BDS 6,2 Frame: TARGET STATE AND STATUS INFORMATION
	// Type Code (TC): 29

	// check type code
	if (int(mb[0])&0b11111000)>>3 != 29 {
		err = errors.New("type code not 29")
		return
	}

	// check subtype code
	if (int(mb[0])&0b00000110)>>1 != 0 {
		err = errors.New("subtype code not 0")
		return
	}

	// check reserved bits
	if (int(mb[1])&0b00100000) != 0 || (int(mb[5])&0b00000011) != 0 || (int(mb[6])&0b11100000) != 0 {
		err = errors.New("reserved bits not 0")
		return
	}

	// VERTICAL DATA AVAILABLE/SOURCE INDICATOR
	switch ((int(mb[0]) & 0b00000001) << 1) + ((int(mb[0]) & 0b10000000) >> 7) {
	case 0:
		frame.vdaSi = VTSUnavailable
	case 1:
		frame.vdaSi = VTSAutopilotSelectedValue
	case 2:
		frame.vdaSi = VTSHoldingAltitude
	case 3:
		frame.vdaSi = VTSFmsRnavSystem
	}

	// TARGET ALTITUDE TYPE

	switch (int(mb[1]) & 0b01000000) >> 6 {
	case 0:
		frame.tat = TATReferencedToFL
	case 1:
		frame.tat = TATReferencedToMSL
	}

	// TARGET ALTITUDE CAPABILITY
	switch (int(mb[1]) & 0b00011000) >> 3 {
	case 0:
		frame.tac = TACReportingHoldingOnly
	case 1:
		frame.tac = TACReportingHoldingAutopilotSelected
	case 2:
		frame.tac = TACReportingHoldingAutopilotSelectedFmsRnavLevelOff
	case 3:
		err = errors.New("target altitude capability set to reserved value")
		return
	}

	// VERTICAL MODE INDICATOR
	switch (int(mb[1]) & 0b00000110) >> 1 {
	case 0:
		frame.vmi = VMIUnknown
	case 1:
		frame.vmi = VMIAcquiring
	case 2:
		frame.vmi = VMICapturingOrMaintaining
	case 3:
		err = errors.New("vertical mode indicator set to reserved value")
	}

	// TARGET ALTITUDE
	frame.ta = -1000 + (100 * (((int(mb[1]) & 0b00000001) << 9) + (int(mb[2]) << 1) + ((int(mb[3]) & 0b10000000) >> 7)))

	// HORIZONTAL DATA AVAILABLE/SOURCE INDICATOR
	switch (int(mb[3]) & 0b01100000) >> 5 {
	case 0:
		frame.hdaSi = HTSUnavailable
	case 1:
		frame.hdaSi = HTSAutopilotSelectedValue
	case 2:
		frame.hdaSi = HTSMaintainingHeading
	case 3:
		frame.hdaSi = HTSFmsRnavSystem
	}

	// TARGET HEADING/TRACK ANGLE
	frame.thta = ((int(mb[3]) & 0b00011111) << 4) + ((int(mb[4]) & 0b11110000) >> 4)
	if frame.thta >= 360 {
		err = errors.New("target heading / track angle invalid")
		return
	}

	// TARGET HEADING/TRACK INDICATOR
	switch (int(mb[4]) & 0b00001000) >> 3 {
	case 0:
		frame.thti = THTITargetHeadingAngle
	case 1:
		frame.thti = THTITrackAngle
	}

	// HORIZONTAL MODE INDICATOR
	switch (int(mb[4]) & 0b00000110) >> 1 {
	case 0:
		frame.hmi = HMIUnknown
	case 1:
		frame.hmi = HMIAcquiring
	case 2:
		frame.hmi = HMICapturingOrMaintaining
	case 3:
		err = errors.New("horizontal mode indicator set to reserved value")
		return
	}

	// NAVIGATION ACCURACY CATEGORY FOR POSITION (NACP)
	switch ((int(mb[4]) & 0b00000001) << 3) + ((int(mb[5]) & 0b11100000) >> 5) {
	case 0:
		frame.nacp = NACpUnknown
	case 1:
		frame.nacp = NACp10NM
	case 2:
		frame.nacp = NACp4NM
	case 3:
		frame.nacp = NACp2NM
	case 4:
		frame.nacp = NACp1NM
	case 5:
		frame.nacp = NACp0NM5
	case 6:
		frame.nacp = NACp0NM3
	case 7:
		frame.nacp = NACp0NM1
	case 8:
		frame.nacp = NACp0NM05
	case 9:
		frame.nacp = NACpGPS
	case 10:
		frame.nacp = NACpWAAS
	case 11:
		frame.nacp = NACpLAAS
	default:
		err = errors.New("NACp set to reserved value")
		return
	}

	// NAVIGATION INTEGRITY CATEGORY FOR BARO (NICBARO)
	switch (int(mb[5]) & 0b00010000) >> 4 {
	case 0:
		frame.nicbaro = NICBaroNotCrossChecked
	case 1:
		frame.nicbaro = NICBaroCrossChecked
	}

	// SURVEILLANCE INTEGRITY LEVEL (SIL)
	frame.sil = (int(mb[5]) & 0b00001100) >> 2

	// CAPABILITY/MODE CODES
	switch (int(mb[6]) & 0b00010000) >> 4 {
	case 0:
		frame.cmc = TcasAcasOperationalOrUnknown
	case 1:
		frame.cmc = TcasAcasNotOperational
	}
	switch (int(mb[6]) & 0b00001000) >> 3 {
	case 0:
		frame.cmcRa = TcasAcasResolutionAdvisoryInactive
	case 1:
		frame.cmcRa = TcasAcasResolutionAdvisoryActive
	}

	// EMERGENCY/PRIORITY STATUS
	frame.eps, err = decodeEmergencyState(int(mb[6]) & 0b00000111)

	return
}
