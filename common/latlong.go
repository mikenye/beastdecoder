package common

import (
	"math"

	"github.com/rs/zerolog/log"
)

func AirborneLatGloballyUnambiguous(latCprEven, latCprOdd float64) (latEven, latOdd float64) {
	// calculate globally unambigouous latitude
	// https://mode-s.org/decode/content/ads-b/3-airborne-position.html#sec:cpr_airborne_global_lat

	lce := latCprEven / math.Pow(2, 17)
	lco := latCprOdd / math.Pow(2, 17)

	// latitude zone sizes
	dLatEven := 360 / (4 * float64(Nz))
	dLatOdd := 360 / ((4 * float64(Nz)) - 1)

	// calculate the latitude zone index
	j := math.Floor((59 * lce) - (60 * lco) + 0.5)

	// computed latitudes
	latEven = dLatEven * (Modulo(j, 60) + lce)
	latOdd = dLatOdd * (Modulo(j, 59) + lco)

	// For the southern hemisphere, values returned from previous equations range from 270 to 360 degrees.
	if latEven >= 270 {
		latEven -= 360
	}
	if latOdd >= 270 {
		latOdd -= 360
	}

	return latEven, latOdd

}

func AirborneLonGloballyUnambiguous(lonCprEven, lonCprOdd, NL float64) (lonEven, lonOdd float64) {
	// calculate globally unambigouous longitude
	// https://mode-s.org/decode/content/ads-b/3-airborne-position.html#sec:cpr_airborne_global_lon

	lce := lonCprEven / math.Pow(2, 17)
	lco := lonCprOdd / math.Pow(2, 17)

	// calculate the longitude zone index
	m := math.Floor(lce*(NL-1) - lco*NL + 0.5)

	// calculate the number of longitude zones
	nEven := math.Max(NL, 1)
	nOdd := math.Max(NL-1, 1)

	// longitude zone sizes
	dLonEven := 360 / nEven
	dLonOdd := 360 / nOdd

	// calculate longitude
	lonEven = dLonEven * (Modulo(m, nEven) + lce)
	lonOdd = dLonOdd * (Modulo(m, nOdd) + lco)

	if lonEven < -180 {
		lonEven = 360 + lonEven
	}
	if lonEven > 180 {
		lonEven = lonEven - 360
	}
	if lonOdd < -180 {
		lonOdd = 360 + lonEven
	}
	if lonOdd > 180 {
		lonOdd = lonEven - 360
	}

	return lonEven, lonOdd

}

func AirborneLatLocallyUnambiguous(refLat, latCpr float64, f CprFormat) (lat float64) {
	// calculate locally unambiguous latitude
	// https://mode-s.org/decode/content/ads-b/3-airborne-position.html#calculation-of-latitude
	var dLat float64

	lc := latCpr / math.Pow(2, 17)

	// latitude zone size is different depending on the message type
	switch f {
	case CprFormatEvenFrame:
		dLat = 360.0 / (4.0 * Nz)
	case CprFormatOddFrame:
		dLat = 360 / (4*Nz - 1)
	}

	// latitude zone index
	j := math.Floor(refLat/dLat) + math.Floor(0.5+(Modulo(refLat, dLat)/dLat)-lc)

	// Knowing the latitude zone index, the latitude of the new position is
	lat = dLat * (j + lc)

	return lat

}

func AirborneLonLocallyUnambiguous(refLon, lonCpr, NL float64, f CprFormat) (lon float64) {
	// calculate locally unambiguous longitude
	// https://mode-s.org/decode/content/ads-b/3-airborne-position.html#calculation-of-longitude

	var dLon float64

	lc := lonCpr / math.Pow(2, 17)

	switch f {
	case CprFormatEvenFrame:
		dLon = 360 / (math.Max(NL, 1))

	case CprFormatOddFrame:
		dLon = 360 / (math.Max(NL-1, 1))

	}

	m := math.Floor(refLon/dLon) + math.Floor(0.5+(Modulo(refLon, dLon)/dLon)-lc)

	lon = dLon * (m + lc)

	if log.Debug().Enabled() {
		log.Debug().Float64("NL", NL).Float64("dLon", dLon).Float64("m", m).Msg("calcAirborneLonLocallyUnambiguous")
	}

	return lon

}

func SurfaceLatGloballyUnambiguous(refLat, latCprEven, latCprOdd float64) (latEven, latOdd float64) {
	// calculate globally unambigouous latitude
	// https://mode-s.org/decode/content/ads-b/4-surface-position.html

	lce := latCprEven / math.Pow(2, 17)
	lco := latCprOdd / math.Pow(2, 17)

	// latitude zone sizes
	dLatEven := 90 / (4 * float64(Nz))
	dLatOdd := 90 / ((4 * float64(Nz)) - 1)

	// calculate the latitude zone index
	j := math.Floor((59 * lce) - (60 * lco) + 0.5)

	// computed latitudes
	latEvenBase := dLatEven * (Modulo(j, 60) + lce)
	latOddBase := dLatOdd * (Modulo(j, 59) + lco)

	latEvenNorthernHemisphere := latEvenBase
	latEvenSouthernHemisphere := latEvenBase - 90
	latOddNorthernHemisphere := latOddBase
	latOddSouthernHemisphere := latOddBase - 90

	if refLat > 0 {
		return latEvenNorthernHemisphere, latOddNorthernHemisphere
	} else {
		return latEvenSouthernHemisphere, latOddSouthernHemisphere
	}
}

