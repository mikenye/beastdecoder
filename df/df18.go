package df

type DF18message struct {
	// DF18: ADS-B Extended Squitter not sent from a Mode S transponder

	cf int // Control field. This 3-bit (6-8) downlink field in DF = 18 shall be used to define the format of the 112-bit transmission as follows.
	// Code 0 = ADS-B ES/NT devices that report the ICAO 24-bit address in the AA field (3.1.2.8.7)
	// Code 1 = Reserved for ADS-B for ES/NT devices that use other addressing techniques in the AA field (3.1.2.8.7.3)
	// Code 2 = Fine format TIS-B message
	// Code 3 = Coarse format TIS-B message
	// Code 4 = Reserved for TIS-B management messages
	// Code 5 = TIS-B messages that relay ADS-B messages that use other addressing techniques in the AA field
	// Code 6 = ADS-B rebroadcast using the same type codes and message formats as defined for DF = 17 ADS-B messages
	// Code 7 = Reserved

	pi []byte // Parity/Interrogator ID

	Tc   int    // message type code
	ICAO int    // Address announced: The address refers to the 24-bit transponder address (icao).
	ME   []byte // Message, extended squitter
}

func DecodeDF18(data []byte) (msg DF18message) {
	// ca = Transponder capability
	// icao = ICAO aircraft address
	// tc = message type code
	// me = Message, extended squitter
	// pi = Parity/Interrogator ID
	msg.cf = (int(data[0]) & 0b00000111)                                   // transponder capability
	msg.ICAO = (int(data[3])) + (int(data[2]) << 8) + (int(data[1]) << 16) // icao aircraft address
	msg.ME = data[4:11]                                                    // message, extended squitter
	msg.Tc = ((int(data[4]) & 0b11111000) >> 3)                            // type code
	msg.pi = data[11:]                                                     // parity/interrogator id
	return msg
}
