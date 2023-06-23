package bds

// func TestDecodeBDS44windSpeedDirection(t *testing.T) {
// 	data := []byte{0xA0, 0x00, 0x16, 0x92, 0x18, 0x5B, 0xD5, 0xCF, 0x40, 0x00, 0x00, 0xDF, 0xC6, 0x96}
// 	df := getDF(data)
// 	assert.Equal(t, DF20, df)
// 	msg := decodeDF20(data)
// 	ws, wd, err := decodeBDS44windSpeedDirection(msg.mb)
// 	assert.NoError(t, err)
// 	assert.Equal(t, float64(22), ws)
// 	assert.Equal(t, float64(344.53125), wd)
// }

// func TestDecodeBDS44staticAirTemperature(t *testing.T) {
// 	data := []byte{0xA0, 0x00, 0x16, 0x92, 0x18, 0x5B, 0xD5, 0xCF, 0x40, 0x00, 0x00, 0xDF, 0xC6, 0x96}
// 	df := getDF(data)
// 	assert.Equal(t, DF20, df)
// 	msg := decodeDF20(data)
// 	sat, err := decodeBDS44staticAirTemperature(msg.mb)
// 	assert.NoError(t, err)
// 	assert.Equal(t, -48.75, sat)
// }
