package bds

import (
	"errors"
	"math"
)

func decodeBDS44figureOfMerit(mb []byte) (fom int) {
	// decode figure of merit from BDS 4,4 message
	// https://mode-s.org/decode/content/mode-s/8-meteo.html#meteorological-routine-air-report-bds-44
	fom = int(mb[0]&0b11110000) >> 4
	return
}

func decodeBDS44windSpeedDirection(mb []byte) (windSpeed, windDirection float64, err error) {
	// decode wind speed (knots) and direction (degrees) from BDS 4,4 message
	// https://mode-s.org/decode/content/mode-s/8-meteo.html#meteorological-routine-air-report-bds-44
	if mb[0]&0b00001000 == 0 {
		err = errors.New("no wind speed or wind direction data present")
		return
	}

	windSpeed = float64(((int(mb[0]) & 0b00000111) << 6) + ((int(mb[1]) & 0b11111100) >> 2))

	if windSpeed < 0 || windSpeed > 511 {
		err = errors.New("wind speed out of range [0,511]")
		return
	}

	windDirection = float64(((int(mb[1])&0b00000011)<<7)+((int(mb[2])&0b11111110)>>1)) * (180.0 / 256.0)

	if windDirection < 0 || windDirection > 360 {
		err = errors.New("wind direction out of range [0,360]")
		return
	}

	return
}

func decodeBDS44staticAirTemperature(mb []byte) (staticAirTemperature float64, err error) {
	// decode static air temperature (degrees C) from BDS 4,4 message
	// https://mode-s.org/decode/content/mode-s/8-meteo.html#meteorological-routine-air-report-bds-44

	sign := int(mb[2]) & 0b00000001
	staticAirTemperature = (float64(((int(mb[3]))<<2)+((int(mb[4])&0b11000000)>>6)) - math.Pow(2, 10)) * 0.25

	if sign != 1 {
		staticAirTemperature *= -1
	}

	if staticAirTemperature < -80 || staticAirTemperature > 60 {
		err = errors.New("static air temperature out of range [-80,60]")
	}
	return
}
