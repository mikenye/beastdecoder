package bds

// Airborne velocity
// Type Code (TC): 19
// https://mode-s.org/decode/content/ads-b/5-airborne-velocity.html

import (
	"errors"
	"math"
)

type BDS09Frame struct {
	// https://mode-s.org/decode/content/ads-b/5-airborne-velocity.html
	tc int // Type Code

	// Subtypes
	//
	// Subtypes 1 (AirborneVelocityGroundSpeedSubsonic) and 2 (AirborneVelocityGroundSpeedSupersonic) of the airborne velocity format
	// shall be used when the transmitting aircraft’s velocity over ground is known.
	// Subtype 1 (AirborneVelocityGroundSpeedSubsonic) shall be used at subsonic velocities while subtype 2 (AirborneVelocityGroundSpeedSupersonic)
	// shall be used when the velocity exceeds 1022 kt.
	//
	// Subtypes 3 (AirborneVelocityAirSpeedSubsonic) and 4 (AirborneVelocityAirSpeedSupersonic) of the airborne velocity \
	// format shall be used when the transmitting aircraft’s velocity over ground is not known.
	// These subtypes substitute airspeed and heading for the velocity over ground.
	// Subtype 3 (AirborneVelocityAirSpeedSubsonic) shall be used at subsonic velocities, while subtype 4 (AirborneVelocityAirSpeedSupersonic)
	// shall be used when the velocity exceeds 1022 kt.
	st AirborneVelocitySubType

	// Intent Change Flag in Airborne Velocity Messages
	//
	// An intent change event shall be triggered 4 seconds after the detection of new information in:
	//  - Selected vertical intention
	//  - Next waypoint position
	// The code shall remain set for 18 ±1 second following an intent change.
	ic bool

	// IFR Capability Flag (IFR) in Airborne Velocity Messages
	//
	// The IFR capability is a subfield in the subtypes 1, 2, 3 and 4 airborne velocity messages.
	// IFR = true shall signify that the transmitting aircraft has a capability for applications requiring ADS-B equipage class A1 or above.
	// Otherwise, IFR shall be set to false.
	ifr bool

	// Sub-type specific fields for ground speed
	groundSpeedFields     BDS09FrameGroundSpeed
	groundSpeed           float64
	groundTrack           float64
	groundSpeedTrackValid bool

	// Sub-type specific fields for airspeed
	airSpeedFields     BDS09FrameAirSpeed
	airSpeed           float64
	airTrack           float64
	airSpeedTrackValid bool

	nuc int // Navigation uncertainty category for velocity

	vrSrc verticalRateSource // Source bit for vertical rate
	svr   verticalRateSign   // Sign bit for vertical rate
	vr    int                // vertical rate
	sDif  int                // Sign bit for GNSS and Baro altitudes difference

	// Difference From Barometric Altitude in Airborne Velocity Messages
	//
	// This field shall contain the signed difference between barometric and GNSS altitude.
	//
	// The difference between barometric altitude and GNSS height above ellipsoid (HAE) shall be used if available.
	// If GNSS HAE is not available, GNSS altitude (MSL) shall be used when airborne position is being reported using format TYPE Codes 11 through 18.
	//
	// If airborne position is being reported using format TYPE Code 9 or 10, only GNSS (HAE) shall be used.
	// For format TYPE Code 9 or 10, if GNSS (HAE) is not available, the field shall be coded with all zeros.
	// The basis for the barometric altitude difference (either GNSS (HAE) or GNSS altitude MSL) shall be used consistently for the reported difference.
	dAlt int
}

type BDS09FrameGroundSpeed struct {
	// https://mode-s.org/decode/content/ads-b/5-airborne-velocity.html#sub-type-1-and-2-ground-speed-decoding
	dew int // Direction for E-W velocity component
	vew int // East-West velocity component
	dns int // Direction for N-S velocity component
	vns int // North-South velocity component
}

type BDS09FrameAirSpeed struct {
	// Magnetic heading status
	//
	// This field shall define the availability of the magnetic heading value.
	// Coding for this field shall be:
	//   false = magnetic heading data not available
	//    true = magnetic heading data available
	sh int

	// Magnetic heading value
	//
	// This field shall contain the aircraft magnetic heading (in degrees clockwise from magnetic north) when velocity over ground is not available.
	// The magnetic heading shall be encoded as an unsigned angular weighted binary numeral with an MSB of 180 degrees and an LSB of 360/1024 degrees,
	// with zero indicating magnetic north. The data in the field shall be rounded to the nearest multiple of 360/1024 degrees.
	hdg int // Magnetic heading

	t  int // Airspeed type
	as int // Airspeed
}

// Airborne Velocity Subtype
type AirborneVelocitySubType uint8

const AirborneVelocityGroundSpeedSubsonic = AirborneVelocitySubType(1)
const AirborneVelocityGroundSpeedSupersonic = AirborneVelocitySubType(2)
const AirborneVelocityAirSpeedSubsonic = AirborneVelocitySubType(3)
const AirborneVelocityAirSpeedSupersonic = AirborneVelocitySubType(4)

// Vertical Rate Source
type verticalRateSource uint8

const verticalRateSourceGNSS = verticalRateSource(0)
const verticalRateSourceBaro = verticalRateSource(1)

// Vertical Rate Sign
type verticalRateSign uint8

