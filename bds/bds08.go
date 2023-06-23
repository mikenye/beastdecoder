package bds

// Aircraft identification and category
// Type Code (TC): 1 to 4
// https://mode-s.org/decode/content/ads-b/2-identification.html

import (
	"errors"
	"fmt"
	"strings"
)

type BDS08Frame struct {
	// https://mode-s.org/decode/content/ads-b/2-identification.html#aircraft-identification-and-category
	tc         int    // Type Code
	ca         int    // Aircraft category
	capability string // Aircraft category string
	Callsign   string // Callsign
}

func DecodeBDS08(mb []byte) (frame BDS08Frame, err error) {
	// decode aircraft identification bytes into struct
	// https://mode-s.org/decode/content/ads-b/2-identification.html#aircraft-identification-and-category

	frame = BDS08Frame{}

	frame.tc = (int(mb[0]) & 0b11111000) >> 3

	if frame.tc < 1 || frame.tc > 5 {
		err = errors.New("type code not from 1 to 4")
		return
	}

	frame.ca = (int(mb[0]) & 0b00000111)

	callsignIndexes := make([]int, 8)
	callsignIndexes[0] = (int(mb[1]) & 0b11111100) >> 2
	callsignIndexes[1] = ((int(mb[1]) & 0b00000011) << 4) + ((int(mb[2]) & 0b11110000) >> 4)
	callsignIndexes[2] = ((int(mb[2]) & 0b00001111) << 2) + ((int(mb[3]) & 0b11000000) >> 6)
	callsignIndexes[3] = (int(mb[3]) & 0b00111111)
	callsignIndexes[4] = (int(mb[4]) & 0b11111100) >> 2
	callsignIndexes[5] = ((int(mb[4]) & 0b00000011) << 4) + ((int(mb[5]) & 0b11110000) >> 4)
	callsignIndexes[6] = ((int(mb[5]) & 0b00001111) << 2) + ((int(mb[6]) & 0b11000000) >> 6)
	callsignIndexes[7] = (int(mb[6]) & 0b00111111)

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

	// determine category code
	switch frame.tc {
	case 4:
		frame.capability = fmt.Sprintf("A%d", frame.ca)
	case 3:
		frame.capability = fmt.Sprintf("B%d", frame.ca)
	case 2:
		frame.capability = fmt.Sprintf("C%d", frame.ca)
	case 1:
		frame.capability = fmt.Sprintf("D%d", frame.ca)
	}

	// sanity check
	if !validCallsign.Match([]byte(frame.Callsign)) {
		err = errors.New("callsign contains invalid characters")
	}

	return

}
