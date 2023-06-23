package bds

import "errors"

// BDS code 6,1 — Aircraft status

type BDS61Frame struct {
	st int // Sub-type code

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

func DecodeBDS61(mb []byte) (frame BDS61Frame, err error) {
	// BDS code 6,1 — Aircraft status
	// Type Code (TC): 28

	// check type code
	if (int(mb[0])&0b11111000)>>3 != 28 {
		err = errors.New("type code not 28")
		return
	}

	// SUBTYPE CODE
	frame.st = (int(mb[0]) & 0b00000111)

	switch frame.st {

	// Subtype code 1, Emergency/priority status
	case 1:
		frame.eps, err = decodeEmergencyState((int(mb[1]) & 0b11100000) >> 5)
		if err != nil {
			return
		}
	}

	return
}
