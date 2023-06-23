package bds

import (
	"errors"
	"regexp"
)

type BDScode int

const (
	// ADS-B
	BDS05 = BDScode(5)  // BDS 0,5 - Extended squitter airborne position
	BDS06 = BDScode(6)  // BDS 0,6 - Extended squitter surface position
	BDS08 = BDScode(8)  // BDS 0,8 - Extended squitter aircraft identification and category
	BDS09 = BDScode(9)  // BDS 0,9 - Extended squitter airborne velocity
	BDS61 = BDScode(61) // BDS 6,1 - Extended squitter emergency/priority status
	BDS62 = BDScode(62) // BDS 6,2 - Extended squitter Target State and Status
	BDS65 = BDScode(65) // BDS 6,5 - Extended squitter aircraft operational status

	// ELS
	BDS10 = BDScode(10) // BDS 1,0 - Data link capability report
	BDS17 = BDScode(17) // BDS 1,7 - Common usage GICB capability report
	BDS20 = BDScode(20) // BDS 2,0 - Aircraft identification
	BDS30 = BDScode(30) // BDS 3,0 - ACAS active resolution advisory

	// EHS
	BDS40 = BDScode(40) // BDS 4,0 - Selected vertical intention
	BDS50 = BDScode(50) // BDS 5,0 - Track and turn report
	BDS60 = BDScode(60) // BDS 6,0 - Heading and speed report

	// MRAR & MHR
	BDS44 = BDScode(44) // BDS 4,4 - Meteorological routine air report
	BDS45 = BDScode(45) // BDS 4,5 - Meteorological hazard report

	// Currently Unknown
	BDS07 = BDScode(7) // BDS 0,7 - Extended squitter status
)

// Emergency/priority status
type EmergencyPriorityStatus uint8

const EmergencyNone = EmergencyPriorityStatus(0)                 // No emergency
const EmergencyGeneral = EmergencyPriorityStatus(1)              // General emergency
const EmergencyLifeguardMedical = EmergencyPriorityStatus(2)     // Lifeguard/medical emergency
const EmergencyMinimumFuel = EmergencyPriorityStatus(3)          // Minimum fuel
const EmergencyNoCommunications = EmergencyPriorityStatus(4)     // No communications
const EmergencyUnlawfulInterference = EmergencyPriorityStatus(5) // Unlawful interference
const EmergencyDownedAircraft = EmergencyPriorityStatus(6)       // Downed aircraft

// Time Synchronization
type TimeSynchronization uint8

const TimeNotSynchronizedToUTC = TimeSynchronization(0)
const TimeSynchronizedToUTC = TimeSynchronization(1)

// --------------------

// character map for callsigns
var callsignCharMap = []string{
	"#",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"#", "#", "#", "#", "#",
	" ",
	"#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#",
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"#", "#", "#", "#", "#", "#"}

// regex to check for valid callsign
var validCallsign = regexp.MustCompile("^[A-Z0-9]*$")

// character map for waypoints
var waypointCharMap = callsignCharMap

// regex to check for valid waypoints
var validWaypoint = validCallsign

// --------------------

func decodeEmergencyState(emergencyStateBits int) (emergencyState EmergencyPriorityStatus, err error) {
	// EMERGENCY/PRIORITY STATUS
	switch emergencyStateBits {
	case 0:
		emergencyState = EmergencyNone
	case 1:
		emergencyState = EmergencyGeneral
	case 2:
		emergencyState = EmergencyLifeguardMedical
	case 3:
		emergencyState = EmergencyMinimumFuel
	case 4:
		emergencyState = EmergencyNoCommunications
	case 5:
		emergencyState = EmergencyUnlawfulInterference
	case 6:
		emergencyState = EmergencyDownedAircraft
	case 7:
		err = errors.New("emergency/priority status set to reserved value")
	}
	return
}
