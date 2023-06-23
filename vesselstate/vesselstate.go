package vesselstate

import (
	"beastdecoder/bds"
	"beastdecoder/common"
	"beastdecoder/df"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/umahmood/haversine"
)

type VesselState struct {
	mu sync.RWMutex // sync mutex

	MsgCount int // message count

	// Vessel Information - Squawk Code
	SquawkCodeKnown bool
	SquawkCode      int

	// Vessel Information - Callsign
	CallsignKnown bool
	Callsign      string

	// Position Information - Airborne Status
	AirborneStatusKnown bool
	Airborne            bool

	// Airborne Position Information - Altitude
	AltitudeKnown bool
	Altitude      int

	// Position Information
	Lat, Lon         float64
	LatLonMethod     string // "global" / "local"
	LatLonKnown      bool
	LastPositionData time.Time

	// Storing Airborne Position odd/even lat/lon CPR
	airborneLatCprOdd, airborneLonCprOdd   int // lats/lons/NL used for actual lat/lon calculation
	airborneLatLonCprOddKnown              bool
	airborneLatCprEven, airborneLonCprEven int // lats/lons/NL used for actual lat/lon calculation
	airborneLatLonCprEvenKnown             bool
	airborneLatLonCprTypeHist              []common.CprFormat

	// Storing Surface Position odd/even lat/lon CPR
	surfaceLatCprOdd, surfaceLonCprOdd   int // lats/lons/NL used for actual lat/lon calculation
	surfaceLatLonCprOddKnown             bool
	surfaceLatCprEven, surfaceLonCprEven int // lats/lons/NL used for actual lat/lon calculation
	surfaceLatLonCprEvenKnown            bool
	surfaceLatLonCprTypeHist             []common.CprFormat

	// Ground Speed (mov)
	GroundSpeed      string
	GroundSpeedKnown bool

	// Ground Track (trk)
	GroundTrack      string
	GroundTrackKnown bool

	// Last message received from vessel
	LastUpdated time.Time
}

type Vessels struct {
	mu      sync.RWMutex         // sync mutex
	Vessels map[int]*VesselState // map of vessels, key is ICAO

	// reference lat/lon for location calculations
	refLatLonKnown bool
	refLat, refLon float64
}

func (vdb *Vessels) RLock() {
	vdb.mu.RLock()
}

func (vdb *Vessels) RUnlock() {
	vdb.mu.RUnlock()
}

func (vdb *Vessels) Init() {
	// run once on program start to init the vessel db
	vdb.Vessels = make(map[int]*VesselState)
	go vdb.evictor()
}

func (vdb *Vessels) SetRefLatLon(refLat, refLon float64) {
	vdb.mu.Lock()
	defer vdb.mu.Unlock()
	vdb.refLat = refLat
	vdb.refLon = refLon
	vdb.refLatLonKnown = true
}

func (vdb *Vessels) evictor() {
	// evicts stale (no updates in >60sec) entries from vdb
	for {
		time.Sleep(time.Second * 1)

		icaosToEvict := []int{}

		// find expired icaos (no updates in 60 sec)
		vdb.mu.RLock()
		for icao := range vdb.Vessels {

			clearPositionData := false

			vdb.Vessels[icao].mu.RLock()
			// determine of record should be evicted
			if time.Now().Sub(vdb.Vessels[icao].LastUpdated) > (time.Second * 60) {
				icaosToEvict = append(icaosToEvict, icao)
			} else {
				// In the event that the navigation input ceases, the extrapolation described in
				// §A.2.3.2.3.1 and §A.2.3.2.3.2 shall be limited to no more than 2 seconds.
				//
				// At the end of this time-out of 2 seconds,
				// all fields of the airborne position register,
				// except the altitude field, shall be cleared (set to zero).
				//
				// When the appropriate register fields are cleared,
				// the zero TYPE Code field shall serve to notify ADS-B receiving equipment
				// that the data in the latitude and longitude fields are invalid.
				if time.Now().Sub(vdb.Vessels[icao].LastPositionData) > (time.Second * 2) {
					clearPositionData = true
				}
			}
			vdb.Vessels[icao].mu.RUnlock()

			// clear position data if needed
			if clearPositionData {
				vdb.clearPositionData(icao)
			}
		}
		vdb.mu.RUnlock()

		// delete them
		vdb.mu.Lock()
		for _, icao := range icaosToEvict {
			if log.Debug().Enabled() {
				log.Info().Str("icao", fmt.Sprintf("%06x", icao)).Msg("removing expired")
			}
			delete(vdb.Vessels, icao)
		}
		vdb.mu.Unlock()

	}
}

