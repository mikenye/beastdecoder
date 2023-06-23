package bds

import (
	"beastdecoder/df"
	"errors"

	"github.com/rs/zerolog/log"
)

const DF17 = df.DF17
const DF18 = df.DF18

func isBDS10(mb []byte) bool {
	// returns true if message is likely to be BDS 1,0
	// https://mode-s.org/decode/content/mode-s/9-inference.html

	// MB bits 1-8: equal to 0001 0000
	if int(mb[0]) != 0b00010000 {
		log.Debug().Str("reason", "bds code mismatch").Msg("not BDS10")
		return false
	}

	// MB bits 10-14: equal to all zeroes (reserved bits)
	if int(mb[1])&0b01111100 != 0 {
		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS10")
		return false
	}

	// Overlay capability conflict
	if (int(mb[1])&0b00000010)>>1 == 1 {
		if (int(mb[2])&0b11111110)>>1 < 5 {
			log.Debug().Str("reason", "overlay capability conflict").Msg("not BDS10")
			return false
		}
	}

	if (int(mb[1])&0b00000010)>>1 == 0 {
		if (int(mb[2])&0b11111110)>>1 > 4 {
			log.Debug().Str("reason", "overlay capability conflict").Msg("not BDS10")
			return false
		}
	}

	return true
}

func isBDS17(mb []byte) bool {
	// returns true if message is likely to be BDS 1,7
	// https://mode-s.org/decode/content/mode-s/9-inference.html

	// MB bit 7 should be equal to 1
	if int(mb[0])&0b10 != 0b10 {
		log.Debug().Str("reason", "bds code mismatch").Msg("not BDS17")
		return false
	}

	// MB bits 29-56 all zeroes (reserved bits)
	if int(mb[3])&0b00001111 != 0 {
		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS17")
		return false
	}
	if int(mb[4]) != 0 {
		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS17")
		return false
	}
	if int(mb[5]) != 0 {
		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS17")
		return false
	}
	if int(mb[6]) != 0 {
		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS17")
		return false
	}

	return true
}

func isBDS20(mb []byte) bool {
	// returns true if message is likely to be BDS 2,0
	// https://mode-s.org/decode/content/mode-s/9-inference.html

	// MB bits 1-8 equal to 0010 0000
	if int(mb[0]) != 0b00100000 {
		log.Debug().Str("reason", "bds code mismatch").Msg("not BDS20")
		return false
	}

	// Callsign	only contains 0–9, A–Z
	frame, err := DecodeBDS20(mb)
	if err != nil {
		log.Debug().AnErr("reason", err).Str("callsign", frame.Callsign).Msg("not BDS20")
		return false
	}

	return true
}

func isBDS30(mb []byte) bool {
	// returns true if message is likely to be BDS 3,0
	// https://mode-s.org/decode/content/mode-s/9-inference.html

	// MB bits 1-8 equal to 0011 0000
	if int(mb[0]) != 0b00110000 {
		log.Debug().Str("reason", "bds code mismatch").Msg("not BDS30")
		return false
	}

	// MB bits 29-39 not equal to 11
	if (int(mb[3])&0b00001100)>>2 == 0b11 {
		log.Debug().Str("reason", "threat type mismatch").Msg("not BDS30")
		return false
	}

	// MB bits 16-22 less than decimal 48
	if (((int(mb[1]) & 0b00000001) << 6) + ((int(mb[2]) & 0b11111100) >> 2)) >= 48 {
		log.Debug().Str("reason", "acas mismatch").Msg("not BDS30")
		return false
	}

	return true
}

