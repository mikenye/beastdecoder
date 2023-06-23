package bds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBDS10(t *testing.T) {

	// define test data
	var testTable = []struct {
		data          []byte
		expectedError bool
	}{
		{
			data:          []byte{0x10, 0x03, 0x0a, 0x80, 0xf5, 0x00, 0x00},
			expectedError: false,
		},
	}

	assert := assert.New(t)
	for _, testData := range testTable {
		_, err := DecodeBDS10(testData.data)
		testMsg := fmt.Sprintf("data: %014x, ", testData.data)
		if !testData.expectedError {
			assert.NoError(err, testMsg+"DecodeBDS10 error")
		} else {
			assert.Error(err, testMsg+"DecodeBDS10 no error")
		}
	}

}
