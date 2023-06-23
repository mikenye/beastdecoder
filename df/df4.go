package df

// DF4: Surveillance, Altitude Reply

type DF4message struct {
	// DF4: Surveillance, Altitude Reply

	fs int    // Flight status (FS): Shows status of alert, special position pulse (SPI, in Mode A only) and aircraft status (airborne or on-ground).
	dr int    // Downlink request (DR): Contains the type of request. In surveillance replies, only values 0, 1, 4, and 5 are used.
	um int    // Utility message (UM): 6 bits, contains transponder communication status information
	ac int    // Altitude Code (AC): Encodes the altitude of the aircraft.
	ap []byte // Address parity bytes

	Airborne bool    // airborne status
	Altitude float64 // decoded Altitude
	ICAO     int     // ICAO aircraft address
}

func DecodeDF4(data []byte) (msg DF4message, err error) {
	// Surveillance altitude reply - https://mode-s.org/decode/content/mode-s/3-surveillance.html
	// fs = Flight status (FS): 3 bits, shows status of alert, special position pulse (SPI, in Mode A only) and aircraft status (airborne or on-ground). The field is interpreted as:
	//        000: no alert, no SPI, aircraft is airborne
	//        001: no alert, no SPI, aircraft is on-ground
	//        010: alert, no SPI, aircraft is airborne
	//        011: alert, no SPI, aircraft is on-ground
	//        100: alert, SPI, aircraft is airborne or on-ground
	//        101: no alert, SPI, aircraft is airborne or on-ground
	//        110: reserved
	//        111: not assigned
	// dr = Downlink request (DR): 5 bits, contains the type of request. In surveillance replies, only values 0, 1, 4, and 5 are used. The field can be decoded as:
	//        00000: no downlink request
	//        00001: request to send Comm-B message
	//        00100: Comm-B broadcast message 1 available
	//        00101: Comm-B broadcast message 2 available
	// um = Utility message (UM): 6 bits, contains transponder communication status information.
	//        IIS: The first 4 bits of UM indicate the interrogator identifier code.
	//        IDS: The last 2 bits indicate the type of reservation made by the interrogator.
	//              00: no information
	//              01: IIS contains Comm-B interrogator identifier code
	//              10: IIS contains Comm-C interrogator identifier code
	//              11: IIS contains Comm-D interrogator identifier code
	// ac = Altitude code
	// ap = Address parity
	msg.fs = (int(data[0]) & 0b00000111)
	msg.dr = (int(data[1]) & 0b11111000) >> 3
	msg.um = ((int(data[1]) & 0b00000111) << 3) + ((int(data[2]) & 0b11100000) >> 5)
	msg.ac = ((int(data[2]) & 0b00011111) << 8) + (int(data[3]) & 0b11111111)
	msg.ap = []byte{data[4], data[5], data[6]}
	msg.ICAO = icaoFromCRC(data)

	// set airborne based on Flight Status bits
	msg.Airborne, err = airborneFromFlightStatus(msg.fs)

	// set decoded altitude
	msg.Altitude, err = altitudeFromAltitudeCode13bit(msg.ac)

	return
}
