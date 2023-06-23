package main

// import (
// 	"fmt"
// 	"math"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/umahmood/haversine"
// )

// func TestAltitudeFromAltitudeCode13bit(t *testing.T) {
// 	var testTable = []struct {
// 		data          []byte
// 		altFt         float64
// 		expectedError bool
// 	}{
// 		{
// 			data:          []byte{0x20, 0x00, 0x06, 0x1b, 0x2d, 0xd4, 0xdd},
// 			altFt:         8875,
// 			expectedError: false,
// 		},
// 		{
// 			data:          []byte{0x20, 0x00, 0x02, 0xbd, 0x19, 0x1f, 0x7e},
// 			altFt:         3725,
// 			expectedError: false,
// 		},
// 		{
// 			data:          []byte{0x20, 0x00, 0x00, 0x00, 0xfc, 0x71, 0x05},
// 			expectedError: true,
// 		},
// 	}

// 	assert := assert.New(t)
// 	for _, testData := range testTable {
// 		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
// 		msg := decodeDF4(testData.data)
// 		altFt, err := altitudeFromAltitudeCode13bit(msg.ac)
// 		if testData.expectedError {
// 			assert.Error(err, testMsg+"altitudeFromAltitudeCode13bit error expected")
// 		} else {
// 			assert.NoError(err, testMsg+"altitudeFromAltitudeCode13bit error")
// 		}
// 		assert.Equal(testData.altFt, altFt, testMsg+"altFt")
// 	}
// }

// // TODO:
// // need to capture data for testing
// // func TestCalcAirSpeedAndHeading(t *testing.T) {
// // }

// func TestCalcAirborneLatLonGloballyUnambiguous(t *testing.T) {
// 	var testTable = []struct {
// 		dataOdd     []byte
// 		dataEven    []byte
// 		latCprOdd   int
// 		latCprEven  int
// 		lonCprOdd   int
// 		lonCprEven  int
// 		latestFrame string
// 		lat, lon    float64
// 	}{
// 		{
// 			// *8d7c534d5813a7182997f866edbc;
// 			// CRC: 000000
// 			// RSSI: -20.5 dBFS
// 			// Time: 426397071011.58us
// 			// DF:17 AA:7C534D CA:5 ME:5813A7182997F8
// 			//  Extended Squitter Airborne position (barometric altitude) (11)
// 			//   ICAO Address:  7C534D (Mode S / ADS-B)
// 			//   Air/Ground:    airborne
// 			//   Baro altitude: 2850 ft
// 			//   CPR type:      Airborne
// 			//   CPR odd flag:  odd
// 			//   CPR latitude:  (101396)
// 			//   CPR longitude: (104440)
// 			//   CPR decoding:  none
// 			//   NIC-B:         0
// 			//   NACp:          8
// 			//   SIL:           2 (p <= 0.001%, unknown type)

// 			// *8d7c534d5813a2bd7a3d0e834424;
// 			// CRC: 000000
// 			// RSSI: -14.7 dBFS
// 			// Time: 426397559097.42us
// 			// DF:17 AA:7C534D CA:5 ME:5813A2BD7A3D0E
// 			//  Extended Squitter Airborne position (barometric altitude) (11)
// 			//   ICAO Address:  7C534D (Mode S / ADS-B)
// 			//   Air/Ground:    airborne
// 			//   Baro altitude: 2850 ft
// 			//   CPR type:      Airborne
// 			//   CPR odd flag:  even
// 			//   CPR latitude:  -31.88979 (89789)
// 			//   CPR longitude: 116.05858 (15630)
// 			//   CPR decoding:  global
// 			//   NIC:           8
// 			//   Rc:            0.186 km / 0.1 NM
// 			//   NIC-B:         0
// 			//   NACp:          8
// 			//   SIL:           2 (p <= 0.001%, unknown type)
// 			dataOdd:     []byte{0x58, 0x13, 0xA7, 0x18, 0x29, 0x97, 0xF8},
// 			dataEven:    []byte{0x58, 0x13, 0xA2, 0xBD, 0x7A, 0x3D, 0x0E},
// 			latCprOdd:   101396,
// 			latCprEven:  89789,
// 			lonCprOdd:   104440,
// 			lonCprEven:  15630,
// 			latestFrame: "even",
// 			lat:         -31.88979,
// 			lon:         116.05858,
// 		},
// 		{
// 			// *8d76e7245831b31544dd371b123f;
// 			// CRC: 000000
// 			// RSSI: -32.6 dBFS
// 			// Time: 426397165639.58us
// 			// DF:17 AA:76E724 CA:5 ME:5831B31544DD37
// 			// Extended Squitter Airborne position (barometric altitude) (11)
// 			// ICAO Address:  76E724 (Mode S / ADS-B)
// 			// Air/Ground:    airborne
// 			// Baro altitude: 8875 ft
// 			// CPR type:      Airborne
// 			// CPR odd flag:  even
// 			// CPR latitude:  (101026)
// 			// CPR longitude: (56631)
// 			// CPR decoding:  none
// 			// NIC-B:         0
// 			// NACp:          8
// 			// SIL:           2 (p <= 0.001%, unknown type)

