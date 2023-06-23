package df

type DF20message struct {
	// COMM-B ALTITUDE REPLY, DOWNLINK FORMAT 20
	fs int    // Flight status
	dr int    // Downlink request
	um int    // Utility message
	ac int    // Altitude code
	p  []byte // Parity

	Airborne bool
	Altitude float64
	ICAO     int    // Address announced: The address refers to the 24-bit transponder address (icao).
	MB       []byte // Message, Comm-B
}

func DecodeDF20(data []byte) (msg DF20message, err error) {
	// COMM-B ALTITUDE REPLY, DOWNLINK FORMAT 20
	// fs = Flight status (3.1.2.6.5.1)
	// dr = Downlink request (3.1.2.6.5.2)
	// um = Utility message (3.1.2.6.5.3)
	// ac = Altitude code (3.1.2.6.5.4)
	// mb = Message, Comm-B (3.1.2.6.6.1)
	// p = parity
	msg.fs = (int(data[0]) & 0b00000111)
	msg.dr = (int(data[1]) & 0b11111000) >> 3
	msg.um = ((int(data[1]) & 0b00000111) << 3) + ((int(data[2]) & 0b11100000) >> 5)
	msg.ac = ((int(data[2]) & 0b00011111) << 8) + (int(data[3]) & 0b11111111)
	msg.MB = []byte{data[4], data[5], data[6], data[7], data[8], data[9], data[10]}
	msg.p = []byte{data[11], data[12], data[13]}
	msg.ICAO = icaoFromCRC(data)
	msg.Airborne, err = airborneFromFlightStatus(msg.fs)
	if msg.ac != 0 {
		msg.Altitude, err = altitudeFromAltitudeCode13bit(msg.ac)
	}
	return
}
