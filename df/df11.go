package df

// DF11: All Call Reply

type DF11message struct {
	// DF11: All Call Reply

	// Capability: The definition of transponder capability is the same as in ADS-B messages.
	ca int

	// Parity/interrogator identifier: The decoding of PI is similar to the decoding of ADS-B parity.
	pi []byte

	// Address announced: The address refers to the 24-bit transponder address (icao).
	ICAO int
}

func DecodeDF11(data []byte) (msg DF11message) {
	// All-call reply - https://mode-s.org/decode/content/mode-s/2-allcall.html
	// ca = Capability: The definition of transponder capability is the same as in ADS-B messages.
	// aa = Address announced: The address refers to the 24-bit transponder address (icao).
	// pi = Parity/interrogator identifier: The decoding of PI is similar to the decoding of ADS-B parity.
	msg.ca = (int(data[0]) & 0b00000111)
	msg.ICAO = (int(data[1]) << 16) + (int(data[2]) << 8) + int(data[3]) // aa
	msg.pi = []byte{data[4], data[5], data[6]}
	return msg
}