func (vdb *Vessels) incrementMessageCount(icao int) {
	// increments the message count for an icao

	// check if vessel is in db
	// if it isn't just bail out
	if !vdb.isVesselTracked(icao) {
		return
	}

	// increment counter
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()
	vdb.Vessels[icao].MsgCount++
}

func (vdb *Vessels) isVesselTracked(icao int) bool {
	// check if vessel is in db
	vdb.mu.RLock()
	_, ok := vdb.Vessels[icao]
	defer vdb.mu.RUnlock()
	if !ok {
		return false
	}
	return true
}

func (vdb *Vessels) addVessel(icao int) {
	// adds a vessel to the vessel db if it does not yet exist
	if icao == 0 {
		return
	}

	// if vessel is not in db, add it
	if !vdb.isVesselTracked(icao) {
		vdb.mu.Lock()
		if log.Debug().Enabled() {
			log.Debug().Str("icao", fmt.Sprintf("%06x", icao)).Msg("addVessel")
		}
		vdb.Vessels[icao] = &VesselState{
			LastUpdated: time.Now(),
			// StoredFrames: make(map[cprFormat]map[BDScode]interface{}),
		}
		vdb.mu.Unlock()
	} else {
		vdb.updateLastSeen(icao)
	}
}

func (vdb *Vessels) updateLastSeen(icao int) {
	// sets vessel last seen to now
	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return
	}
	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()
	// set time
	vdb.Vessels[icao].LastUpdated = time.Now()
	if log.Debug().Enabled() {
		log.Debug().Time("LastUpdated", vdb.Vessels[icao].LastUpdated).Str("icao", fmt.Sprintf("%06x", icao)).Msg("updateLastSeen")
	}
}

func (vdb *Vessels) setCallsign(icao int, callsign string) {
	// set ground speed
	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return
	}
	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()
	// set ground speed
	vdb.Vessels[icao].Callsign = callsign
	vdb.Vessels[icao].CallsignKnown = true
}

func (vdb *Vessels) setGroundSpeed(icao int, groundSpeed string) {
	// set ground speed
	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return
	}
	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()
	// set ground speed
	vdb.Vessels[icao].GroundSpeed = groundSpeed
	vdb.Vessels[icao].GroundSpeedKnown = true
}

func (vdb *Vessels) setGroundTrack(icao int, groundTrack string) {
	// set ground speed
	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return
	}
	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()
	// set ground speed
	vdb.Vessels[icao].GroundTrack = groundTrack
	vdb.Vessels[icao].GroundTrackKnown = true
}

func (vdb *Vessels) setSquawkCode(icao int, squawk int) {
	// sets airborne status
	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return
	}
	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()
	vdb.Vessels[icao].SquawkCodeKnown = true
	vdb.Vessels[icao].SquawkCode = squawk
	if log.Debug().Enabled() {
		log.Debug().Bool("SquawkCodeKnown", true).Int("SquawkCode", squawk).Str("icao", fmt.Sprintf("%06x", icao)).Msg("setSquawkCode")
	}
}

func (vdb *Vessels) setAirborneStatus(icao int, airborne bool) {
	// sets airborne status
	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return
	}
	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()
	vdb.Vessels[icao].AirborneStatusKnown = true
	vdb.Vessels[icao].Airborne = airborne
	if log.Debug().Enabled() {
		log.Debug().Bool("AirborneStatusKnown", true).Bool("Airborne", airborne).Str("icao", fmt.Sprintf("%06x", icao)).Msg("setAirborneStatus")
	}
}