// 			// *8d76e7245831b76eb8383f2717b9;
// 			// CRC: 000000
// 			// RSSI: -30.1 dBFS
// 			// Time: 426397722652.75us
// 			// DF:17 AA:76E724 CA:5 ME:5831B76EB8383F
// 			// Extended Squitter Airborne position (barometric altitude) (11)
// 			// ICAO Address:  76E724 (Mode S / ADS-B)
// 			// Air/Ground:    airborne
// 			// Baro altitude: 8875 ft
// 			// CPR type:      Airborne
// 			// CPR odd flag:  odd
// 			// CPR latitude:  -31.37416 (112476)
// 			// CPR longitude: 115.99096 (14399)
// 			// CPR decoding:  global
// 			// NIC:           8
// 			// Rc:            0.186 km / 0.1 NM
// 			// NIC-B:         0
// 			// NACp:          8
// 			// SIL:           2 (p <= 0.001%, unknown type)
// 			dataOdd:     []byte{0x58, 0x31, 0xB7, 0x6E, 0xB8, 0x38, 0x3F},
// 			dataEven:    []byte{0x58, 0x31, 0xB3, 0x15, 0x44, 0xDD, 0x37},
// 			latCprOdd:   112476,
// 			latCprEven:  101026,
// 			lonCprOdd:   14399,
// 			lonCprEven:  56631,
// 			latestFrame: "odd",
// 			lat:         -31.37416,
// 			lon:         115.99096,
// 		},
// 	}

// 	assert := assert.New(t)
// 	for _, testData := range testTable {
// 		testMsg := fmt.Sprintf("dataOdd: %014x, dataEven: %014x, ", testData.dataOdd, testData.dataEven)

// 		oddFrame, err := decodeBDS05(testData.dataOdd)
// 		assert.NoError(err, testMsg+"decodeBDS05 error odd frame")
// 		assert.Equal(cprFormatOddFrame, oddFrame.f)

// 		evenFrame, err := decodeBDS05(testData.dataEven)
// 		assert.NoError(err, testMsg+"decodeBDS05 error even frame")
// 		assert.Equal(cprFormatEvenFrame, evenFrame.f)

// 		latEven, latOdd := calcAirborneLatGloballyUnambiguous(float64(evenFrame.latCpr), float64(oddFrame.latCpr))

// 		NLeven := calcLongitudeZoneNumber(latEven)
// 		NLodd := calcLongitudeZoneNumber(latOdd)
// 		assert.EqualValues(NLeven, NLodd, testMsg+"NL for odd and even values not equal")

// 		lonEven, lonOdd := calcAirborneLonGloballyUnambiguous(float64(evenFrame.lonCpr), float64(oddFrame.lonCpr), NLeven)

// 		switch testData.latestFrame {
// 		case "odd":
// 			assert.Equal(testData.lat, math.Round(latOdd*100000)/100000, testMsg+"latOdd")
// 			assert.Equal(testData.lon, math.Round(lonOdd*100000)/100000, testMsg+"lonOdd")
// 		case "even":
// 			assert.Equal(testData.lat, math.Round(latEven*100000)/100000, testMsg+"latEven")
// 			assert.Equal(testData.lon, math.Round(lonEven*100000)/100000, testMsg+"lonEven")
// 		default:
// 			t.Error("unknown latestFrame in testData")
// 		}
// 	}
// }