func isBDS40(mb []byte) bool {
	// returns true if message is likely to be BDS 4,0
	// https://mode-s.org/decode/content/mode-s/9-inference.html

	// if MB bit 1 == 0, then bits 2-13 must be zero
	if int(mb[0])&0b10000000 == 0 {
		if int(mb[0])&0b01111111 != 0 {
			log.Debug().Str("reason", "MCP/FCU selected altitude not status consistent").Msg("not BDS40")
			return false
		}
		if int(mb[1])&0b11111000 != 0 {
			log.Debug().Str("reason", "MCP/FCU selected altitude not status consistent").Msg("not BDS40")
			return false
		}
	}

	// if MB bit 14 == 0, then bits 15-26 must be zero
	if int(mb[1])&0b00000100 == 0 {
		if int(mb[1])&0b00000011 != 0 {
			log.Debug().Str("reason", "FMS selected altitude not status consistent").Msg("not BDS40")
			return false
		}
		if int(mb[2])&0b11111111 != 0 {
			log.Debug().Str("reason", "FMS selected altitude not status consistent").Msg("not BDS40")
			return false
		}
		if int(mb[3])&0b11000000 != 0 {
			log.Debug().Str("reason", "FMS selected altitude not status consistent").Msg("not BDS40")
			return false
		}
	}

	// if MB bit 27 == 0, then bits 28-39 must be zero
	if int(mb[3])&0b00100000 == 0 {
		if int(mb[3])&0b00011111 != 0 {
			log.Debug().Str("reason", "Barometric pressure not status consistent").Msg("not BDS40")
			return false
		}
		if int(mb[4])&0b11111110 != 0 {
			log.Debug().Str("reason", "Barometric pressure not status consistent").Msg("not BDS40")
			return false
		}
	}

	// MB bits 40-47 should all be zero (reserved bits)
	if int(mb[4])&0b00000001 != 0 {
		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS40")
		return false
	}
	if int(mb[5])&0b11111110 != 0 {
		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS40")
		return false
	}

	// MB bits 52-53 should all be zero (reserved bits)
	if int(mb[6])&0b00011000 != 0 {
		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS40")
		return false
	}

	return true
}

func isBDS44(mb []byte) bool {
	// returns true if message is likely to be BDS 4,4
	// https://mode-s.org/decode/content/mode-s/9-inference.html

	// figure of merit must be less than 5
	fom := decodeBDS44figureOfMerit(mb)
	if fom >= 5 {
		log.Debug().Str("reason", "reserved bits in fom/source not zero").Msg("not BDS44")
		return false
	}

	// wind speed / direction bits must be status consistent, with speed less than 250 kt
	if int(mb[0])&0b00001000 == 0 {
		if int(mb[0])&0b00000111 != 0 {
			log.Debug().Str("reason", "wind speed not status consistent").Msg("not BDS44")
			return false
		}
		if int(mb[1]) != 0 {
			log.Debug().Str("reason", "wind speed not status consistent").Msg("not BDS44")
			return false
		}
		if int(mb[2])&0b11111110 != 0 {
			log.Debug().Str("reason", "wind speed not status consistent").Msg("not BDS44")
			return false
		}
	} else {
		ws, _, err := decodeBDS44windSpeedDirection(mb)
		if err != nil {
			log.Debug().AnErr("err", err).Str("reason", "wind speed").Msg("not BDS44")
			return false
		}
		if ws > 250 {
			log.Debug().Str("reason", "wind speed out of range").Msg("not BDS44")
			return false
		}
	}

	// static air temperature must be between -80 and 60 deg C
	sat, err := decodeBDS44staticAirTemperature(mb)
	if err != nil {
		return false
	}
	if sat < -80 || sat > 60 {
		return false
	}

	// average static pressure bits must be status consistent
	if int(mb[4])&0b00100000 == 0 {
		if int(mb[4])&0b00011111 != 0 {
			log.Debug().Str("reason", "static pressure not status consistent").Msg("not BDS44")
			return false
		}
		if int(mb[5])&0b11111100 != 0 {
			log.Debug().Str("reason", "static pressure not status consistent").Msg("not BDS44")
			return false
		}
	}

	// turbulence bits must be status consistent
	if int(mb[5])&0b00000010 == 0 {
		if int(mb[5])&0b00000001 != 0 {
			log.Debug().Str("reason", "turbulence not status consistent").Msg("not BDS44")
			return false
		}
		if int(mb[6])&0b10000000 != 0 {
			log.Debug().Str("reason", "turbulence not status consistent").Msg("not BDS44")
			return false
		}
	}

	// humidity bits must be status consistent
	if int(mb[6])&0b01000000 == 0 {
		if int(mb[6])&0b00111111 != 0 {
			log.Debug().Str("reason", "humidity not status consistent").Msg("not BDS44")
			return false
		}
	}

	return true
}

