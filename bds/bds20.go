package bds

// All-call reply
// https://mode-s.org/decode/content/mode-s/2-allcall.html

import (
	"errors"
	"strings"
)

type BDS20Frame struct {
	Callsign string
}

func DecodeBDS20(mb []byte) (frame BDS20Frame, err error) {
	// decode callsign from BDS 2,0 message
	// https://mode-s.org/decode/content/mode-s/6-els.html#aircraft-identification-bds-20

	frame = BDS20Frame{}

	callsignIndexes := make([]int, 8)

	// check bds code is 0b00100000
	if mb[0] != 0b00100000 {
		err = errors.New("BDS code not [0010 0000]")
		return
	} else {
		callsignIndexes[0] = (int(mb[1]) & 0b11111100) >> 2
		callsignIndexes[1] = ((int(mb[1]) & 0b00000011) << 4) + ((int(mb[2]) & 0b11110000) >> 4)
		callsignIndexes[2] = ((int(mb[2]) & 0b00001111) << 2) + ((int(mb[3]) & 0b11000000) >> 6)
		callsignIndexes[3] = (int(mb[3]) & 0b00111111)
		callsignIndexes[4] = (int(mb[4]) & 0b11111100) >> 2
		callsignIndexes[5] = ((int(mb[4]) & 0b00000011) << 4) + ((int(mb[5]) & 0b11110000) >> 4)
		callsignIndexes[6] = ((int(mb[5]) & 0b00001111) << 2) + ((int(mb[6]) & 0b11000000) >> 6)
		callsignIndexes[7] = (int(mb[6]) & 0b00111111)
	}

	// sanity check & build callsign
	for i := range callsignIndexes {
		if callsignIndexes[i] >= len(callsignCharMap) {
			err = errors.New("callsign character out of range")
			return
		} else {
			frame.Callsign += callsignCharMap[callsignIndexes[i]]
		}
	}

	// trim whitespace
	frame.Callsign = strings.TrimSpace(frame.Callsign)

	// sanity check
	if !validCallsign.Match([]byte(frame.Callsign)) {
		err = errors.New("callsign contains invalid characters")
	}

	return
}
