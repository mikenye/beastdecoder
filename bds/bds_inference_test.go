package bds

// func TestNotInferrable(t *testing.T) {
// 	// a921109446da704cd0690dffe93e
// 	// a800091000800081c081f052a261
// 	// a8001e9800800081c081f0b0ed63
// 	// msg := []byte{0xa9, 0x21, 0x10, 0x94, 0x46, 0xda, 0x70, 0x4c, 0xd0, 0x69, 0x0d, 0xff, 0xe9, 0x3e}
// 	// msg := []byte{0xa8, 0x00, 0x09, 0x10, 0x00, 0x80, 0x00, 0x81, 0xc0, 0x81, 0xf0, 0x52, 0xa2, 0x61}
// 	data := []byte{0xa8, 0x00, 0x1e, 0x98, 0x00, 0x80, 0x00, 0x81, 0xc0, 0x81, 0xf0, 0xb0, 0xed, 0x63}
// 	df := getDF(data)
// 	assert.Equal(t, uint8(21), df)
// 	msg := decodeDF21(data)
// 	inferBDS(msg.mb)
// 	t.Log(fmt.Sprintf("%08b", msg.mb))
// 	frame, e := decodeBDS06(msg.mb)
// 	assert.NoError(t, e)
// 	t.Log(frame)

// 	refLat := -31.897
// 	refLon := 115.928

// 	lat := calcSurfaceLatLocallyUnambiguous(refLat, float64(frame.latCpr), frame.f)
// 	NL := calcLongitudeZoneNumber(lat)
// 	lon := calcSurfaceLonLocallyUnambiguous(refLon, float64(frame.lonCpr), NL, frame.f)
// 	t.Log(lat, lon)

// 	valid := isSurfacePosValid(refLat, refLon, lat, lon)
// 	t.Log(valid)

// }

// func TestInferBDS60(t *testing.T) {

// 	data := []byte{0xA8, 0x00, 0x1E, 0xBC, 0xFF, 0xFB, 0x23, 0x28, 0x60, 0x04, 0xA7, 0x3F, 0x6A, 0x5B}
// 	df := getDF(data)
// 	assert.Equal(t, uint8(21), df)
// 	msg := decodeDF20(data)
// 	bc, err := inferBDS(msg.mb)
// 	assert.NoError(t, err)
// 	assert.Equal(t, []BDScode{BDS60}, bc)

// }
