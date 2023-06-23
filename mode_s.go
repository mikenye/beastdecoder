package main

// // Operation Status Sub-Type
// type operationStatusSubType uint8

// const operationStatusAirborne = operationStatusSubType(0)
// const operationStatusSurface = operationStatusSubType(1)

// // --------------------------------------

// // --------------------------------------

// func calcAltitudeDifferenceGNSSBaro(sDif, dAlt int) (delta int, err error) {
// 	// https://mode-s.org/decode/content/ads-b/5-airborne-velocity.html#gnss-and-barometric-altitudes-difference
// 	switch dAlt {
// 	case 0:
// 		err = errors.New("information not available")
// 	default:
// 		switch sDif {
// 		case 0:
// 			delta = (1 * ((dAlt) - 1) * 25)
// 		case 1:
// 			delta = (-1 * ((dAlt) - 1) * 25)
// 		}
// 	}
// 	return delta, err
// }

// func calcVerticalRate(svr verticalRateSign, vr int) (vs int) {
// 	// https://mode-s.org/decode/content/ads-b/5-airborne-velocity.html#vertical-rate

// 	switch svr {
// 	case verticalRateClimb:
// 		vs = 64 * (vr - 1)
// 	case verticalRateDescent:
// 		vs = -64 * (vr - 1)
// 	}
// 	return vs
// }

// func capabilityDefinition(ca int) (def string, err error) {
// 	// Return transponder capability definition from capability bits
// 	// https://mode-s.org/decode/content/ads-b/1-basics.html#capability
// 	switch {
// 	case ca == 0:
// 		def = "Level 1 transponder"
// 		break
// 	case ca >= 1 && ca <= 3:
// 		def = "Reserved"
// 		err = errors.New("reserved")
// 		break
// 	case ca == 4:
// 		def = "Level 2+ transponder, with ability to set CA to 7, on-ground"
// 	case ca == 5:
// 		def = "Level 2+ transponder, with ability to set CA to 7, airborne"
// 	case ca == 6:
// 		def = "Level 2+ transponder, with ability to set CA to 7, either on-ground or airborne"
// 	case ca == 7:
// 		def = "Signifies the Downlink Request value is 0, or the Flight Status is 2, 3, 4, or 5, either airborne or on the ground"
// 	}
// 	return def, err
// }

// func decodeFlightStatus(fs int) (alert, spi, airborne bool, err error) {
// 	// decode Flight status (FS) - https://mode-s.org/decode/content/mode-s/3-surveillance.html#message-structure
// 	// shows status of alert, special position pulse (SPI, in Mode A only) and aircraft status (airborne or on-ground).
// 	switch fs {
// 	case 0b000:
// 		alert = false
// 		spi = false
// 		airborne = true
// 	case 0b001:
// 		alert = false
// 		spi = false
// 		airborne = false
// 	case 0b010:
// 		alert = true
// 		spi = false
// 		airborne = true
// 	case 0b011:
// 		alert = true
// 		spi = false
// 		airborne = false
// 	case 0b100:
// 		alert = true
// 		spi = true
// 		err = errors.New("aircraft is airborne or on-ground")
// 	case 0b101:
// 		alert = false
// 		spi = true
// 		err = errors.New("aircraft is airborne or on-ground")
// 	case 0b110:
// 		err = errors.New("reserved")
// 	case 0b111:
// 		err = errors.New("not assigned")
// 	}
// 	return alert, spi, airborne, err
// }

// func getDF(data []byte) df.DownlinkFormat {
// 	// Returns DF (downlink format) from message
// 	return df.DownlinkFormat(((int(data[0]) & 0b11111000) >> 3) + 100)
// }

// func isSurfacePosValid(refLat, refLon, surfaceLat, surfaceLon float64) bool {

// 	surfacePos := haversine.Coord{Lat: surfaceLat, Lon: surfaceLon}
// 	receiverPos := haversine.Coord{Lat: refLat, Lon: refLon}

// 	_, km := haversine.Distance(surfacePos, receiverPos)
// 	nm := km * 0.539957

// 	if nm < 45 {
// 		return true
// 	}
// 	return false
// }