const verticalRateClimb = verticalRateSign(0)
const verticalRateDescent = verticalRateSign(1)

// --------------------

func DecodeBDS09(mb []byte) (frame BDS09Frame, err error) {

	frame = BDS09Frame{}

	frame.tc = (int(mb[0]) & 0b11111000) >> 3

	if frame.tc != 19 {
		err = errors.New("type code not 19")
		return
	}

	frame.st = AirborneVelocitySubType(int(mb[0]) & 0b00000111)
	switch (int(mb[1]) & 0b10000000) >> 7 {
	case 0:
		frame.ic = false
	case 1:
		frame.ic = true
	}

	switch (int(mb[1]) & 0b01000000) >> 6 {
	case 0:
		frame.ifr = false
	case 1:
		frame.ifr = true
	}

	frame.nuc = ((int(mb[1]) & 0b00111000) >> 3)

	subTypeBits := (int(mb[1]) & 0b00000111 << 19) + (int(mb[2]) << 11) + (int(mb[3]) << 3) + ((int(mb[4]) & 0b11100000) >> 5)

	frame.vrSrc = verticalRateSource((int(mb[4]) & 0b00010000) >> 4)
	frame.svr = verticalRateSign((int(mb[4]) & 0b00001000) >> 3)
	switch frame.svr {
	case 0:
		frame.vr = 64 * ((((int(mb[4]) & 0b00000111) << 6) + ((int(mb[5]) & 0b11111100) >> 2)) - 1)
	case 1:
		frame.vr = -64 * ((((int(mb[4]) & 0b00000111) << 6) + ((int(mb[5]) & 0b11111100) >> 2)) - 1)
	}

	// reserved bits
	// if (int(mb[5]) & 0b00000011) != 0 {
	// 	err = errors.New("reserved bits not zero")
	// 	return
	// }

	frame.sDif = ((int(mb[6]) & 0b10000000) >> 7)
	switch frame.sDif {
	case 0:
		frame.dAlt = ((int(mb[6]) & 0b01111111) - 1) * 25
	case 1:
		frame.dAlt = -((int(mb[6]) & 0b01111111) - 1) * 25
	}

	switch {
	case frame.st <= 2:
		frame.groundSpeedFields = BDS09FrameGroundSpeed{
			dew: ((subTypeBits & 0b1000000000000000000000) >> 21),
			vew: ((subTypeBits & 0b0111111111100000000000) >> 11),
			dns: ((subTypeBits & 0b0000000000010000000000) >> 10),
			vns: (subTypeBits & 0b0000000000001111111111),
		}
		frame.groundSpeed, frame.groundTrack, _, _ = calcGroundSpeedAndHeading(frame.st, &frame.groundSpeedFields)
		frame.groundSpeedTrackValid = true
	case frame.st <= 4:
		frame.airSpeedFields = BDS09FrameAirSpeed{
			sh:  ((subTypeBits & 0b1000000000000000000000) >> 21),
			hdg: ((subTypeBits & 0b0111111111100000000000) >> 11),
			t:   ((subTypeBits & 0b0000000000010000000000) >> 10),
			as:  (subTypeBits & 0b0000000000001111111111),
		}
		frame.airSpeed, frame.airTrack = calcAirSpeedAndHeading(frame.st, &frame.airSpeedFields)
		frame.airSpeedTrackValid = true
	}

	return
}

func calcAirSpeedAndHeading(st AirborneVelocitySubType, m *BDS09FrameAirSpeed) (vas, mh float64) {
	// https://mode-s.org/decode/content/ads-b/5-airborne-velocity.html#sub-type-3-and-4-airspeed-decoding

	// magnetic heading
	mh = float64(m.hdg) * (360.0 / 1024.0)

	// airspeed
	switch st {
	case AirborneVelocityAirSpeedSubsonic:
		vas = (float64(m.as) - 1)
	case AirborneVelocityAirSpeedSupersonic:
		vas = 4 * (float64(m.as) - 1)
	}
	return vas, mh
}

func calcGroundSpeedAndHeading(st AirborneVelocitySubType, m *BDS09FrameGroundSpeed) (gs, gta, vx, vy float64) {
	// https://mode-s.org/decode/content/ads-b/5-airborne-velocity.html#sub-type-1-and-2-ground-speed-decoding

	switch st {
	case AirborneVelocityGroundSpeedSubsonic:
		switch m.dew {
		case 0:
			vx = float64((m.vew) - 1)
		case 1:
			vx = float64(-1 * ((m.vew) - 1))
		}
		switch m.dns {
		case 0:
			vy = float64((m.vns) - 1)
		case 1:
			vy = float64(-1 * ((m.vns) - 1))
		}
	case AirborneVelocityGroundSpeedSupersonic:
		switch m.dew {
		case 0:
			vx = float64(4 * ((m.vew) - 1))
		case 1:
			vx = float64(-4 * ((m.vew) - 1))
		}
		switch m.dns {
		case 0:
			vy = float64(4 * ((m.vns) - 1))
		case 1:
			vy = float64(-4 * ((m.vns) - 1))
		}
	}

	// ground speed
	gs = math.Sqrt((math.Pow(vx, 2) + math.Pow(vy, 2)))

	// ground track angle
	gta = math.Atan2(vx, vy) * (360 / (2 * math.Pi))
	if gta < 0 {
		gta += 360
	}

	return gs, gta, vx, vy
}
