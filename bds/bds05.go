package bds

// Airborne position
// Type Code (TC): 9-18 (Baro)
// Type Code (TC): 20-22 (GNSS)
// https://mode-s.org/decode/content/ads-b/3-airborne-position.html

import (
	"beastdecoder/common"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type BDS05Frame struct {
	Tc int // Type Code

	// Time Synchronization (T)
	//
	// This field shall indicate whether or not the time of applicability of the message is synchronized with UTC time.
	//
	//  TimeNotSynchronizedToUTC = time is not synchronized to UTC.
	//  TimeSynchronizedToUTC = time of applicability is synchronized to UTC time.
	//
	// Synchronization shall only be used for airborne position messages having the top two horizontal
	// position precision categories (format TYPE Codes 9, 10, 20 and 21).
	//
	// When T = TimeSynchronizedToUTC, the time of validity in the airborne position message format shall be encoded in the 1-bit F field which,
	// in addition to CPR format type, indicates the 0.2-second time tick for UTC time of position validity.
	// The F bit shall alternate between 0 and 1 for successive 0.2-second time ticks, beginning with F = 0 when the time of
	// applicability is an exact even-numbered UTC second.
	T TimeSynchronization

	// Compact Position Reporting (CPR) Format (F)
	//
	// In order to achieve coding that is unambiguous worldwide,
	// CPR shall use two format types, known as even and odd.
	// This field shall be used to define the CPR format type.
	//  CprFormatEvenFrame = even format coding
	//  CprFormatOddFrame  = odd format coding
	F common.CprFormat

	// Latitude
	//
	// The CPR-encoded latitude value.
	LatCpr int

	// Longitude
	//
	// The CPR-encoded longitude value
	LonCpr int

	// Altitude (encoded)
	//
	//  Type Code (TC) of 9-18: Barometric altitude encoded in 25- or 100-foot increments.
	//  Type Code (TC) of 20-22: GNSS height above ellipsoid (HAE).
	ac int

	// Altitude (ft)
	Altitude float64

	// Single Antenna Flag (SAF)
	//
	// This field shall indicate the type of antenna system that is being used to transmit extended squitters.
	//  SingleTransmitAntenna = single transmit antenna.
	//  DualTransmitAntenna = dual transmit antenna system.
	//
	// At any time that the diversity configuration cannot guarantee that both antenna channels are functional,
	// then the single antenna subfield shall be set to SingleTransmitAntenna.
	saf SingleAntennaFlag

	// Surveillance Status
	//
	// The surveillance status field in the airborne position message format shall encode information from the
	// aircraft’s Mode A code and SPI condition indication.
	ss SurveillanceStatus
}

// Surveillance status
type SurveillanceStatus uint8

const SurveillanceStatusNoCondition = SurveillanceStatus(0)    // No condition information
const SurveillanceStatusPermanentAlert = SurveillanceStatus(1) // Permanent alert (emergency condition)
const SurveillanceStatusTemporaryAlert = SurveillanceStatus(2) // Temporary alert (change in Mode A identity code other than emergency condition)
const SurveillanceStatusSPICondition = SurveillanceStatus(3)   // SPI condition

// Single Antenna Flag (SAF)
type SingleAntennaFlag uint8

const SingleTransmitAntenna = SingleAntennaFlag(1) // single transmit antenna
const DualTransmitAntenna = SingleAntennaFlag(0)   // dual transmit antenna system

func DecodeBDS05(mb []byte) (frame BDS05Frame, err error) {
	// decode airborne position bytes into struct
	// https://mode-s.org/decode/content/ads-b/3-airborne-position.html

	frame = BDS05Frame{}

	// Type Code
	frame.Tc = (int(mb[0]) & 0b11111000) >> 3
	if (frame.Tc < 9 || frame.Tc > 18) && (frame.Tc < 20 || frame.Tc > 22) {
		err = errors.New("type code not 9-18 or 20-22")
		return
	}

	// COMPACT POSITION REPORTING (CPR) FORMAT (F)
	switch (int(mb[2]) & 0b00000100) >> 2 {
	case 0:
		frame.F = common.CprFormatEvenFrame
	case 1:
		frame.F = common.CprFormatOddFrame
	}

	// TIME SYNCHRONIZATION (T)
	switch (int(mb[2]) & 0b00001000) >> 3 {
	case 0:
		frame.T = TimeNotSynchronizedToUTC
	case 1:
		frame.T = TimeSynchronizedToUTC
	}

	// SURVEILLANCE STATUS
	frame.ss = SurveillanceStatus((int(mb[0]) & 0b00000110) >> 1)

	switch int(mb[0]) & 0b00000001 {
	case 0:
		frame.saf = DualTransmitAntenna
	case 1:
		frame.saf = SingleTransmitAntenna
	}

	// ALTITUDE
	frame.ac = ((int(mb[1]) << 4) + ((int(mb[2]) & 0b11110000) >> 4)) // Encoded altitude

	frame.LatCpr = (((int(mb[2]) & 0b00000011) << 15) + (int(mb[3]) << 7) + ((int(mb[4]) & 0b11111110) >> 1)) // Encoded latitude
	frame.LonCpr = (((int(mb[4]) & 0b00000001) << 16) + (int(mb[5]) << 8) + int(mb[6]))                       // Encoded longitude

	if log.Debug().Enabled() {
		log.Debug().Int("tc", frame.Tc).Str("ac", fmt.Sprintf("%012b", frame.ac)).Msg("ac & tc for alt tshooting")
	}

	frame.Altitude, err = DecodeBDS05Altitude(frame.ac, frame.Tc)
	if err != nil {
		return
	}

	return
}

func DecodeBDS05Altitude(ac, tc int) (altFt float64, err error) {

	log.Debug().Msg("start")
	log.Debug().Str("ac", fmt.Sprintf("%013b", ac)).Msg("ac bits BEFORE")

	if ac == 0b0000000000000 {
		err = errors.New("altitude information is not available or has been determined invalid")
		return
	}

	log.Debug().Str("ac", fmt.Sprintf("%013b", ac)).Msg("ac bits AFTER")
	log.Debug().Int("tc", tc).Msg("type code")
	defer log.Debug().Msg("finish")

	switch {

	case tc < 19:

		acTopBits := ac & 0b111111000000
		acBtmBits := ac & 0b000000111111
		ac = (acTopBits << 1) + acBtmBits

		// determine M and Q bits
		mBit := (ac & 0b0000001000000) >> 6
		qBit := (ac & 0b0000000010000) >> 4

		switch {

		// If the M bit (bit 26) and the Q bit (bit 28) equal 0,
		// the altitude shall be coded according to the pattern for Mode C replies.
		// Starting with bit 20 the sequence shall be C1, A1, C2, A2, C4, A4, ZERO, B1, ZERO, B2, D2, B4, D4.
		case mBit == 0 && qBit == 0:
			// bitfield format:
			// +----+----+----+----+----+----+---+----+---+----+----+----+----+
			// | C1 | A1 | C2 | A2 | C4 | A4 | 0 | B1 | 0 | B2 | D2 | B4 | D4 |
			// +----+----+----+----+----+----+---+----+---+----+----+----+----+

			// graycode format required:
			// D2 D4 A1 A2 A4 B1 B2 B4     C1 C2 C4
			d2 := (ac & 0b0000000000100) >> 5
			d4 := (ac & 0b0000000000001) << 6
			a1 := (ac & 0b0100000000000) >> 6
			a2 := (ac & 0b0001000000000) >> 5
			a4 := (ac & 0b0000010000000) >> 4
			b1 := (ac & 0b0000000100000) >> 3
			b2 := (ac & 0b0000000001000) >> 2
			b4 := (ac & 0b0000000000010) >> 1

			c1 := (ac & 0b1000000000000) >> 10
			c2 := (ac & 0b0010000000000) >> 9
			c4 := (ac & 0b0000100000000) >> 8

			// see: https://github.com/junzis/pyModeS/blob/8e2051af688defd26537ca52928b0a15ba734a49/pyModeS/py_common.py#L359
			graycodeDAB := d2 + d4 + a1 + a2 + a4 + b1 + b2 + b4
			graycodeC := c1 + c2 + c4

			n500 := graycodeDAB
			n500 ^= n500 >> 8
			n500 ^= n500 >> 4
			n500 ^= n500 >> 2
			n500 ^= n500 >> 1

			n100 := graycodeC
			n100 ^= n100 >> 8
			n100 ^= n100 >> 4
			n100 ^= n100 >> 2
			n100 ^= n100 >> 1

			switch {
			case n100 == 0 || n100 == 5 || n100 == 6:
				err = errors.New("altitude information invalid")
				return
			case n100 == 7:
				n100 = 5
			case common.Modulo(float64(n500), 2.0) > 0:
				n100 = 6 - n100
			}

			altFt = (float64(n500)*500.0 + float64(n100)*100.0) - 1300.0

			if log.Debug().Enabled() {
				log.Debug().Int("n100", n100).Int("n500", n500).Float64("altFt", altFt).Int("M", mBit).Int("Q", qBit).Int("A1", a1).Int("A2", a2).Int("A4", a4).Int("B1", b1).Int("B2", b2).Int("B4", b4).Int("C1", c1).Int("C2", c2).Int("C4", c4).Int("D2", d2).Int("D4", d4).Msg("altitude calculation mBit == 0 && qBit == 0")
			}

		// If the M bit equals 0 and the Q bit equals 1,
		// the 11-bit field represented by bits 20 to 25, 27 and 29 to 32
		// shall represent a binary coded field with a least significant bit (LSB) of 25 ft.
		// The binary value of the positive decimal integer “N” shall be encoded to report pressure-altitude in the range [(25 N – 1 000) plus or minus 12.5 ft].
		// The coding of (mBit == 0 && qBit == 0) shall be used to report pressure-altitude above 50 187.5 ft.
		case mBit == 0 && qBit == 1:
			altFt = 25.0*float64(((ac&0b1111110000000)>>2)+((ac&0b0000000100000)>>1)+(ac&0b0000000001111)) - 1000.0

			if log.Debug().Enabled() {
				log.Debug().Float64("altFt", altFt).Int("M", mBit).Int("Q", qBit).Msg("altitude calculation mBit == 0 && qBit == 1")
			}

		// If the M bit equals 1, the 12-bit field represented by bits 20 to 25 and 27 to 31 shall be reserved for encoding altitude in metric units.
		case mBit == 1:
			altFt = float64(((ac&0b1111110000000)>>2)+((ac&0b0000000111110)>>1)) * 3.28084

			if log.Debug().Enabled() {
				log.Debug().Float64("altFt", altFt).Int("M", mBit).Int("Q", qBit).Msg("altitude calculation mBit == 1")
			}
		}
		break

	default:
		err = errors.New("GNSS altitude returned, not accurate enough for position reporting")
	}

	return
}

// func DecodeBDS05Altitude(ac int, tc int) (altFt float64, err error) {
// 	// Returns altitude (if available) from altitude code
// 	// https://mode-s.org/decode/content/mode-s/3-surveillance.html#sec:alt_code
// 	// 3.1.2.6.5.4 AC: Altitude code.

// 	if ac == 0b000000000000 {
// 		err = errors.New("altitude information is not available or has been determined invalid")
// 		return
// 	}

// 	switch {

// 	case tc < 19:

// 		// determine M and Q bits
// 		qBit := (ac & 0b000000010000) >> 4

// 		switch {

// 		// When Q=0, the altitude is encoded with a 100 feet increment
// 		case qBit == 0:
// 			altFt = 100.0*float64(((ac&0b111111100000)>>1)+(ac&0b000000001111)) - 1000.0

// 		// When Q=1, the altitude is encoded with a 25 feet increment.
// 		case qBit == 1:
// 			altFt = 25.0*float64(((ac&0b111111100000)>>1)+(ac&0b000000001111)) - 1000.0

// 		}

// 	default:
// 		altFt = float64(ac) * 3.28084
// 	}

// 	return

// }
