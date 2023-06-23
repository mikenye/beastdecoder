package bds

import "errors"

type BDS65Frame struct {
	cc4 int // En-Route Capabilities (CC-4)
	cc3 int // Terminal Area Operational Capabilities (CC-3)
	cc2 int // Approach/Landing Operational Capabilities (CC-2)
	cc1 int // Surface Operational Capabilities (CC-1)

	om4 int // En-Route Operational Capability Status (OM-4)
	om3 int // Terminal Area Operational Capability Status (OM-3)
	om2 int // Approach/Landing Operational Capability Status (OM-2)
	om1 int // Surface Operational Capability Status (OM-1)
}

func DecodeBDS65(mb []byte) (frame BDS65Frame, err error) {

	// check format code
	if (int(mb[0])&0b11111000)>>3 != 31 {
		err = errors.New("type code not 31")
		return
	}

	// check subtype code
	if (int(mb[0]) & 0b00000111) != 0 {
		err = errors.New("subtype code not 0")
		return
	}

	// check reserved bits
	if int(mb[5]) != 0 || int(mb[6]) != 0 {
		err = errors.New("reserved bits not zero")
		return
	}

	// En-Route Capabilities (CC-4)
	frame.cc4 = (int(mb[1]) & 0b11110000) >> 4

	// Terminal Area Operational Capabilities (CC-3)
	frame.cc3 = (int(mb[1]) & 0b00001111)

	// Approach/Landing Operational Capabilities (CC-2)
	frame.cc2 = (int(mb[2]) & 0b11110000) >> 4

	// Surface Operational Capabilities (CC-1)
	frame.cc1 = (int(mb[2]) & 0b00001111)

	// En-Route Operational Capability Status (OM-4)
	frame.om4 = (int(mb[3]) & 0b11110000) >> 4

	// Terminal Area Operational Capability Status (OM-3)
	frame.om3 = (int(mb[3]) & 0b00001111)

	// Approach/Landing Operational Capability Status (OM-2)
	frame.om2 = (int(mb[4]) & 0b11110000) >> 4

	// Surface Operational Capability Status (OM-1)
	frame.om1 = (int(mb[4]) & 0b00001111)

	return
}
