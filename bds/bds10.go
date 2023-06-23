package bds

import (
	"errors"
)

type BDS10Frame struct {
}

func DecodeBDS10(mb []byte) (frame BDS10Frame, err error) {

	// check the bds code matches
	if int(mb[0]) != 0b00010000 {
		err = errors.New("bds code mismatch")
		return
	}

	// check reserved bits
	if int(mb[1])&0b01111100 != 0 {
		err = errors.New("reserved bits not zero")
		return
	}

	// There's lots of other info that we probably don't care about
	// See "Technical Provisions for Mode S Services and Extended Squitter", Doc 9871
	// Table A-2-16. BDS code 1,0 — Data link capability report

	return
}