func (vdb *Vessels) setAltitude(icao int, altitude int) {
	// sets airborne status
	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return
	}
	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()

	// sanity check
	if vdb.Vessels[icao].AltitudeKnown {
		if math.Abs(float64(vdb.Vessels[icao].Altitude-altitude)) >= 2000 {
			log.Warn().Bool("AltitudeKnown", true).Int("OldAltitude", vdb.Vessels[icao].Altitude).Int("NewAltitude", altitude).Str("icao", fmt.Sprintf("%06x", icao)).Msg("sanity check fail on altitude change")
			return
		}
	}

	// set altitude
	vdb.Vessels[icao].AltitudeKnown = true
	vdb.Vessels[icao].Altitude = altitude
	if log.Debug().Enabled() {
		log.Debug().Bool("AltitudeKnown", true).Int("Altitude", altitude).Str("icao", fmt.Sprintf("%06x", icao)).Msg("setAltitude")
	}
}

func (vdb *Vessels) clearPositionData(icao int) {
	if !vdb.isVesselTracked(icao) {
		return
	}
	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()

	vdb.Vessels[icao].LatLonKnown = false

	vdb.Vessels[icao].airborneLatLonCprOddKnown = false
	vdb.Vessels[icao].airborneLatLonCprEvenKnown = false

	vdb.Vessels[icao].surfaceLatLonCprOddKnown = false
	vdb.Vessels[icao].surfaceLatLonCprEvenKnown = false

	vdb.Vessels[icao].airborneLatLonCprTypeHist = make([]common.CprFormat, 3)
	vdb.Vessels[icao].surfaceLatLonCprTypeHist = make([]common.CprFormat, 3)
}

func (vdb *Vessels) storeAirborneLatLonCPR(icao int, latCpr, lonCpr int, f common.CprFormat) {
	// stores latCpr/lonCpr in the vessel db entry
	// used for position decoding, where we need both an "odd" and "even" frames to be able to determine position accurately

	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return
	}

	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()

	// store the data
	switch f {
	case common.CprFormatEvenFrame:
		vdb.Vessels[icao].airborneLatCprEven = latCpr
		vdb.Vessels[icao].airborneLonCprEven = lonCpr
		vdb.Vessels[icao].airborneLatLonCprEvenKnown = true
	case common.CprFormatOddFrame:
		vdb.Vessels[icao].airborneLatCprOdd = latCpr
		vdb.Vessels[icao].airborneLonCprOdd = lonCpr
		vdb.Vessels[icao].airborneLatLonCprOddKnown = true
	}
	vdb.Vessels[icao].LastPositionData = time.Now()
	vdb.Vessels[icao].airborneLatLonCprTypeHist = append(vdb.Vessels[icao].airborneLatLonCprTypeHist, f)

	// trim airborneLatLonCprTypeHist
	if len(vdb.Vessels[icao].airborneLatLonCprTypeHist) > 3 {
		vdb.Vessels[icao].airborneLatLonCprTypeHist = vdb.Vessels[icao].airborneLatLonCprTypeHist[len(vdb.Vessels[icao].airborneLatLonCprTypeHist)-3:]
	}
}

func (vdb *Vessels) storeSurfaceLatLonCPR(icao int, latCpr, lonCpr int, f common.CprFormat) {
	// stores latCpr/lonCpr in the vessel db entry
	// used for position decoding, where we need both an "odd" and "even" frames to be able to determine position accurately

	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return
	}

	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()

	// store the data
	switch f {
	case common.CprFormatEvenFrame:
		vdb.Vessels[icao].surfaceLatCprEven = latCpr
		vdb.Vessels[icao].surfaceLonCprEven = lonCpr
		vdb.Vessels[icao].surfaceLatLonCprEvenKnown = true
	case common.CprFormatOddFrame:
		vdb.Vessels[icao].surfaceLatCprOdd = latCpr
		vdb.Vessels[icao].surfaceLonCprOdd = lonCpr
		vdb.Vessels[icao].surfaceLatLonCprOddKnown = true
	}
	vdb.Vessels[icao].LastPositionData = time.Now()
	vdb.Vessels[icao].surfaceLatLonCprTypeHist = append(vdb.Vessels[icao].surfaceLatLonCprTypeHist, f)

	// trim airborneLatLonCprTypeHist
	if len(vdb.Vessels[icao].surfaceLatLonCprTypeHist) > 3 {
		vdb.Vessels[icao].surfaceLatLonCprTypeHist = vdb.Vessels[icao].surfaceLatLonCprTypeHist[len(vdb.Vessels[icao].surfaceLatLonCprTypeHist)-3:]
	}
}

