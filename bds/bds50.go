package bds

// BDS5,0 Track and turn report

import (
	"errors"
	"math"
)

type BDS50Frame struct {
	RollAngleValid bool
	RollAngle      float64

	TrueTrackAngleValid bool
	TrueTrackAngle      float64

	GroundSpeedValid bool
	GroundSpeed      float64

	TrueTrackAngleRateValid bool
	TrueTrackAngleRate      float64

	TrueAirspeedValid bool
	TrueAirspeed      float64
}

func DecodeBDS50(mb []byte) (frame BDS50Frame, err error) {

	// Roll Angle
	switch (int(mb[0]) & 0b10000000) >> 7 {
	case 0:
		frame.RollAngleValid = false
	case 1:
		frame.RollAngleValid = true
		frame.RollAngle, err = decodeBDS50roll(mb)
		if err != nil {
			return
		}
	}

	// True track angle
	switch (int(mb[1]) & 0b00010000) >> 4 {
	case 0:
		frame.TrueTrackAngleValid = false
	case 1:
		frame.TrueTrackAngleValid = true
		frame.TrueTrackAngle, err = decodeBDS50trueTrackAngle(mb)
		if err != nil {
			return
		}
	}

	// Ground speed
	switch int(mb[2]) & 0b00000001 {
	case 0:
		frame.GroundSpeedValid = false
	case 1:
		frame.GroundSpeedValid = true
		frame.GroundSpeed, err = decodeBDS50groundSpeed(mb)
		if err != nil {
			return
		}
	}

	// Track angle rate
	switch int(mb[4]) & 0b00100000 >> 5 {
	case 0:
		frame.TrueTrackAngleRateValid = false
	case 1:
		frame.TrueTrackAngleRateValid = true
		frame.TrueTrackAngleRate, err = decodeBDS50trackAngleRate(mb)
		if err != nil {
			return
		}
	}

	// True airspeed
	switch int(mb[5]) & 0b00000100 >> 2 {
	case 0:
		frame.TrueAirspeedValid = false
	case 1:
		frame.TrueAirspeedValid = true
		frame.TrueAirspeed, err = decodeBDS50trueAirspeed(mb)
		if err != nil {
			return
		}

	}

	return
}

func decodeBDS50trueTrackAngle(mb []byte) (trueTrackAngle float64, err error) {
	// decode true track angle (ground track) from BDS 5,0 message
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#track-and-turn-report-bds-50

	if (int(mb[1])&0b00010000)>>4 == 0 {
		err = errors.New("no true track angle data present")
		return
	}

	sign := (int(mb[1]) & 0b00001000) >> 3
	trueTrackAngle = float64(((int(mb[1])&0b00000111)<<7)+((int(mb[2])&0b11111110)>>1)) * (90.0 / 512.0)

	if sign != 0 {
		trueTrackAngle = 180 + trueTrackAngle
	}

	return

}

func decodeBDS50roll(mb []byte) (rollAngleDegrees float64, err error) {
	// decode roll angle from BDS 5,0 message
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#track-and-turn-report-bds-50

	// check roll angle status bit is set
	if int(int(mb[0])&0b10000000) == 0 {
		err = errors.New("no roll angle data present")
		return
	}

	sign := (int(mb[0]) & 0b01000000) >> 6
	rollAngleDegrees = (float64(((int(mb[0])&0b00111111)<<3)+((int(mb[1])&0b11100000)>>5)) - math.Pow(2, 9)) * (45.0 / 256.0)

	// log.Debug().Int("sign", sign).Msg("sign")

	if sign == 0 {
		rollAngleDegrees = (rollAngleDegrees + 90)
	}

	if rollAngleDegrees > 90 || rollAngleDegrees < -90 {
		err = errors.New("roll angle out of range [-90,+90] degrees")
	}

	return
}

func decodeBDS50groundSpeed(mb []byte) (groundSpeedKnots float64, err error) {
	// decode true airspeed from BDS 5,0 message
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#track-and-turn-report-bds-50

	// check true ground speed status bit is set
	if int(mb[2])&0b00000001 == 0 {
		err = errors.New("no ground speed data present")
		return
	}

	groundSpeedKnots = float64((int(mb[3])<<2)+((int(mb[4])&0b11000000)>>6)) * 2

	if groundSpeedKnots < 0 || groundSpeedKnots > 2046 {
		err = errors.New("true ground speed out of range [0,2046]")
	}

	return

}

func decodeBDS50trackAngleRate(mb []byte) (trackAngleRate float64, err error) {
	// decode track angle rate from BDS 5,0 message
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#track-and-turn-report-bds-50

	// check track angle rate status bit is set
	if (int(mb[4])&0b00100000)>>5 == 0 {
		err = errors.New("no track angle rate data present")
		return
	}

	sign := (int(mb[6]) & 0b00010000) >> 4
	trackAngleRate = float64(((int(mb[4])&0b00001111)<<5)+((int(mb[5])&0b11111000)>>3)) * (8.0 / 256.0)

	if sign != 0 {
		trackAngleRate = -trackAngleRate
	}

	return
}

func decodeBDS50trueAirspeed(mb []byte) (trueAirspeedKnots float64, err error) {
	// decode true airspeed from BDS 5,0 message
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#track-and-turn-report-bds-50

	// check true airspeed status bit is set
	if int(mb[5])&0b00000100 == 0 {
		err = errors.New("no true airspeed data present")
		return
	}

	trueAirspeedKnots = float64(((int(mb[5])&0b00000011)<<8)+(int(mb[6]))) * 2

	if trueAirspeedKnots < 0 || trueAirspeedKnots > 2046 {
		err = errors.New("true airspeed out of range [0,2046]")
	}

	return

}