// func TestCalcAirborneLatLonLocallyUnambiguous(t *testing.T) {

// 	var testTable = []struct {
// 		data           []byte
// 		latCpr         int
// 		lonCpr         int
// 		f              cprFormat
// 		refLat, refLon float64
// 		lat, lon       float64
// 	}{
// 		{

// 			// *8d7c0a2b581f16c5bb893a062273;
// 			// CRC: 000000
// 			// RSSI: -34.2 dBFS
// 			// Time: 426410647233.17us
// 			// DF:17 AA:7C0A2B CA:5 ME:581F16C5BB893A
// 			//  Extended Squitter Airborne position (barometric altitude) (11)
// 			//   ICAO Address:  7C0A2B (Mode S / ADS-B)
// 			//   Air/Ground:    airborne
// 			//   Baro altitude: 5025 ft
// 			//   CPR type:      Airborne
// 			//   CPR odd flag:  odd
// 			//   CPR latitude:  -32.38113 (90845)
// 			//   CPR longitude: 115.84668 (100666)
// 			//   CPR decoding:  global
// 			//   NIC:           8
// 			//   Rc:            0.186 km / 0.1 NM
// 			//   NIC-B:         0

// 			// *8d7c0a2b581de6c67b89b2e1c251;
// 			// CRC: 000000
// 			// RSSI: -36.1 dBFS
// 			// Time: 426424095414.50us
// 			// DF:17 AA:7C0A2B CA:5 ME:581DE6C67B89B2
// 			//  Extended Squitter Airborne position (barometric altitude) (11)
// 			//   ICAO Address:  7C0A2B (Mode S / ADS-B)
// 			//   Air/Ground:    airborne
// 			//   Baro altitude: 4950 ft
// 			//   CPR type:      Airborne
// 			//   CPR odd flag:  odd
// 			//   CPR latitude:  -32.37666 (90941)
// 			//   CPR longitude: 115.85341 (100786)
// 			//   CPR decoding:  local
// 			//   NIC:           8
// 			//   Rc:            0.186 km / 0.1 NM
// 			//   NIC-B:         0

// 			data:   []byte{0x58, 0x1D, 0xE6, 0xC6, 0x7B, 0x89, 0xB2},
// 			latCpr: 90941,
// 			lonCpr: 100786,
// 			f:      cprFormatOddFrame,
// 			refLat: -32.38113, // Previous position
// 			refLon: 115.84668, // Previous position
// 			lat:    -32.37666,
// 			lon:    115.85341,
// 		},
// 	}

// 	assert := assert.New(t)
// 	for _, testData := range testTable {
// 		testMsg := fmt.Sprintf("data: %014x, ", testData.data)

// 		previousPosition := haversine.Coord{Lat: testData.refLat, Lon: testData.refLon}
// 		currentPosition := haversine.Coord{Lat: testData.lat, Lon: testData.lon}

// 		_, km := haversine.Distance(previousPosition, currentPosition)
// 		nm := km * 0.539957
// 		assert.LessOrEqual(nm, float64(180), testMsg+"The reference position should be close to the actual position, which must be within a 180 NM range.")

// 		lat := calcAirborneLatLocallyUnambiguous(testData.refLat, float64(testData.latCpr), testData.f)
// 		NL := calcLongitudeZoneNumber(lat)
// 		lon := calcAirborneLonLocallyUnambiguous(testData.refLon, float64(testData.lonCpr), NL, testData.f)

// 		assert.Equal(testData.lat, math.Round(lat*100000)/100000, testMsg+"lat")
// 		assert.Equal(testData.lon, math.Round(lon*100000)/100000, testMsg+"lon")
// 	}
// }

// func TestCalcSurfaceLatLonGloballyUnambiguous(t *testing.T) {

// 	// *907cf5cc3020043d0a4711e51ac9;
// 	// CRC: 000000
// 	// RSSI: -16.8 dBFS
// 	// Time: 16082399992415.17us
// 	// DF:18 AA:7CF5CC CF:0 ME:3020043D0A4711
// 	//  Extended Squitter (Non-Transponder) Surface position (6)
// 	//   ICAO Address:  7CF5CC (ADS-B, non-transponder)
// 	//   Air/Ground:    ground
// 	//   Groundspeed:   0.2 kt (v2: 0.1 kt)
// 	//   CPR type:      Surface
// 	//   CPR odd flag:  odd
// 	//   CPR latitude:  (7813)
// 	//   CPR longitude: (18193)
// 	//   CPR decoding:  none
// 	//   NACp:          10
// 	//   SIL:           2 (p <= 0.001%, unknown type)