func isBDS45(mb []byte) bool {

	// turbulence bits must be status consistent
	if int(mb[0])&0b10000000 == 0 {
		if int(mb[0])&0b01100000 != 0 {
			log.Debug().Str("reason", "turbulence not status consistent").Msg("not BDS45")
			return false
		}
	}

	// wind shear bits must be status consistent
	if int(mb[0])&0b00010000 == 0 {
		if int(mb[0])&0b00001100 != 0 {
			log.Debug().Str("reason", "wind shear not status consistent").Msg("not BDS45")
			return false
		}
	}

	// microburst bits must be status consistent
	if int(mb[0])&0b00000010 == 0 {
		if int(mb[0])&0b00000001 != 0 {
			log.Debug().Str("reason", "microburst not status consistent").Msg("not BDS45")
			return false
		}
		if int(mb[1])&0b10000000 != 0 {
			log.Debug().Str("reason", "microburst not status consistent").Msg("not BDS45")
			return false
		}
	}

	// icing bits must be status consistent
	if int(mb[1])&0b01000000 == 0 {
		if int(mb[1])&0b00110000 != 0 {
			log.Debug().Str("reason", "icing not status consistent").Msg("not BDS45")
			return false
		}
	}

	// wake vortex bits must be status consistent
	if int(mb[1])&0b00001000 == 0 {
		if int(mb[1])&0b00000110 != 0 {
			log.Debug().Str("reason", "wake vortex not status consistent").Msg("not BDS45")
			return false
		}
	}

	// static air temp bits must be status consistent
	if int(mb[1])&0b00000001 == 0 {
		if int(mb[2])&0b01111111 != 0 {
			log.Debug().Str("reason", "static air temp not status consistent").Msg("not BDS45")
			return false
		}
		if int(mb[3])&0b11000000 != 0 {
			log.Debug().Str("reason", "static air temp not status consistent").Msg("not BDS45")
			return false
		}
	}

	// average static pressure bits must be status consistent
	if int(mb[3])&0b00100000 == 0 {
		if int(mb[3])&0b00011111 != 0 {
			log.Debug().Str("reason", "average static pressure not status consistent").Msg("not BDS45")
			return false
		}
		if int(mb[4])&0b11111100 != 0 {
			log.Debug().Str("reason", "average static pressure not status consistent").Msg("not BDS45")
			return false
		}
	}

	// radio height bits must be status consistent
	if int(mb[4])&0b00000010 == 0 {
		if int(mb[4])&0b00000001 != 0 {
			log.Debug().Str("reason", "radio height not status consistent").Msg("not BDS45")
			return false
		}
		if int(mb[5]) != 0 {
			log.Debug().Str("reason", "radio height not status consistent").Msg("not BDS45")
			return false
		}
		if int(mb[6])&0b11100000 != 0 {
			log.Debug().Str("reason", "radio height not status consistent").Msg("not BDS45")
			return false
		}
	}

	// reserved bits must be zero
	if int(mb[6])&0b00011111 == 0 {
		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS45")
		return false
	}

	// TODO: if temp > 60 or temp < -80: return false

	return true
}

func isBDS50(mb []byte) bool {
	// returns true if message is likely to be BDS 5,0
	// https://mode-s.org/decode/content/mode-s/9-inference.html

	// roll angle must have status consistent bits, and be between -50 and 50 degrees
	if int(mb[0])&0b10000000 == 0 {
		if int(mb[0])&0b00111111 != 0 {
			log.Debug().Str("reason", "roll angle bits not status consistent").Msg("not BDS50")
			return false
		}
		if int(mb[1])&0b01100000 != 0 {
			log.Debug().Str("reason", "roll angle bits not status consistent").Msg("not BDS50")
			return false
		}
	} else {
		rollAngle, err := decodeBDS50roll(mb)
		if err != nil {
			log.Debug().AnErr("err", err).Msg("not BDS50")
			return false
		}
		if rollAngle < -50 || rollAngle > 50 {
			log.Debug().Str("reason", "roll angle out of range").Float64("rollAngle", rollAngle).Msg("not BDS50")
			return false
		}
	}

	// if MB bit 12 == 0, then bits 13-23 must be 0
	if int(mb[1])&0b00010000 == 0 {
		if int(mb[1])&0b00001111 != 0 {
			log.Debug().Str("reason", "true track angle bits not status consistent").Msg("not BDS50")
			return false
		}
		if int(mb[2])&0b11111110 != 0 {
			log.Debug().Str("reason", "true track angle bits not status consistent").Msg("not BDS50")
			return false
		}
	}

	// check ground speed is valid and between 0 and 600 knots
	if int(mb[2])&0b00000001 == 0 {
		if int(mb[3]) != 0 {
			log.Debug().Str("reason", "ground speed bits not status consistent").Msg("not BDS50")
			return false
		}
		if int(mb[4])&0b11000000 != 0 {
			log.Debug().Str("reason", "ground speed bits not status consistent").Msg("not BDS50")
			return false
		}
	} else {
		gs, err := decodeBDS50groundSpeed(mb)
		if err != nil {
			log.Debug().AnErr("err", err).Msg("not BDS50")
			return false
		}
		if gs < 0 || gs > 500 {
			log.Debug().Str("reason", "ground speed out of range").Float64("gs", gs).Msg("not BDS50")
			return false
		}
	}

	// if MB bit 35 == 0, then bits 36-45 must be 0
	if int(mb[4])&0b00100000 == 0 {
		if int(mb[4])&0b00011111 != 0 {
			log.Debug().Str("reason", "track angle rate bits not status consistent").Msg("not BDS50")
			return false
		}
		if int(mb[5])&0b11111000 != 0 {
			log.Debug().Str("reason", "track angle rate bits not status consistent").Msg("not BDS50")
			return false
		}
	}

	// check true airspeed status consistent & is valid and between 0 and 500 knots
	if (int(mb[5])&0b00000100)>>2 == 0 {
		if int(mb[5])&0b00000011 != 0 {
			log.Debug().Str("reason", "true airspeed bits not status consistent").Msg("not BDS50")
			return false
		}
		if int(mb[6]) != 0 {
			log.Debug().Str("reason", "true airspeed bits not status consistent").Msg("not BDS50")
			return false
		}
	} else {
		tas, err := decodeBDS50trueAirspeed(mb)
		if err != nil {
			log.Debug().AnErr("err", err).Msg("not BDS50")
			return false
		}
		if tas < 0 || tas > 500 {
			log.Debug().Str("reason", "true airspeed bits out of range").Float64("tas", tas).Msg("not BDS50")
			return false
		}
	}

	// fmt.Println(reason)
	return true

}

