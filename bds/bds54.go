package bds

import (
	"errors"
	"strings"
)

type BDS54Frame struct {
	Waypoint string  // waypoint name
	ETA      float64 // Estimated Time of Arrival (normal flight)
	EFL      int     // Estimated Flight Level (normal flight)
	TTG      float64 // Time to Go (direct route)
}

func DecodeBDS54(mb []byte) (frame BDS54Frame, err error) {

	// check status bit
	if (int(mb[0])&0b10000000)>>7 == 1 {
		err = errors.New("status bit denotes invalid parameters")
		return
	}

	// check reserved bit
	if (int(mb[6]) & 0b00000001) != 0 {
		err = errors.New("reserved bit not zero")
		return
	}

	// waypoint characters
	waypointIndexes := make([]int, 5)
	waypointIndexes[0] = ((int(mb[0]) & 0b01111110) >> 1)
	waypointIndexes[1] = ((int(mb[0]) & 0b00000001) << 5) + ((int(mb[1]) & 0b11111000) >> 3)
	waypointIndexes[2] = ((int(mb[1]) & 0b00000111) << 3) + ((int(mb[2]) & 0b11100000) >> 5)
	waypointIndexes[3] = ((int(mb[2]) & 0b00011111) << 1) + ((int(mb[3]) & 0b10000000) >> 7)
	waypointIndexes[4] = ((int(mb[3]) & 0b01111110) >> 1)

	// sanity check & build waypoint
	for i := range waypointIndexes {
		if waypointIndexes[i] >= len(waypointCharMap) {
			err = errors.New("waypoint character out of range")
			return
		} else {
			frame.Waypoint += waypointCharMap[waypointIndexes[i]]
		}
	}

	// trim whitespace
	frame.Waypoint = strings.TrimSpace(frame.Waypoint)

	// check waypoint
	if !validWaypoint.Match([]byte(frame.Waypoint)) {
		err = errors.New("waypoint contains invalid characters")
	}

	// estimated time of arrival
	frame.ETA = float64(((int(mb[3])&0b00000001)<<8)+int(mb[4])) * (60.0 / 512.0)
	if frame.ETA < 0 || frame.ETA > 60 {
		err = errors.New("estimated time of arrival out of range (>60 mins)")
		return
	}

	// estimated flight level (ft)
	frame.EFL = ((int(mb[5]) & 0b11111100) >> 2) * 10

	// time to go
	frame.TTG = float64(((int(mb[5])&0b00000011)<<7)+((int(mb[6])&0b11111110)>>1)) * (60.0 / 512.0)

	return
}