// 	// *907cf5cc302002d198dace16a694;
// 	// CRC: 000000
// 	// RSSI: -16.8 dBFS
// 	// Time: 16082404851754.08us
// 	// DF:18 AA:7CF5CC CF:0 ME:302002D198DACE
// 	//  Extended Squitter (Non-Transponder) Surface position (6)
// 	//   ICAO Address:  7CF5CC (ADS-B, non-transponder)
// 	//   Air/Ground:    ground
// 	//   Groundspeed:   0.2 kt (v2: 0.1 kt)
// 	//   CPR type:      Surface
// 	//   CPR odd flag:  even
// 	//   CPR latitude:  (92364)
// 	//   CPR longitude: (56014)
// 	//   CPR decoding:  none
// 	//   NACp:          10
// 	//   SIL:           2 (p <= 0.001%, unknown type)

// 	var testTable = []struct {
// 		dataOdd        []byte
// 		dataEven       []byte
// 		latCprOdd      int
// 		latCprEven     int
// 		lonCprOdd      int
// 		lonCprEven     int
// 		latestFrame    string
// 		refLat, refLon float64
// 		lat, lon       float64
// 	}{
// 		{
// 			dataOdd:     []byte{0x30, 0x20, 0x04, 0x3D, 0x0A, 0x47, 0x11},
// 			dataEven:    []byte{0x30, 0x20, 0x02, 0xD1, 0x98, 0xDA, 0xCE},
// 			latCprOdd:   7813,
// 			latCprEven:  92364,
// 			lonCprOdd:   18193,
// 			lonCprEven:  56014,
// 			latestFrame: "even",
// 			refLat:      -31.9487, // YPPH lat
// 			refLon:      115.9733, // YPPH lon
// 			lat:         -31.9430,
// 			lon:         115.9692,
// 		},
// 	}

// 	assert := assert.New(t)
// 	for _, testData := range testTable {
// 		testMsg := fmt.Sprintf("dataOdd: %014x, dataEven: %014x, ", testData.dataOdd, testData.dataEven)

// 		oddFrame, err := decodeBDS06(testData.dataOdd)
// 		assert.NoError(err, testMsg+"decodeBDS06 error odd frame")
// 		assert.Equal(cprFormatOddFrame, oddFrame.f, testMsg+"odd f")

// 		evenFrame, err := decodeBDS06(testData.dataEven)
// 		assert.NoError(err, testMsg+"decodeBDS06 error even frame")
// 		assert.Equal(cprFormatEvenFrame, evenFrame.f, testMsg+"even f")

// 		latEven, latOdd := calcSurfaceLatGloballyUnambiguous(testData.refLat, float64(testData.latCprEven), float64(testData.latCprOdd))

// 		NLeven := calcLongitudeZoneNumber(latEven)
// 		NLodd := calcLongitudeZoneNumber(latOdd)
// 		assert.Equal(NLeven, NLodd, testMsg+"NLeven should equal NLodd")

// 		lonEven, lonOdd := calcSurfaceLonGloballyUnambiguous(testData.refLon, float64(testData.lonCprEven), float64(testData.lonCprOdd), NLeven)

// 		switch testData.latestFrame {
// 		case "odd":
// 			assert.Equal(testData.lat, math.Round(latEven*10000)/10000, testMsg+"lat")
// 			assert.Equal(testData.lon, math.Round(lonEven*10000)/10000, testMsg+"lon")
// 		case "even":
// 			assert.Equal(testData.lat, math.Round(latOdd*10000)/10000, testMsg+"lat")
// 			assert.Equal(testData.lon, math.Round(lonOdd*10000)/10000, testMsg+"lon")
// 		}
// 	}
// }

// func TestCalcSurfaceLatLonLocallyUnambiguous(t *testing.T) {