func (vdb *Vessels) calculateAirbornePosition(icao int) error {

	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return nil
	}

	vdb.mu.RLock()
	defer vdb.mu.RUnlock()
	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()

	oldLat := vdb.Vessels[icao].Lat
	oldLon := vdb.Vessels[icao].Lon

	// TODO check times that lat/lons were received, to work out whether they are still valid

	// need at least two previous positions
	if len(vdb.Vessels[icao].airborneLatLonCprTypeHist) >= 2 {

		// for local decoding, need one CPR position and a previous lat/lon within 150NM
		// for global decoding: need odd and even data, and last two frames must be different format (ie: odd + even || even + odd frames)
		if vdb.Vessels[icao].airborneLatLonCprTypeHist[len(vdb.Vessels[icao].airborneLatLonCprTypeHist)-2] == vdb.Vessels[icao].airborneLatLonCprTypeHist[len(vdb.Vessels[icao].airborneLatLonCprTypeHist)-1] {

			// local if last position known
			if vdb.Vessels[icao].LatLonKnown {

				var lat, lon float64

				// calculate lat/lon
				switch vdb.Vessels[icao].airborneLatLonCprTypeHist[len(vdb.Vessels[icao].airborneLatLonCprTypeHist)-1] {
				case common.CprFormatEvenFrame:
					lat = common.AirborneLatLocallyUnambiguous(vdb.Vessels[icao].Lat, float64(vdb.Vessels[icao].airborneLatCprEven), common.CprFormatEvenFrame)
					NL := common.LongitudeZoneNumber(lat)
					lon = common.AirborneLonLocallyUnambiguous(vdb.Vessels[icao].Lon, float64(vdb.Vessels[icao].airborneLonCprEven), NL, common.CprFormatEvenFrame)
				case common.CprFormatOddFrame:
					lat = common.AirborneLatLocallyUnambiguous(vdb.Vessels[icao].Lat, float64(vdb.Vessels[icao].airborneLatCprOdd), common.CprFormatOddFrame)
					NL := common.LongitudeZoneNumber(lat)
					lon = common.AirborneLonLocallyUnambiguous(vdb.Vessels[icao].Lon, float64(vdb.Vessels[icao].airborneLonCprOdd), NL, common.CprFormatOddFrame)
				}

				// new position must be no more than 180NM from reference position
				posOld := haversine.Coord{Lat: vdb.Vessels[icao].Lat, Lon: vdb.Vessels[icao].Lon}
				posNew := haversine.Coord{Lat: lat, Lon: lon}
				_, km := haversine.Distance(posOld, posNew)
				nm := km * 0.539957
				if nm > 180 {
					return errors.New("local decoding, new position greater than 180NM from old position")
				}

				// update lat/lon
				vdb.Vessels[icao].Lat = lat
				vdb.Vessels[icao].Lon = lon
				vdb.Vessels[icao].LatLonMethod = "airborne,local"

				// Need two lat/lon within ~1NM of each other before saying that we know the position
				previousPosition := haversine.Coord{Lat: oldLat, Lon: oldLon}
				currentPosition := haversine.Coord{Lat: vdb.Vessels[icao].Lat, Lon: vdb.Vessels[icao].Lon}
				mi, _ := haversine.Distance(previousPosition, currentPosition)
				if mi < 1.15 {
					vdb.Vessels[icao].LatLonKnown = true
				}

				return nil

			} else {
				return errors.New("insufficient data for local position decoding")
			}
		} else {

			// global unambiguous decoding

			// calc latitude
			latEven, latOdd := common.AirborneLatGloballyUnambiguous(float64(vdb.Vessels[icao].airborneLatCprEven), float64(vdb.Vessels[icao].airborneLatCprOdd))

			// calc longitude zone number and check sanity
			NLeven := common.LongitudeZoneNumber(latEven)
			NLodd := common.LongitudeZoneNumber(latOdd)
			if NLeven != NLodd {
				return errors.New("longitude zone numbers do not match between even and odd cpr values")
			}

			// calculate longitude
			lonEven, lonOdd := common.AirborneLonGloballyUnambiguous(float64(vdb.Vessels[icao].airborneLonCprEven), float64(vdb.Vessels[icao].airborneLonCprOdd), NLeven)

			// work out which frame type to return
			switch vdb.Vessels[icao].airborneLatLonCprTypeHist[len(vdb.Vessels[icao].airborneLatLonCprTypeHist)-1] {
			case common.CprFormatEvenFrame:
				vdb.Vessels[icao].Lat = latEven
				vdb.Vessels[icao].Lon = lonEven
			case common.CprFormatOddFrame:
				vdb.Vessels[icao].Lat = latOdd
				vdb.Vessels[icao].Lon = lonOdd
			}
			vdb.Vessels[icao].LatLonMethod = "airborne,global"

			// Need two lat/lon within ~1NM of each other before saying that we know the position
			previousPosition := haversine.Coord{Lat: oldLat, Lon: oldLon}
			currentPosition := haversine.Coord{Lat: vdb.Vessels[icao].Lat, Lon: vdb.Vessels[icao].Lon}
			mi, _ := haversine.Distance(previousPosition, currentPosition)
			if mi < 1.15 {
				vdb.Vessels[icao].LatLonKnown = true
			}

			return nil
		}
	}
	return errors.New("not enough position messages received")
}

