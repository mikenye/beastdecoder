package df

// DF16: Long Air-Air ACAS

type DF16message struct {
	// DF16: Long Air-Air ACAS

	vs int    // Vertical status (VS): 1 bit, it indicates whether the aircraft is airborne (0) or on the ground (1).
	sl int    // Sensitivity level (SL): 3 bits, it represents the sensitivity level of the ACAS system, except that 0 indicates the ACAS is inoperative.
	ri int    // Reply information (RI): 4 bits, it indicates the type of reply to interrogating aircraft. For ACAS message, valid values are 0 and from 2 to 4.
	ac int    // Altitude Code (AC): 13 bits, it encodes the altitude of the aircraft.
	mv []byte // Message, V
	ap []byte // Address parity

	Airborne bool
	Altitude float64
	ICAO     int // Address announced: The address refers to the 24-bit transponder address (icao).

}

func DecodeDF16(data []byte) (msg DF16message, err error) {
	// ACAS coordination reply - https://mode-s.org/decode/content/mode-s/4-acas.html
	// vs = Vertical status (VS): 1 bit, it indicates whether the aircraft is airborne (0) or on the ground (1).
	// sl = Sensitivity level (SL): 3 bits, it represents the sensitivity level of the ACAS system, except that 0 indicates the ACAS is inoperative.
	// ri = Reply information (RI): 4 bits, it indicates the type of reply to interrogating aircraft. For ACAS message, valid values are 0 and from 2 to 4. Other values are not part of the ACAS:
	//        0000: No operating ACAS
	//        0010: ACAS with resolution capability inhibited
	//        0011: ACAS with vertical-only resolution capability
	//        0111: ACAS with vertical and horizontal resolution capability
	// ac = Altitude Code (AC): 13 bits, it encodes the altitude of the aircraft.
	// mv = Message, V
	// ap = Address parity
	msg.vs = (int(data[0]) & 0b00000100) >> 2
	// RESERVED = (int(data[0]) & 0b00000011)
	msg.sl = (int(data[1]) & 0b11100000) >> 5
	// RESERVED = (int(data[1]) & 0b00011000) >> 3
	msg.ri = ((int(data[1]) & 0b00000111) << 1) + ((int(data[2]) & 0b10000000) >> 7)
	// RESERVED = (int(data[2]) & 0b01100000) >> 5
	msg.ac = ((int(data[2]) & 0b00011111) << 8) + (int(data[3]))
	msg.mv = []byte{data[4], data[5], data[6], data[7], data[8], data[9], data[10]}
	msg.ap = []byte{data[11], data[12], data[13]}
	msg.ICAO = icaoFromCRC(data)

	// Set msg.Airborne
	switch msg.vs {
	case 0:
		msg.Airborne = true
	case 1:
		msg.Airborne = false
	}

	// Set altitude
	msg.Altitude, err = altitudeFromAltitudeCode13bit(msg.ac)
	if err != nil {
		return
	}

	return
}