func SurfaceLatLocallyUnambiguous(refLat float64, latCpr float64, f CprFormat) (lat float64) {
	// calculate locally unambiguous latitude
	// https://mode-s.org/decode/content/ads-b/4-surface-position.html

	var dLat float64

	lc := latCpr / math.Pow(2, 17)

	// latitude zone size is different depending on the message type
	switch f {
	case CprFormatEvenFrame:
		dLat = 90 / (4 * Nz)
	case CprFormatOddFrame:
		dLat = 90 / (4*Nz - 1)
	}

	// latitude zone index
	j := math.Floor(refLat/dLat) + math.Floor(((Modulo(refLat, dLat) / dLat) - lc + 0.5))

	// Knowing the latitude zone index, the latitude of the new position is
	lat = dLat * (j + lc)

	// log.Debug().Float64("dLat", dLat).Float64("j", j).Msg("tshoot")

	return lat

}

func SurfaceLonGloballyUnambiguous(refLon, lonCprEven, lonCprOdd, NL float64) (lonEven, lonOdd float64) {
	// calculate globally unambigouous longitude
	// https://mode-s.org/decode/content/ads-b/4-surface-position.html

	lce := lonCprEven / math.Pow(2, 17)
	lco := lonCprOdd / math.Pow(2, 17)

	// calculate the longitude zone index
	m := math.Floor(lce*(NL-1) - lco*NL + 0.5)

	// calculate the number of longitude zones
	nEven := math.Max(NL, 1)
	nOdd := math.Max(NL-1, 1)

	// longitude zone sizes
	dLonEven := 90 / nEven
	dLonOdd := 90 / nOdd

	// calculate longitude
	lonEvenBase := dLonEven * (Modulo(m, nEven) + lce)
	lonOddBase := dLonOdd * (Modulo(m, nOdd) + lco)

	// four possible longitude solutions
	possibleLonsEven := []float64{
		lonEvenBase,
		lonEvenBase + 90,
		lonEvenBase + 180,
		lonEvenBase + 270,
	}

	possibleLonsOdd := []float64{
		lonOddBase,
		lonOddBase + 90,
		lonOddBase + 180,
		lonOddBase + 270,
	}

	// make sure lons are between -180 and 180
	for i := range possibleLonsEven {
		if possibleLonsEven[i] >= 180 {
			possibleLonsEven[i] -= 360
		}
	}
	for i := range possibleLonsOdd {
		if possibleLonsOdd[i] >= 180 {
			possibleLonsOdd[i] -= 360
		}
	}

	// the closest solution to receiver is the correct one
	distancesEven := []float64{
		math.Abs(refLon - possibleLonsEven[0]),
		math.Abs(refLon - possibleLonsEven[1]),
		math.Abs(refLon - possibleLonsEven[2]),
		math.Abs(refLon - possibleLonsEven[3]),
	}

	distancesOdd := []float64{
		math.Abs(refLon - possibleLonsOdd[0]),
		math.Abs(refLon - possibleLonsOdd[1]),
		math.Abs(refLon - possibleLonsOdd[2]),
		math.Abs(refLon - possibleLonsOdd[3]),
	}

	distanceIndex := 0
	distanceMin := float64(0)
	for i, e := range distancesEven {
		if i == 0 || e < distanceMin {
			distanceIndex = i
			distanceMin = e
		}
	}
	lonEven = possibleLonsEven[distanceIndex]

	distanceIndex = 0
	distanceMin = float64(0)
	for i, e := range distancesOdd {
		if i == 0 || e < distanceMin {
			distanceIndex = i
			distanceMin = e
		}
	}
	lonOdd = possibleLonsOdd[distanceIndex]

	return lonEven, lonOdd

}

func SurfaceLonLocallyUnambiguous(refLon float64, lonCpr, NL float64, f CprFormat) (lon float64) {
	// calculate locally unambiguous longitude
	// https://mode-s.org/decode/content/ads-b/4-surface-position.html

	var dLon float64

	lc := lonCpr / math.Pow(2, 17)

	switch f {
	case CprFormatEvenFrame:
		dLon = 90 / (math.Max(NL, 1))
	case CprFormatOddFrame:
		dLon = 90 / (math.Max(NL-1, 1))
	}

	m := math.Floor(refLon/dLon) + math.Floor((Modulo(refLon, dLon)/dLon)-lc+0.5)

	lon = dLon * (m + lc)

	// log.Debug().Float64("dLon", dLon).Float64("m", m).Msg("tshoot")

	return lon

}

func LongitudeZoneNumber(lat float64) (NL float64) {
	// calculate Longitude zone number
	// Given the latitude, this function yields the number of longitude zones between 1 and 59.
	switch lat := lat; {
	case lat == 0:
		return float64(59)
	case lat == 87:
		return float64(2)
	case lat == -87:
		return float64(2)
	case lat > 87:
		return float64(1)
	case lat < -87:
		return float64(1)
	default:
		return math.Floor((2 * math.Pi) / (math.Acos(1 - (1-math.Cos(math.Pi/(2*Nz)))/(math.Pow(math.Cos((math.Pi/180)*lat), 2)))))
	}
}