func (vdb *Vessels) calculateSurfacePosition(icao int) error {

	// ensure vessel exists before attempting to update
	if !vdb.isVesselTracked(icao) {
		return nil
	}

	vdb.mu.RLock()
	defer vdb.mu.RUnlock()

	if !vdb.refLatLonKnown {
		return errors.New("cannot decode surface position without receiver lat/lon")
	}

	vdb.Vessels[icao].mu.Lock()
	defer vdb.Vessels[icao].mu.Unlock()

	oldLat := vdb.Vessels[icao].Lat
	oldLon := vdb.Vessels[icao].Lon

	// TODO check times that lat/lons were received, to work out whether they are still valid

	// need at least two previous positions
	if len(vdb.Vessels[icao].surfaceLatLonCprTypeHist) >= 2 {

		// for local decoding, need one CPR position and a previous lat/lon within 150NM
		// for global decoding: need odd and even data, and last two frames must be different format (ie: odd + even || even + odd frames)
		if vdb.Vessels[icao].surfaceLatLonCprTypeHist[len(vdb.Vessels[icao].surfaceLatLonCprTypeHist)-2] == vdb.Vessels[icao].surfaceLatLonCprTypeHist[len(vdb.Vessels[icao].surfaceLatLonCprTypeHist)-1] {

			// local if last position known
			if vdb.Vessels[icao].LatLonKnown {

				var lat, lon float64

				// calculate lat/lon
				switch vdb.Vessels[icao].surfaceLatLonCprTypeHist[len(vdb.Vessels[icao].surfaceLatLonCprTypeHist)-1] {
				case common.CprFormatEvenFrame:
					lat = common.SurfaceLatLocallyUnambiguous(vdb.Vessels[icao].Lat, float64(vdb.Vessels[icao].surfaceLatCprEven), common.CprFormatEvenFrame)
					NL := common.LongitudeZoneNumber(lat)
					lon = common.SurfaceLonLocallyUnambiguous(vdb.Vessels[icao].Lon, float64(vdb.Vessels[icao].surfaceLonCprEven), NL, common.CprFormatEvenFrame)
				case common.CprFormatOddFrame:
					lat = common.SurfaceLatLocallyUnambiguous(vdb.Vessels[icao].Lat, float64(vdb.Vessels[icao].surfaceLatCprOdd), common.CprFormatOddFrame)
					NL := common.LongitudeZoneNumber(lat)
					lon = common.SurfaceLonLocallyUnambiguous(vdb.Vessels[icao].Lon, float64(vdb.Vessels[icao].surfaceLonCprOdd), NL, common.CprFormatOddFrame)
				}

				// new position must be no more than 180NM from reference position
				posOld := haversine.Coord{Lat: vdb.Vessels[icao].Lat, Lon: vdb.Vessels[icao].Lon}
				posNew := haversine.Coord{Lat: lat, Lon: lon}
				_, km := haversine.Distance(posOld, posNew)
				nm := km * 0.539957
				if nm > 180 {
					return errors.New("local decoding, new position greater than 180NM from old position")
				}

				// update lat/lon
				vdb.Vessels[icao].Lat = lat
				vdb.Vessels[icao].Lon = lon
				vdb.Vessels[icao].LatLonMethod = "surface,local"

				// Need two lat/lon within ~1NM of each other before saying that we know the position
				previousPosition := haversine.Coord{Lat: oldLat, Lon: oldLon}
				currentPosition := haversine.Coord{Lat: vdb.Vessels[icao].Lat, Lon: vdb.Vessels[icao].Lon}
				mi, _ := haversine.Distance(previousPosition, currentPosition)
				if mi < 1.15 {
					vdb.Vessels[icao].LatLonKnown = true
				}

				return nil

			} else {
				return errors.New("insufficient data for local position decoding")
			}
		} else {

			// global unambiguous decoding

			// calc latitude
			latEven, latOdd := common.SurfaceLatGloballyUnambiguous(vdb.refLat, float64(vdb.Vessels[icao].surfaceLatCprEven), float64(vdb.Vessels[icao].surfaceLatCprOdd))

			// calc longitude zone number and check sanity
			NLeven := common.LongitudeZoneNumber(latEven)
			NLodd := common.LongitudeZoneNumber(latOdd)
			if NLeven != NLodd {
				return errors.New("longitude zone numbers do not match between even and odd cpr values")
			}

			// calculate longitude
			lonEven, lonOdd := common.SurfaceLonGloballyUnambiguous(vdb.refLon, float64(vdb.Vessels[icao].surfaceLonCprEven), float64(vdb.Vessels[icao].surfaceLonCprOdd), NLeven)

			// work out which frame type to return
			switch vdb.Vessels[icao].surfaceLatLonCprTypeHist[len(vdb.Vessels[icao].surfaceLatLonCprTypeHist)-1] {
			case common.CprFormatEvenFrame:
				vdb.Vessels[icao].Lat = latEven
				vdb.Vessels[icao].Lon = lonEven
			case common.CprFormatOddFrame:
				vdb.Vessels[icao].Lat = latOdd
				vdb.Vessels[icao].Lon = lonOdd
			}
			vdb.Vessels[icao].LatLonMethod = "surface,global"

			// Need two lat/lon within ~1NM of each other before saying that we know the position
			previousPosition := haversine.Coord{Lat: oldLat, Lon: oldLon}
			currentPosition := haversine.Coord{Lat: vdb.Vessels[icao].Lat, Lon: vdb.Vessels[icao].Lon}
			mi, _ := haversine.Distance(previousPosition, currentPosition)
			if mi < 1.15 {
				vdb.Vessels[icao].LatLonKnown = true
			}

			return nil
		}
	}
	return errors.New("not enough position messages received")
}