// 	var testTable = []struct {
// 		data           []byte
// 		latCpr         int
// 		lonCpr         int
// 		f              cprFormat
// 		refLat, refLon float64
// 		lat, lon       float64
// 	}{
// 		{

// 			// *907cf5cc3020043d0a4711e51ac9;
// 			// CRC: 000000
// 			// RSSI: -16.8 dBFS
// 			// Time: 16082399992415.17us
// 			// DF:18 AA:7CF5CC CF:0 ME:3020043D0A4711
// 			//  Extended Squitter (Non-Transponder) Surface position (6)
// 			//   ICAO Address:  7CF5CC (ADS-B, non-transponder)
// 			//   Air/Ground:    ground
// 			//   Groundspeed:   0.2 kt (v2: 0.1 kt)
// 			//   CPR type:      Surface
// 			//   CPR odd flag:  odd
// 			//   CPR latitude:  (7813)
// 			//   CPR longitude: (18193)
// 			//   CPR decoding:  none
// 			//   NACp:          10
// 			//   SIL:           2 (p <= 0.001%, unknown type)

// 			data:   []byte{0x30, 0x20, 0x04, 0x3D, 0x0A, 0x47, 0x11},
// 			latCpr: 7813,
// 			lonCpr: 18193,
// 			f:      cprFormatOddFrame,
// 			refLat: -31.94297, // Previous position
// 			refLon: 115.96923, // Previous position
// 			lat:    -31.94297,
// 			lon:    115.96923,
// 		},
// 	}

// 	assert := assert.New(t)
// 	for _, testData := range testTable {
// 		testMsg := fmt.Sprintf("data: %014x, ", testData.data)

// 		previousPosition := haversine.Coord{Lat: testData.refLat, Lon: testData.refLon}
// 		currentPosition := haversine.Coord{Lat: testData.lat, Lon: testData.lon}

// 		_, km := haversine.Distance(previousPosition, currentPosition)
// 		nm := km * 0.539957
// 		assert.LessOrEqual(nm, float64(180), testMsg+"The reference position should be close to the actual position, which must be within a 180 NM range.")

// 		lat := calcSurfaceLatLocallyUnambiguous(testData.refLat, float64(testData.latCpr), testData.f)
// 		NL := calcLongitudeZoneNumber(lat)
// 		lon := calcSurfaceLonLocallyUnambiguous(testData.refLon, float64(testData.lonCpr), NL, testData.f)

// 		assert.Equal(testData.lat, math.Round(lat*100000)/100000, testMsg+"lat")
// 		assert.Equal(testData.lon, math.Round(lon*100000)/100000, testMsg+"lon")
// 	}
// }

// func TestCalcSurfaceMovementSpeed(t *testing.T) {

// 	var testTable = []struct {
// 		data []byte
// 		gs   float64
// 	}{
// 		{
// 			// *907cf66b41dcc2d70cd8b5374da0;
// 			// CRC: 000000
// 			// RSSI: -16.5 dBFS
// 			// Time: 16082400083962.25us
// 			// DF:18 AA:7CF66B CF:0 ME:41DCC2D70CD8B5
// 			//  Extended Squitter (Non-Transponder) Surface position (8)
// 			//   ICAO Address:  7CF66B (ADS-B, non-transponder)
// 			//   Air/Ground:    ground
// 			//   Track/Heading  213.8
// 			//   Groundspeed:   10.2 kt
// 			//   CPR type:      Surface
// 			//   CPR odd flag:  even
// 			//   CPR latitude:  (93062)
// 			//   CPR longitude: (55477)
// 			//   CPR decoding:  none
// 			//   NACp:          0
// 			//   SIL:           2 (p <= 0.001%, unknown type)
// 			data: []byte{0x41, 0xDC, 0xC2, 0xD7, 0x0C, 0xD8, 0xB5},
// 			gs:   10.00,
// 		},
// 		{
// 			// *8c7c40053a3ef445384494579d82;
// 			// CRC: 000000
// 			// RSSI: -9.3 dBFS
// 			// Time: 16082400418458.92us
// 			// DF:17 AA:7C4005 CA:4 ME:3A3EF445384494
// 			//  Extended Squitter Surface position (7)
// 			//   ICAO Address:  7C4005 (Mode S / ADS-B)
// 			//   Air/Ground:    ground
// 			//   Track/Heading  312.2
// 			//   Groundspeed:   13.2 kt
// 			//   CPR type:      Surface
// 			//   CPR odd flag:  odd
// 			//   CPR latitude:  (8860)
// 			//   CPR longitude: (17556)
// 			//   CPR decoding:  none
// 			data: []byte{0x3A, 0x3E, 0xF4, 0x45, 0x38, 0x44, 0x94},
// 			gs:   13.00,
// 		},
// 	}
// 	assert := assert.New(t)
// 	for _, testData := range testTable {
// 		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
// 		frame, err := decodeBDS06(testData.data)
// 		assert.NoError(err, testMsg+"decodeBDS06 error")