func isBDS53(mb []byte) bool {
	// returns true if the message is likely to be BDS 5,3

	// check magnetic heading bits are status consistent
	if (int(mb[0]) & 0b10000000) == 0 {
		if (int(mb[0]) & 0b00111111) != 0 {
			log.Debug().Str("reason", "magnetic heading bits not status consistent").Msg("not BDS53")
			return false
		}
		if (int(mb[1]) & 0b11110000) != 0 {
			log.Debug().Str("reason", "magnetic heading bits not status consistent").Msg("not BDS53")
			return false
		}
	}

	// check indicated airspeed bits are status consistent
	if (int(mb[1]) & 0b00001000) == 0 {
		if (int(mb[1]) & 0b00000111) != 0 {
			log.Debug().Str("reason", "indicated airspeed bits not status consistent").Msg("not BDS53")
			return false
		}
		if (int(mb[2]) & 0b11111110) != 0 {
			log.Debug().Str("reason", "indicated airspeed bits not status consistent").Msg("not BDS53")
			return false
		}
	}

	// check mach number bits are status consistent
	if (int(mb[2]) & 0b00000001) == 0 {
		if int(mb[3]) != 0 {
			log.Debug().Str("reason", "mach number bits not status consistent").Msg("not BDS53")
			return false
		}
		if (int(mb[4]) & 0b10000000) != 0 {
			log.Debug().Str("reason", "mach number bits not status consistent").Msg("not BDS53")
			return false
		}
	}

	// check true airspeed bits are status consistent
	if (int(mb[4]) & 0b01000000) == 0 {
		if (int(mb[4]) & 0b00111111) != 0 {
			log.Debug().Str("reason", "true airspeed bits not status consistent").Msg("not BDS53")
			return false
		}
		if (int(mb[5]) & 0b11111100) != 0 {
			log.Debug().Str("reason", "true airspeed bits not status consistent").Msg("not BDS53")
			return false
		}
	}

	// check altitude rate bits are status consistent
	if (int(mb[5]) & 0b00000010) == 0 {
		if int(mb[6]) != 0 {
			log.Debug().Str("reason", "altitude rate bits not status consistent").Msg("not BDS53")
			return false
		}
	}

	return true
}

func isBDS545556(mb []byte) bool {
	// returns true if the message is likely to be BDS 5,4, 5,5, 5,6
	_, err := DecodeBDS54(mb)
	if err != nil {
		log.Debug().AnErr("reason", err).Msg("not BDS54")
		return false
	}
	return true
}

