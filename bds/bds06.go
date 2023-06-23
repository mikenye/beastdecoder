package bds

// Surface position
// Type Code (TC): 5 to 8
// https://mode-s.org/decode/content/ads-b/4-surface-position.html

import (
	"beastdecoder/common"
	"errors"
	"fmt"
	"math"
)

type BDS06Frame struct {
	tc int // Type Code

	// Movement
	//
	// This field shall provide information on the ground speed of the aircraft.
	// A non-linear scale shall be used as defined in the following table where speeds are given in km/h and kt.
	//
	//        0 = No information available
	//        1 = Aircraft stopped (ground speed <0.2315 km/h (0.125 kt))
	//      2-8 = 0.2315 km/h (0.125 kt) ≤ ground speed < 1.852 km/h (1 kt); LSB: 0.2315 km/h (0.125 kt)
	//     9-12 = 1.852 km/h (1 kt) ≤ ground speed < 3.704 km/h (2 kt); LSB: 0.463 km/h (0.25 kt)
	//    13-38 = 3.704 km/h (2 kt) ≤ ground speed < 27.78 km/h (15 kt); LSB: 0.926 km/h (0.5 kt)
	//    39-93 = 27.78 km/h (15 kt) ≤ ground speed < 129.64 km/h (70 kt); LSB: 1.852 km/h (1.0 kt)
	//   94-108 = 129.64 km/h (70 kt) ≤ ground speed < 185.2 km/h (100 kt); LSB: 3.704 km/h (2.0 kt)
	//  109-123 = 185.2 km/h (100 kt) ≤ ground speed < 324.1 km/h (175 kt); LSB: 9.26 km/h (5.0 kt)
	//      124 = Ground speed ≥ 324.1 km/h (175 kt)
	mov int

	// Ground Speed (decoded)
	GroundSpeed string

	// Ground track status
	//
	// This field shall define the validity of the ground track value. Coding for this field shall be as follows:
	//  GroundTrackStatusInvalid = invalid
	//  GroundTrackStatusValid = valid
	S GroundTrackStatus

	// Ground track value
	// This field shall define the direction (in degrees clockwise from true north) of aircraft motion on the surface.
	// The ground track shall be encoded as an unsigned angular weighted binary numeral,
	// with an MSB of 180 degrees and an LSB of 360/128 degrees, with zero indicating true north.
	// The data in the field shall be rounded to the nearest multiple of 360/128 degrees.
	trk int

	// Ground track value (decoded)
	GroundTrack string

	// Compact Position Reporting (CPR) Format (F)
	//
	// The CPR format field for the surface position message shall be encoded as specified for the airborne message.
	//  F = CprFormatEvenFrame = even format coding
	//  F = CprFormatOddFrame = odd format coding
	F common.CprFormat

	// Time Synchronization (T)
	//
	// This field shall indicate whether or not the time of applicability of the message is synchronized with UTC time.
	//
	//  TimeNotSynchronizedToUTC = time is not synchronized to UTC.
	//  TimeSynchronizedToUTC = time of applicability is synchronized to UTC time.
	//
	// Synchronization shall only be used for surface position messages having the top two horizontal position
	// precision categories (format TYPE Codes 5 and 6).
	//
	// When T = TimeSynchronizedToUTC, the time of validity in the surface position message format shall be encoded in the 1-bit F field which,
	// in addition to CPR format type, indicates the 0.2-second time tick for UTC time of position validity.
	// The F bit shall alternate between 0 and 1 for successive 0.2-second time ticks, beginning with F = 0 when the time of
	// applicability is an exact even-numbered UTC second.
	T TimeSynchronization

	// Latitude
	//
	// The CPR-encoded latitude value.
	LatCpr int

	// Longitude
	//
	// The CPR-encoded longitude value
	LonCpr int
}

// Ground Track Status
type GroundTrackStatus uint8

const GroundTrackStatusInvalid = GroundTrackStatus(0)
const GroundTrackStatusValid = GroundTrackStatus(1)

func DecodeBDS06(mb []byte) (frame BDS06Frame, err error) {
	// decode surface position bytes into struct
	// https://mode-s.org/decode/content/ads-b/4-surface-position.html

	frame = BDS06Frame{}

	frame.tc = (int(mb[0]) & 0b11111000) >> 3 // Type Code

	if frame.tc < 5 || frame.tc > 8 {
		err = errors.New("type code not from 5 to 8")
		return
	}

	frame.mov = (((int(mb[0]) & 0b00000111) << 4) + ((int(mb[1]) & 0b11110000) >> 4))
	frame.S = GroundTrackStatus((int(mb[1]) & 0b00001000) >> 3)
	frame.trk = (((int(mb[1]) & 0b00000111) << 4) + (int(mb[2]) & 0b11110000 >> 4))

	switch (int(mb[2]) & 0b00001000) >> 3 {
	case 0:
		frame.T = TimeNotSynchronizedToUTC
	case 1:
		frame.T = TimeSynchronizedToUTC
	}

	switch (int(mb[2]) & 0b00000100) >> 2 {
	case 0:
		frame.F = common.CprFormatEvenFrame
	case 1:
		frame.F = common.CprFormatOddFrame
	}

	frame.LatCpr = (((int(mb[2]) & 0b00000011) << 15) + (int(mb[3]) << 7) + ((int(mb[4]) & 0b11111110) >> 1))
	frame.LonCpr = (((int(mb[4]) & 0b00000001) << 16) + (int(mb[5]) << 8) + int(mb[6]))

	frame.GroundTrack = DecodeBDS05GroundTrack(frame.S, frame.trk)
	frame.GroundSpeed = DecodeBDS05SurfaceMovementSpeed(frame.mov)

	return
}

func DecodeBDS05GroundTrack(s GroundTrackStatus, trk int) (groundTrack string) {
	// https://mode-s.org/decode/content/ads-b/4-surface-position.html#ground-track

	switch s {
	case GroundTrackStatusInvalid:
		groundTrack = "ground track status invalid"
	case GroundTrackStatusValid:
		hdg := int(math.Round((360 * float64(trk)) / 128))
		groundTrack = fmt.Sprintf("%d°", hdg)
	}
	return
}

func DecodeBDS05SurfaceMovementSpeed(mov int) (speed string) {
	// https://mode-s.org/decode/content/ads-b/4-surface-position.html#movement

	var speedKnots, speedKmh float64

	switch {
	case mov == 0:
		speed = "No information available"
		return
	case mov == 1:
		speed = "Aircraft stopped (< 0.2315 km/h (0.125 kt))"
		return
	case mov >= 2 && mov <= 8:
		speedKnots = 0.125 + ((float64(mov) - 2) * 0.125)
	case mov >= 9 && mov <= 12:
		speedKnots = 1 + ((float64(mov) - 9) * 0.25)
	case mov >= 13 && mov <= 38:
		speedKnots = 2 + ((float64(mov) - 13) * 0.5)
	case mov >= 39 && mov <= 93:
		speedKnots = 15 + (float64(mov) - 39)
	case mov >= 94 && mov <= 108:
		speedKnots = 70 + ((float64(mov) - 94) * 2)
	case mov >= 109 && mov <= 123:
		speedKnots = 100 + ((float64(mov) - 109) * 5)
	case mov >= 124:
		speed = "≥ 324.1 km/h (175 kt)"
		return
	}
	speedKmh = 1.852 * speedKnots
	speed = fmt.Sprintf("%.4f km/h (%.3f kt)", speedKmh, speedKnots)
	return
}
