package df

// DF0: Short Air-Air Surveillance

type DF0message struct {
	// DF0: Short Air-Air Surveillance

	vs int    // Vertical status (VS): Indicates whether the aircraft is airborne (0) or on the ground (1).
	cc int    // Cross-link capability (CC): Refers to the capability of reply DF=16 upon request of UF=0. When set to 1, the cross-link is supported. Otherwise, the field is set to 0.
	sl int    // Sensitivity level (SL): Represents the sensitivity level of the ACAS system, except that 0 indicates the ACAS is inoperative.
	ri int    // Reply information (RI): Indicates the type of reply to interrogating aircraft. For ACAS message, valid values are 0 and from 2 to 4. Other values are not part of the ACAS.
	ac int    // Altitude Code (AC): Encodes the altitude of the aircraft.
	ap []byte // Address parity bytes

	Airborne bool    // airborne status
	Altitude float64 // decoded Altitude
	ICAO     int     // ICAO aircraft address
}

func DecodeDF0(data []byte) (msg DF0message, err error) {
	// ACAS surveillance reply (downlink) - https://mode-s.org/decode/content/mode-s/4-acas.html
	// vs = Vertical status (VS): 1 bit, it indicates whether the aircraft is airborne (0) or on the ground (1).
	// cc = Cross-link capability (CC): 1 bit, it refers to the capability of reply DF=16 upon request of UF=0. When this 1-bit field is set to 1, the cross-link is supported. Otherwise, the field is set to 0.
	// sl = Sensitivity level (SL): 3 bits, it represents the sensitivity level of the ACAS system, except that 0 indicates the ACAS is inoperative.
	// ri = Reply information (RI): 4 bits, it indicates the type of reply to interrogating aircraft. For ACAS message, valid values are 0 and from 2 to 4. Other values are not part of the ACAS:
	//        0000: No operating ACAS
	//        0010: ACAS with resolution capability inhibited
	//        0011: ACAS with vertical-only resolution capability
	//        0111: ACAS with vertical and horizontal resolution capability
	// ac = Altitude Code (AC): 13 bits, it encodes the altitude of the aircraft.
	// ap = Address parity
	msg.vs = (int(data[0]) & 0b00000100) >> 2

	// Set msg.Airborne
	switch msg.vs {
	case 0:
		msg.Airborne = true
	case 1:
		msg.Airborne = false
	}

	msg.cc = (int(data[0]) & 0b00000010) >> 1
	// RESERVED = int(data[0]) & 0b0000001
	msg.sl = (int(data[1]) & 0b11100000) >> 5
	// RESERVED = (int(data[1]) & 0b00011000) >> 3
	msg.ri = ((int(data[1]) & 0b00000111) << 1) + ((int(data[2]) & 0b10000000) >> 7)
	// RESERVED = (int(data[2]) & 0b01100000) >> 5
	msg.ac = ((int(data[2]) & 0b00011111) << 8) + (int(data[3]) & 0b11111111)
	msg.Altitude, err = altitudeFromAltitudeCode13bit(msg.ac)
	if err != nil {
		return
	}
	msg.ap = data[len(data)-3:]
	msg.ICAO = icaoFromCRC(data)
	return
}
