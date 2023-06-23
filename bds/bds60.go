package bds

// BDS6,0 Heading and speed report

import (
	"errors"
	"math"
)

type BDS60Frame struct {
	MagneticHeadingValid bool
	MagneticHeading      float64

	IndicatedAirspeedValid bool
	IndicatedAirspeed      float64

	MachNumberValid bool
	MachNumber      float64

	BarometricAltitudeRateValid bool
	BarometricAltitudeRate      float64

	GNSSAltitudeRateValid bool
	GNSSAltitudeRate      float64
}

func DecodeBDS60(mb []byte) (frame BDS60Frame, err error) {

	// magnetic heading
	switch (int(mb[0]) & 0b10000000) >> 7 {
	case 0:
		frame.MagneticHeadingValid = false
	case 1:
		frame.MagneticHeadingValid = true
		frame.MagneticHeading, err = decodeBDS60magneticHeading(mb)
		if err != nil {
			return
		}
	}

	// indicated airspeed
	switch (int(mb[1]) & 0b00001000) >> 3 {
	case 0:
		frame.IndicatedAirspeedValid = false
	case 1:
		frame.IndicatedAirspeedValid = true
		frame.IndicatedAirspeed, err = decodeBDS60indicatedAirspeed(mb)
		if err != nil {
			return
		}
	}

	// mach number
	switch int(mb[2]) & 0b00000001 {
	case 0:
		frame.MachNumberValid = false
	case 1:
		frame.MachNumberValid = true
		frame.MachNumber, err = decodeBDS60machNumber(mb)
		if err != nil {
			return
		}
	}

	// Barometric altitude rate
	switch (int(mb[4]) & 0b00100000) >> 5 {
	case 0:
		frame.BarometricAltitudeRateValid = false
	case 1:
		frame.BarometricAltitudeRateValid = true
		frame.BarometricAltitudeRate, err = decodeBDS60barometricAltitudeRate(mb)
		if err != nil {
			return
		}
	}

	// Inertial vertical velocity
	switch (int(mb[5]) & 0b00000100) >> 2 {
	case 0:
		frame.GNSSAltitudeRateValid = false
	case 1:
		frame.GNSSAltitudeRateValid = true
		frame.GNSSAltitudeRate, err = decodeBDS60GNSSAltitudeRate(mb)
		if err != nil {
			return
		}
	}

	return
}

func decodeBDS60magneticHeading(mb []byte) (magneticHeading float64, err error) {

	if (int(mb[0])&0b10000000)>>7 != 1 {
		err = errors.New("invalid status for status (for magnetic heading)")
		return
	}

	sign := int(mb[0] & 0b10000000 >> 7)
	magneticHeading = float64(((int(mb[0])&0b00111111)<<4)+((int(mb[1])&0b11110000)>>4)) * (90.0 / 512.0)

	if sign != 0 {
		magneticHeading = 180 + magneticHeading
	}

	return
}

func decodeBDS60indicatedAirspeed(mb []byte) (indicatedAirspeedKnots float64, err error) {
	// decode indicated airspeed from BDS 6,0 message
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#heading-and-speed-report-bds-60

	// check ias status bit is set
	if int(mb[1])&0b00001000 == 0 {
		err = errors.New("no true airspeed data present")
		return
	}

	indicatedAirspeedKnots = float64(((int(mb[1]) & 0b00000111) << 7) + ((int(mb[2]) & 0b11111110) >> 1))

	if indicatedAirspeedKnots < 0 || indicatedAirspeedKnots > 1023 {
		err = errors.New("indicated airspeed out of range [0,1023]")
	}

	return
}

func decodeBDS60machNumber(mb []byte) (machNumber float64, err error) {
	// decode mach number from BDS 6,0 message
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#heading-and-speed-report-bds-60

	// check mach number status bit is set
	if int(mb[2])&0b00000001 == 0 {
		err = errors.New("no mach number data present")
		return
	}

	machNumber = float64((int(mb[3])<<2)+((int(mb[4])&0b11000000)>>6)) * 0.004

	if machNumber < 0 || machNumber > 4.092 {
		err = errors.New("mach number out of range [0,4.092]")
	}

	return

}

func decodeBDS60barometricAltitudeRate(mb []byte) (barometricAltitudeRate float64, err error) {
	// decode barometric altitude rate from BDS 6,0 message
	// barometricAltitudeRate units are ft/min
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#heading-and-speed-report-bds-60

	// check barometric altitude rate status bit is set
	if int(mb[4])&0b00100000 == 0 {
		err = errors.New("no barometric altitude rate data present")
		return
	}

	sign := int((mb[4])&0b00010000) >> 4
	barometricAltitudeRate = (float64(((int(mb[4])&0b00001111)<<5)+((int(mb[5])&0b11111000)>>3)) - math.Pow(2, 9)) * 32

	if sign != 1 {
		barometricAltitudeRate *= -1
	}

	if barometricAltitudeRate < -16384 || barometricAltitudeRate > 16352 {
		err = errors.New("barometric altitude rate out of range [-16384, +16352]")
	}

	return

}

func decodeBDS60GNSSAltitudeRate(mb []byte) (inertialVerticalVelocity float64, err error) {
	// decode inertial vertical velocity from BDS 6,0 message
	// inertialVerticalVelocity units are ft/min
	// https://mode-s.org/decode/content/mode-s/7-ehs.html#heading-and-speed-report-bds-60

	// check intertial vertical velocity status bit is set
	if int(mb[5])&0b00000100 == 0 {
		err = errors.New("no intertial vertical velocity data present")
		return
	}

	sign := (int(mb[5]) & 0b00000010) >> 1
	inertialVerticalVelocity = (float64(((int(mb[5])&0b00000001)<<8)+(int(mb[6]))) - math.Pow(2, 9)) * 32

	if sign != 1 {
		inertialVerticalVelocity *= -1
	}

	if inertialVerticalVelocity < -16384 || inertialVerticalVelocity > 16352 {
		err = errors.New("intertial vertical velocity out of range [-16384, +16352]")
	}

	return inertialVerticalVelocity, err

}
