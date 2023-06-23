package df

type DF21message struct {
	// COMM-B IDENTITY REPLY, DOWNLINK FORMAT 21
	fs int    // Flight status
	dr int    // Downlink request
	um int    // Utility message
	id int    // Identity code
	p  []byte // Parity

	Airborne bool
	ICAO     int // Address announced: The address refers to the 24-bit transponder address (icao).
	Squawk   int
	MB       []byte // Message, Comm-B
}

func DecodeDF21(data []byte) (msg DF21message, err error) {
	// Comm-B, identity reply - https://mode-s.org/decode/content/mode-s/5-commb.html
	// COMM-B IDENTITY REPLY, DOWNLINK FORMAT 21
	// fs = Flight status (3.1.2.6.5.1)
	// dr = Downlink request (3.1.2.6.5.2)
	// um = Utility message (3.1.2.6.5.3)
	// id = Identity code (3.1.2.6.7.1)
	// mb = Message, Comm-B (3.1.2.6.6.1)
	// p = parity
	msg.fs = (int(data[0]) & 0b00000111)
	msg.dr = (int(data[1]) & 0b11111000) >> 3
	msg.um = ((int(data[1]) & 0b00000111) << 3) + ((int(data[2]) & 0b11100000) >> 5)
	msg.id = ((int(data[2]) & 0b00011111) << 8) + (int(data[3]) & 0b11111111)
	msg.MB = []byte{data[4], data[5], data[6], data[7], data[8], data[9], data[10]}
	msg.p = []byte{data[11], data[12], data[13]}
	msg.ICAO = icaoFromCRC(data)

	msg.Airborne, err = airborneFromFlightStatus(msg.fs)
	if err != nil {
		return
	}

	msg.Squawk, err = squawkFromIdentityCode(msg.id)
	if err != nil {
		return
	}

	return
}