func (vdb *Vessels) updateFromBDS07(icao int, frame bds.BDS07Frame) {
	// BDS 0,7 - Extended squitter status

	switch frame.Ver {
	case 1:
		switch frame.St {
		case 0:
			vdb.setAirborneStatus(icao, true)
		case 1:
			vdb.setAirborneStatus(icao, false)
		}

	case 2:
		switch frame.St {
		case 0:
			vdb.setAirborneStatus(icao, true)
		case 1:
			vdb.setAirborneStatus(icao, false)
		}
	}
}

func (vdb *Vessels) updateFromCommB(icao int, mb []byte, df df.DownlinkFormat, data []byte) {
	// icao = ICAO aircraft address
	// mb = message, Comm-B
	// tc = Type Code
	// data = entire data bytes

	log := log.With().Str("component", "vesselstate").Str("icao", fmt.Sprintf("%06x", icao)).Str("mb", fmt.Sprintf("%07x", mb)).Str("data", fmt.Sprintf("%x", data)).Logger()

	possibleBDS, err := bds.InferBDS(df, mb)
	if err != nil {
		log.Warn().AnErr("err", err).Hex("data", data).Msg("problem inferring BDS code")
		return
	} else {
		log.Debug().Any("bds", possibleBDS).Msg("inferred bds")
	}

	if len(possibleBDS) != 1 {
		if log.Debug().Enabled() {
			log.Warn().Any("bds", possibleBDS).Msg("need at least one BDS messages")
		}
		return
	}

	switch possibleBDS[0] {

	// if message contains BDS05 frame:
	case bds.BDS05:
		bds05frame, err := bds.DecodeBDS05(mb)
		if err != nil {
			log.Err(err).Msg("error decoding BDS05 frame")
		}
		if bds05frame.Tc < 19 {
			vdb.setAltitude(icao, int(math.Round(bds05frame.Altitude)))
		}
		vdb.storeAirborneLatLonCPR(icao, bds05frame.LatCpr, bds05frame.LonCpr, bds05frame.F)

		// see if we can calculate position
		err = vdb.calculateAirbornePosition(icao)
		if err != nil {
			if log.Debug().Enabled() {
				log.Err(err).Msg("could not determine airborne lat/lon")
			}
		}
		return

	// if message contains BDS06 frame:
	case bds.BDS06:
		bds06frame, err := bds.DecodeBDS06(mb)
		if err != nil {
			log.Err(err).Msg("error decoding BDS06 frame")
		}
		vdb.setGroundSpeed(icao, bds06frame.GroundSpeed)
		vdb.setGroundTrack(icao, bds06frame.GroundTrack)
		vdb.storeSurfaceLatLonCPR(icao, bds06frame.LatCpr, bds06frame.LonCpr, bds06frame.F)

		// see if we can calculate position
		err = vdb.calculateSurfacePosition(icao)
		if err != nil {
			if log.Debug().Enabled() {
				log.Err(err).Msg("could not determine surface lat/lon")
			}
		}

		return

	// if message contains BDS07 frame:
	case bds.BDS07:
		bds07frame, err := bds.DecodeBDS07(mb)
		if err != nil {
			log.Err(err).Msg("error decoding BDS06 frame")
		}
		vdb.updateFromBDS07(icao, bds07frame)
		return

	// if message contains BDS08 frame:
	case bds.BDS08:
		bds08frame, err := bds.DecodeBDS08(mb)
		if err != nil {
			log.Err(err).Msg("error decoding BDS06 frame")
		}
		vdb.setCallsign(icao, bds08frame.Callsign)
		return

	// TODO: if message contains BDS09 frame:
	case bds.BDS09:
		return

	// TODO: if message contains BDS10 frame:
	case bds.BDS10:
		return

	case bds.BDS17:
		return

	// if message contains BDS20 frame:
	case bds.BDS20:
		bds20frame, err := bds.DecodeBDS20(mb)
		if err != nil {
			log.Err(err).Msg("error decoding BDS06 frame")
		}
		vdb.setCallsign(icao, bds20frame.Callsign)
		return

	// TODO: if message contains BDS40 frame:
	case bds.BDS40:
		return

	// if message contains BDS50 frame:
	case bds.BDS50:
		bds50frame, err := bds.DecodeBDS50(mb)
		if err != nil {
			log.Err(err).Msg("error decoding BDS50 frame")
		}

		if bds50frame.GroundSpeedValid {
			kts := bds50frame.GroundSpeed
			kph := kts * 0.539957
			vdb.setGroundSpeed(icao, fmt.Sprintf("%.4f km/h (%.4f kts)", kph, kts))
		}

		if bds50frame.TrueTrackAngleValid {
			tta := int(math.Round(bds50frame.TrueTrackAngle))
			vdb.setGroundTrack(icao, fmt.Sprintf("%d°", tta))
		}

		return

	// TODO: if message contains BDS60 frame:
	case bds.BDS60:
		return

	default:

		// if message contains BDS61 frame:
		_, err = bds.DecodeBDS61(mb)
		if err == nil {
			// TODO
			return
		}

		// if message contains BDS61 frame:
		_, err = bds.DecodeBDS62(mb)
		if err == nil {
			// TODO
			return
		}
	}

	log.Warn().Msg("type code not handled")
}