func isBDS60(mb []byte) bool {
	// returns true if message is likely to be BDS 6,0
	// https://mode-s.org/decode/content/mode-s/9-inference.html

	// if MB bit 1 == 0, then bits 2-12 must be 0
	if int(mb[0])&0b10000000 == 0 {
		if int(mb[0])&0b01111111 != 0 {
			log.Debug().Str("reason", "Magnetic heading bits not status consistent").Msg("not BDS60")
			return false
		}
		if int(mb[1])&0b11110000 != 0 {
			log.Debug().Str("reason", "Magnetic heading bits not status consistent").Msg("not BDS60")
			return false
		}
	}

	// indicated airspeed must have status consistent bits, and be between 0 and 500 knots
	if int(mb[1])&0b00001000 == 0 {
		if int(mb[1])&0b00000111 != 0 {
			log.Debug().Str("reason", "Indicated airspeed bits not status consistent").Msg("not BDS60")
			return false
		}
		if int(mb[2])&0b11111110 != 0 {
			log.Debug().Str("reason", "Indicated airspeed bits not status consistent").Msg("not BDS60")
			return false
		}
	} else {
		ias, err := decodeBDS60indicatedAirspeed(mb)
		if err != nil {
			log.Debug().AnErr("reason", err).Msg("not BDS60")
			return false
		}
		if ias < 0 || ias > 500 {
			log.Debug().Float64("ias", ias).Str("reason", "Indicated airspeed out of range").Msg("not BDS60")
			return false
		}
	}

	// mach number must have status consistent bits, and be between 0 and 1
	if int(mb[2])&0b00000001 == 0 {
		if int(mb[3]) != 0 {
			log.Debug().Str("reason", "Mach number bits not status consistent").Msg("not BDS60")
			return false
		}
		if int(mb[4])&0b11000000 != 0 {
			log.Debug().Str("reason", "Mach number bits not status consistent").Msg("not BDS60")
			return false
		}
	} else {
		mach, err := decodeBDS60machNumber(mb)
		if err != nil {
			log.Debug().AnErr("err", err).Msg("not BDS60")
			return false
		}
		if mach < 0 || mach > 1 {
			log.Debug().Float64("mach", mach).Str("reason", "Mach number out of range").Msg("not BDS60")
			return false
		}
	}

	// Barometric vertical rate must have status consistent bits, and be between -6000 and 6000 fpm
	if int(mb[4])&0b00100000 == 0 {
		if int(mb[4])&0b00011111 != 0 {
			log.Debug().Str("reason", "Barometric vertical rate bits not status consistent").Msg("not BDS60")
			return false
		}
		if int(mb[5])&0b11111000 != 0 {
			log.Debug().Str("reason", "Barometric vertical rate bits not status consistent").Msg("not BDS60")
			return false
		}
	} else {
		bar, err := decodeBDS60barometricAltitudeRate(mb)
		if err != nil {
			log.Debug().AnErr("reason", err).Msg("not BDS60")
			return false
		}
		if bar < -6000 || bar > 6000 {
			log.Debug().Float64("bar", bar).Str("reason", "Barometric vertical rate out of range").Msg("not BDS60")
			return false
		}
	}

	// Inertial vertical velocity must have status consistent bits, and be between -6000 and 6000 fpm
	if int(mb[5])&0b00000100 == 0 {
		if int(mb[5])&0b00000011 != 0 {
			log.Debug().Str("reason", "Inertial vertical rate bits not status consistent").Msg("not BDS60")
			return false
		}
		if int(mb[6]) != 0 {
			log.Debug().Str("reason", "Inertial vertical rate bits not status consistent").Msg("not BDS60")
			return false
		}
	} else {
		ivv, err := decodeBDS60GNSSAltitudeRate(mb)
		if err != nil {
			log.Debug().AnErr("reason", err).Msg("not BDS60")
			return false
		}
		if ivv < -6000 || ivv > 6000 {
			log.Debug().Str("reason", "Inertial vertical (GNSS altitude) rate out of range").Msg("not BDS60")
			return false
		}
	}
	return true
}

// func isBDS61(mb []byte) bool {

// 	// check type code
// 	if (int(mb[0])&0b11111000)>>3 != 28 {
// 		log.Debug().Str("reason", "type code not 28").Msg("not BDS61")
// 		return false
// 	}

// 	// check subtype code
// 	if (int(mb[0]) & 0b00000111) != 1 {
// 		log.Debug().Str("reason", "subtype code not 1").Msg("not BDS61")
// 		return false
// 	}

