package bds

import (
	"errors"
	"math"
)

func decodeBDS45turbulence(mb []byte) (turbulence int, err error) {
	// decode turbulence data from BDS4,5 message
	// https://mode-s.org/decode/content/mode-s/8-meteo.html#meteorological-hazard-report-bds-45

	if (int(mb[0])&0b10000000)>>7 == 0 {
		err = errors.New("no turbulence data present")
	}

	turbulence = (int(mb[0]) & 0b01100000) >> 5

	return
}

func decodeBDS45windShear(mb []byte) (windShear int, err error) {
	// decode wind shear data from BDS4,5 message
	// https://mode-s.org/decode/content/mode-s/8-meteo.html#meteorological-hazard-report-bds-45

	if (int(mb[0])&0b00010000)>>4 == 0 {
		err = errors.New("no wind shear data present")
		return
	}

	windShear = (int(mb[0]) & 0b00001100) >> 2

	return
}

func decodeBDS45staticAirTemperature(mb []byte) (staticAirTemperature float64, err error) {
	// decode static air temperature data from BDS4,5 message
	// https://mode-s.org/decode/content/mode-s/8-meteo.html#meteorological-hazard-report-bds-45

	if (int(mb[1]) & 0b00000001) == 0 {
		err = errors.New("no static air temperature data present")
		return
	}

	sign := (int(mb[2]) & 0b10000000) >> 7

	staticAirTemperature = (float64(((int(mb[2])&0b01111111)<<2)+((int(mb[3])&0b11000000)>>6)) - math.Pow(2, 9)) * 0.25

	if sign != 1 {
		staticAirTemperature *= -1
	}

	return
}