func (vdb *Vessels) UpdateFromDF0(msg df.DF0message) {
	// updates vessel status based on information from DF0 message
	vdb.addVessel(msg.ICAO)
	vdb.incrementMessageCount(msg.ICAO)
	vdb.setAirborneStatus(msg.ICAO, msg.Airborne)
	vdb.setAltitude(msg.ICAO, int(math.Round(msg.Altitude)))
}

func (vdb *Vessels) UpdateFromDF4(msg df.DF4message) {
	// updates vessel status based on information from DF4 message
	vdb.addVessel(msg.ICAO)
	vdb.incrementMessageCount(msg.ICAO)
	vdb.setAirborneStatus(msg.ICAO, msg.Airborne)
	vdb.setAltitude(msg.ICAO, int(math.Round(msg.Altitude)))
}

func (vdb *Vessels) UpdateFromDF5(msg df.DF5message) {
	// updates vessel status based on information from DF5 message
	vdb.addVessel(msg.ICAO)
	vdb.incrementMessageCount(msg.ICAO)
	vdb.setAirborneStatus(msg.ICAO, msg.Airborne)
	vdb.setSquawkCode(msg.ICAO, msg.Squawk)
}

func (vdb *Vessels) UpdateFromDF11(msg df.DF11message) {
	// updates vessel status based on information from DF11 message
	vdb.addVessel(msg.ICAO)
	vdb.incrementMessageCount(msg.ICAO)
}

