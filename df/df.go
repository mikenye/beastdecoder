package df

import (
	"beastdecoder/common"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

// DF: Downlink format.
// This downlink format field (5 bits long except in format 24 where it is 2 bits long)
// shall serve as the downlink format descriptor in all Mode S replies.

type DownlinkFormat uint8

const DF0 = DownlinkFormat(0)
const DF4 = DownlinkFormat(4)
const DF5 = DownlinkFormat(5)
const DF11 = DownlinkFormat(11)
const DF16 = DownlinkFormat(16)
const DF17 = DownlinkFormat(17)
const DF18 = DownlinkFormat(18)
const DF19 = DownlinkFormat(19)
const DF20 = DownlinkFormat(20)
const DF21 = DownlinkFormat(21)
const DF24 = DownlinkFormat(24)

func airborneFromFlightStatus(fs int) (airborne bool, err error) {
	switch fs {
	case 0b000:
		airborne = true
	case 0b001:
		airborne = false
	case 0b010:
		airborne = true
	case 0b011:
		airborne = false
	case 0b100:
		airborne = false
	case 0b101:
		airborne = false
	case 0b110:
		err = errors.New("flight status set to reserved")
	case 0b111:
		err = errors.New("flight status set to not assigned")
	}
	return
}

func altitudeFromAltitudeCode13bit(ac int) (altFt float64, err error) {
	// Returns altitude (if available) from 13-bit altitude code (does not work with 12-bit altitude codes!!)
	// https://mode-s.org/decode/content/mode-s/3-surveillance.html#sec:alt_code
	// 3.1.2.6.5.4 AC: Altitude code.

	if ac == 0b0000000000000 {
		err = errors.New("altitude information is not available or has been determined invalid")
		return
	}

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
			log.Debug().Int("n100", n100).Int("n500", n500).Float64("altFt", altFt).Int("M", mBit).Int("Q", qBit).Int("A1", a1).Int("A2", a2).Int("A4", a4).Int("B1", b1).Int("B2", b2).Int("B4", b4).Int("C1", c1).Int("C2", c2).Int("C4", c4).Int("D2", d2).Int("D4", d4).Msg("altitude calculation")
		}

	// If the M bit equals 0 and the Q bit equals 1,
	// the 11-bit field represented by bits 20 to 25, 27 and 29 to 32
	// shall represent a binary coded field with a least significant bit (LSB) of 25 ft.
	// The binary value of the positive decimal integer “N” shall be encoded to report pressure-altitude in the range [(25 N – 1 000) plus or minus 12.5 ft].
	// The coding of (mBit == 0 && qBit == 0) shall be used to report pressure-altitude above 50 187.5 ft.
	case mBit == 0 && qBit == 1:
		altFt = 25.0*float64(((ac&0b1111110000000)>>2)+((ac&0b0000000100000)>>1)+(ac&0b0000000001111)) - 1000.0

		if log.Debug().Enabled() {
			log.Debug().Float64("altFt", altFt).Int("M", mBit).Int("Q", qBit).Msg("altitude calculation")
		}

	// If the M bit equals 1, the 12-bit field represented by bits 20 to 25 and 27 to 31 shall be reserved for encoding altitude in metric units.
	case mBit == 1:
		altFt = float64(((ac&0b1111110000000)>>2)+((ac&0b0000000111110)>>1)) * 3.28084

		if log.Debug().Enabled() {
			log.Debug().Float64("altFt", altFt).Int("M", mBit).Int("Q", qBit).Msg("altitude calculation")
		}
	}
	return
}

func calcFrameCRC(inframe []byte) (remainder []byte) {
	// calculate ADS-B frame parity, returns remainder bytes
	// https://mode-s.org/decode/content/ads-b/8-error-control.html#ads-b-parity

	startingByte := len(inframe) - 3

	frame := make([]byte, len(inframe))
	copy(frame, inframe)

	// zero out PI bytes
	frame[startingByte] = 0x00
	frame[startingByte+1] = 0x00
	frame[startingByte+2] = 0x00

	// the CRC generator
	G := []int{0b11111111, 0b11111010, 0b00000100, 0b10000000}

	// perform destructive crc
	for ibyte := 0; ibyte <= len(frame)-4; ibyte++ { //-4?
		for ibit := 0; ibit <= 7; ibit++ {
			mask := 0x80 >> ibit
			bits := int(frame[ibyte]) & mask

			if bits > 0 {
				frame[ibyte] = byte(int(frame[ibyte]) ^ (G[0] >> ibit))
				frame[ibyte+1] = byte(int(frame[ibyte+1]) ^ (0xFF & ((G[0] << (8 - ibit)) | (G[1] >> ibit))))
				frame[ibyte+2] = byte(int(frame[ibyte+2]) ^ (0xFF & ((G[1] << (8 - ibit)) | (G[2] >> ibit))))
				frame[ibyte+3] = byte(int(frame[ibyte+3]) ^ (0xFF & ((G[2] << (8 - ibit)) | (G[3] >> ibit))))
			}
		}
	}

	// return results
	return []byte{frame[startingByte], frame[startingByte+1], frame[startingByte+2]}
}

func GetDF(data []byte) DownlinkFormat {
	// Returns DF (downlink format) from message
	return DownlinkFormat((int(data[0]) & 0b11111000) >> 3)
}

func icaoFromCRC(data []byte) (icao int) {
	// Decodes ICAO address from message with address parity
	// https://mode-s.org/decode/content/mode-s/1-basics.html#icao-address-recovery

	// prep parity bytes
	p := []byte{data[len(data)-3], data[len(data)-2], data[len(data)-1]}

	// combine bits
	pint := (int(p[0]) << 16) + (int(p[1]) << 8) + int(p[2])

	// calculate frame CRC bytes
	pcalc := calcFrameCRC(data)

	// combine bits
	pcalcint := (int(pcalc[0]) << 16) + (int(pcalc[1]) << 8) + int(pcalc[2])

	// XOR magic
	icao = pint ^ pcalcint

	// icao.
	return icao
}

func squawkFromIdentityCode(id int) (squawk int, err error) {
	// returns a squawk code from an Identity code (DF5)
	// https://mode-s.org/decode/content/mode-s/3-surveillance.html#sec:id_code

	// The 13-bit identity code encodes the 4 octal digit squawk code (from 0000 to 7777). The structure of this field is shown as follows:

	// +----+----+----+----+----+----+---+----+----+----+----+----+----+
	// | C1 | A1 | C2 | A2 | C4 | A4 | X | B1 | D1 | B2 | D2 | B4 | D4 |
	// +----+----+----+----+----+----+---+----+----+----+----+----+----+
	// The binary representation of the octal digit is:

	// A4 A2 A1 | B4 B2 B1 | C4 C2 C1 | D4 D2 D1

	a := (((id & 0b0000010000000) >> 5) + ((id & 0b0001000000000) >> 8) + ((id & 0b0100000000000) >> 11)) * 1000
	b := (((id & 0b0000000000010) << 1) + ((id & 0b0000000001000) >> 2) + ((id & 0b0000000100000) >> 5)) * 100
	c := (((id & 0b0000100000000) >> 6) + ((id & 0b0010000000000) >> 9) + ((id & 0b1000000000000) >> 12)) * 10
	d := (((id & 0b0000000000001) << 2) + ((id & 0b0000000000100) >> 1) + ((id & 0b0000000010000) >> 4))

	squawk = a + b + c + d
	if squawk < 0 || squawk >= 10000 {
		err = errors.New(fmt.Sprintf("invalid squawk code: %d", squawk))
	}

	return
}
