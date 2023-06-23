package bds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS62(t *testing.T) {
	var testTable = []struct {
		data  []byte
		vdaSi int // vertical data available / source indicator
		tat   int // target altitude type
		bcf   int // backwards capability flag
		tac   int // target altitude compatibility
		vmi   int // vertical mode indicator
		ta    int // target altitude

		hdaSi int // horizontal data available / source indicator
		thta  int // target heading / track angle
		hmi   int // horizontal mode indicator

		nacp    int // navigation accuracy category
		nicbaro int // navigation integrity category - baro
		sil     int // surveillance integrity level

		cmc int // capability / mode codes
		eps int // emergency / priority status
	}{
		{
			data: []byte{0xe9, 0x08, 0x32, 0x00, 0x09, 0x38, 0x10},
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {

		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		frame, err := DecodeBDS62(testData.data)
		assert.NoError(err, testMsg+"decodeBDS62 error")

		fmt.Println("")
		fmt.Println(frame.Sprint())
	}

}