// 	// check emergency state
// 	es := (int(mb[1]) & 0b11100000) >> 5
// 	if es >= 6 {
// 		log.Debug().Str("reason", "emergency state bits set to reserved value").Msg("not BDS61")
// 		return false
// 	}

// 	// check reserved bits
// 	if (int(mb[1])&0b00011111) != 0 || int(mb[2]) != 0 || int(mb[3]) != 0 || int(mb[4]) != 0 || int(mb[5]) != 0 || int(mb[6]) != 0 {
// 		log.Debug().Str("reason", "reserved bits not zero").Msg("not BDS61")
// 		return false
// 	}

// 	return true
// }

// func isBDS62(mb []byte) bool {
// 	_, err := DecodeBDS62(mb)
// 	if err != nil {
// 		log.Debug().AnErr("reason", err).Msg("not BDS62")
// 		return false
// 	}
// 	return true
// }

func InferBDS(df df.DownlinkFormat, mb []byte) (possibleBDScodes []BDScode, err error) {
	// BDS codes identification
	// https://mode-s.org/decode/content/mode-s/9-inference.html

	// For ADS-B / Mode-S extended squitter
	if df == DF17 || df == DF18 {

		// Check BDS05
		_, e := DecodeBDS05(mb)
		if e != nil {
			log.Debug().AnErr("reason", e).Msg("not BDS05")
		} else {
			possibleBDScodes = append(possibleBDScodes, BDS05)
			return
		}

		// Check BDS06
		_, e = DecodeBDS06(mb)
		if e != nil {
			log.Debug().AnErr("reason", e).Msg("not BDS06")
		} else {
			possibleBDScodes = append(possibleBDScodes, BDS06)
			return
		}

		// Check BDS08
		_, e = DecodeBDS08(mb)
		if e != nil {
			log.Debug().AnErr("reason", e).Msg("not BDS08")
		} else {
			possibleBDScodes = append(possibleBDScodes, BDS08)
			return
		}

		// Check BDS09
		_, e = DecodeBDS09(mb)
		if e != nil {
			log.Debug().AnErr("reason", e).Msg("not BDS09")
		} else {
			possibleBDScodes = append(possibleBDScodes, BDS09)
			return
		}

		// Check BDS61
		_, e = DecodeBDS61(mb)
		if e != nil {
			log.Debug().AnErr("reason", e).Msg("not BDS61")
		} else {
			possibleBDScodes = append(possibleBDScodes, BDS61)
			return
		}

		// Check BDS62
		_, e = DecodeBDS62(mb)
		if e != nil {
			log.Debug().AnErr("reason", e).Msg("not BDS62")
		} else {
			possibleBDScodes = append(possibleBDScodes, BDS62)
			return
		}

		// Check BDS65
		_, e = DecodeBDS65(mb)
		if e != nil {
			log.Debug().AnErr("reason", e).Msg("not BDS65")
		} else {
			possibleBDScodes = append(possibleBDScodes, BDS65)
		}

		// Check BDS07
		_, e = DecodeBDS07(mb)
		if e != nil {
			log.Debug().AnErr("reason", e).Msg("not BDS07")
		} else {
			possibleBDScodes = append(possibleBDScodes, BDS07)
		}

		// Sanity checks
		if len(possibleBDScodes) == 0 {
			if log.Debug().Enabled() {
				err = errors.New("could not infer extended squitter bds for DF17")
			}
		}
		if len(possibleBDScodes) > 1 {
			err = errors.New("multiple bds match")
		}

		return
	}

	if isBDS10(mb) {
		possibleBDScodes = append(possibleBDScodes, BDS10)
	}

	if isBDS17(mb) {
		possibleBDScodes = append(possibleBDScodes, BDS17)
	}

	if isBDS20(mb) {
		possibleBDScodes = append(possibleBDScodes, BDS20)
	}

	if isBDS30(mb) {
		possibleBDScodes = append(possibleBDScodes, BDS30)
	}

	if isBDS40(mb) {
		possibleBDScodes = append(possibleBDScodes, BDS40)
	}

	if isBDS50(mb) {
		possibleBDScodes = append(possibleBDScodes, BDS50)
	}

	if isBDS60(mb) {
		possibleBDScodes = append(possibleBDScodes, BDS60)
	}

	// TODO: BDS44/45

	if len(possibleBDScodes) == 0 {
		if log.Debug().Enabled() {
			err = errors.New("could not infer bds")
		}
	}
	if len(possibleBDScodes) > 1 {
		err = errors.New("multiple bds match")
	}

	return possibleBDScodes, err
}