// 		gs, err := calcSurfaceMovementSpeed(frame.mov)
// 		assert.NoError(err, testMsg+"calcSurfaceMovementSpeed error")

// 		assert.Equal(testData.gs, gs, testMsg+"gs")
// 	}
// }

// func TestCalcVerticalRate(t *testing.T) {

// 	var testTable = []struct {
// 		data []byte
// 		vr   int
// 	}{
// 		{
// 			// *8d7c1ac19944c283682c01f69467;
// 			// CRC: 000000
// 			// RSSI: -23.5 dBFS
// 			// Time: 426397000361.58us
// 			// DF:17 AA:7C1AC1 CA:5 ME:9944C283682C01
// 			//  Extended Squitter Airborne velocity over ground, subsonic (19/1)
// 			//   ICAO Address:  7C1AC1 (Mode S / ADS-B)
// 			//   Air/Ground:    airborne
// 			//   Geom - baro:   0 ft
// 			//   Ground track   262.3
// 			//   Groundspeed:   194.7 kt
// 			//   Geom rate:     -640 ft/min
// 			//   NACv:          0
// 			data: []byte{0x99, 0x44, 0xC2, 0x83, 0x68, 0x2C, 0x01},
// 			vr:   -640,
// 		},
// 		{
// 			// *8d7c803399145f1cc85c81f5b8af;
// 			// CRC: 000000
// 			// RSSI: -25.9 dBFS
// 			// Time: 426397061573.58us
// 			// DF:17 AA:7C8033 CA:5 ME:99145F1CC85C81
// 			//  Extended Squitter Airborne velocity over ground, subsonic (19/1)
// 			//   ICAO Address:  7C8033 (Mode S / ADS-B)
// 			//   Air/Ground:    airborne
// 			//   Geom - baro:   0 ft
// 			//   Ground track   337.7
// 			//   Groundspeed:   247.5 kt
// 			//   Geom rate:     -1408 ft/min
// 			//   NACv:          2
// 			data: []byte{0x99, 0x14, 0x5F, 0x1C, 0xC8, 0x5C, 0x81},
// 			vr:   -1408,
// 		},
// 	}

// 	assert := assert.New(t)
// 	for _, testData := range testTable {
// 		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
// 		frame, err := decodeBDS09(testData.data)
// 		assert.NoError(err, testMsg+"decodeBDS09 error")
// 		assert.Equal(testData.vr, frame.vr, testMsg+"vr")
// 	}
// }

// //------------

// func TestGetDF(t *testing.T) {
// 	data := []byte{0x2A, 0x00, 0x51, 0x6D, 0x49, 0x2B, 0x80}
// 	df := getDF(data)
// 	assert.Equal(t, DF5, df)
// }

// func TestIcaoFromCRC(t *testing.T) {
// 	data := []byte{0xA0, 0x00, 0x18, 0x38, 0xCA, 0x38, 0x00, 0x31, 0x44, 0x00, 0x00, 0xF2, 0x41, 0x77}
// 	df := getDF(data)
// 	assert.Equal(t, DF20, df)
// 	icao := icaoFromCRC(data)
// 	assert.Equal(t, 0x3C6DD0, icao)
// }

// func TestSquawkFromIdentityCode(t *testing.T) {
// 	data := []byte{0x2A, 0x00, 0x51, 0x6D, 0x49, 0x2B, 0x80}
// 	msg := decodeDF5(data)
// 	squawk := squawkFromIdentityCode(msg.id)
// 	assert.Equal(t, 356, squawk)
// }
