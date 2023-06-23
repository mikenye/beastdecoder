package df

// DF17: ADS-B Extended Squitter sent from a Mode S transponder - https://mode-s.org/decode/content/ads-b/1-basics.html#message-structure

type DF17message struct {
	// DF17: ADS-B Extended Squitter sent from a Mode S transponder

	ca int    // Vertical status (VS): 1 bit, it indicates whether the aircraft is airborne (0) or on the ground (1).
	Tc int    // message type code
	pi []byte // Parity/Interrogator ID

	ICAO int    // Address announced: The address refers to the 24-bit transponder address (icao).
	ME   []byte // Message, extended squitter
}

func DecodeDF17(data []byte) (msg DF17message) {
	// ca = Transponder capability
	// icao = ICAO aircraft address
	// tc = message type code
	// me = Message, extended squitter
	// pi = Parity/Interrogator ID
	msg.ca = (int(data[0]) & 0b00000111)                                   // transponder capability
	msg.ICAO = (int(data[3])) + (int(data[2]) << 8) + (int(data[1]) << 16) // icao aircraft address
	msg.ME = data[4:11]                                                    // message, extended squitter
	msg.Tc = ((int(data[4]) & 0b11111000) >> 3)                            // type code
	msg.pi = data[11:]                                                     // parity/interrogator id
	return msg
}