func (vdb *Vessels) UpdateFromDF16(msg df.DF16message) {
	// updates vessel status based on information from DF16 message
	vdb.addVessel(msg.ICAO)
	vdb.incrementMessageCount(msg.ICAO)
	vdb.setAirborneStatus(msg.ICAO, msg.Airborne)
	vdb.setAltitude(msg.ICAO, int(math.Round(msg.Altitude)))
}

func (vdb *Vessels) UpdateFromDF17(msg df.DF17message, data []byte) {
	// updates vessel status based on information from DF17 message
	vdb.addVessel(msg.ICAO)
	vdb.incrementMessageCount(msg.ICAO)
	vdb.updateFromCommB(msg.ICAO, msg.ME, df.DF17, data)
}

func (vdb *Vessels) UpdateFromDF18(msg df.DF18message, data []byte) {
	// updates vessel status based on information from DF18 message
	vdb.addVessel(msg.ICAO)
	vdb.incrementMessageCount(msg.ICAO)
	vdb.updateFromCommB(msg.ICAO, msg.ME, df.DF18, data)
}

func (vdb *Vessels) UpdateFromDF20(msg df.DF20message, data []byte) {
	// updates vessel status based on information from DF20 message
	vdb.addVessel(msg.ICAO)
	vdb.incrementMessageCount(msg.ICAO)
	vdb.setAltitude(msg.ICAO, int(msg.Altitude))
	vdb.setAirborneStatus(msg.ICAO, msg.Airborne)
	vdb.updateFromCommB(msg.ICAO, msg.MB, df.DF20, data)

	// TODO: Downlink request
	// TODO: Utility message
}

func (vdb *Vessels) UpdateFromDF21(msg df.DF21message, data []byte) {
	// updates vessel status based on information from DF21 message
	vdb.addVessel(msg.ICAO)
	vdb.incrementMessageCount(msg.ICAO)
	vdb.setAirborneStatus(msg.ICAO, msg.Airborne)
	vdb.setSquawkCode(msg.ICAO, msg.Squawk)
	vdb.updateFromCommB(msg.ICAO, msg.MB, df.DF21, data)

	// TODO: Downlink request
	// TODO: Utility message
}
