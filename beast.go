package main

import (
	"beastdecoder/df"
	"errors"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type syncState uint8

const syncStateUnknown = syncState(0)
const syncStateNot1a = syncState(1)
const syncState1a = syncState(2)
const syncStateFrameType = syncState(3)

func beastDial(addr *net.TCPAddr) {

	// set up logger
	log := log.With().Str("src", addr.String()).Logger()

	var frameType byte
	buf := make([]byte, 2)

	for {

		log.Info().Msg("connecting")
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			log.Info().Err(err).Msg("connection error")
			continue
		}
		defer conn.Close()
		log.Info().Msg("connected, synchronising")

		synced := false
		for {
			if !synced {
				frameType, err = beastSync(conn)
				if err != nil {
					log.Err(err).Msg("synchronisation error")
					break
				} else {
					synced = true
					log.Info().Msg("synchronised, receiving")
				}

			} else {
				_, err := conn.Read(buf)
				if err != nil {
					log.Err(err).Msg("receive error")
					break
				}
				if buf[0] == 0x1a {
					frameType = buf[1]
				} else {
					synced = false
				}
			}

			if synced {

				frameData, err := beastReadFrame(conn, frameType)
				if err != nil {
					log.Err(err).Msg("frame read error")
					break
				}

				var data []byte
				var DF df.DownlinkFormat

				switch frameType {

				// Mode-AC frame
				case 0x31:
					// ignore this frame type
					continue

				// Mode-S short frame
				case 0x32:
					_, _, data = handleFrameData(frameType, frameData)
					DF = df.GetDF(data)

				// Mode-S long frame
				case 0x33:
					_, _, data = handleFrameData(frameType, frameData)
					DF = df.GetDF(data)

				default:
					log.Warn().Hex("frameType", []byte{frameType}).Hex("frameData", frameData).Msg("unknown frameType")
				}

				log.Debug().Msg("START OF FRAME")

				log.Debug().Uint8("DF", uint8(DF)).Hex("data", data).Msg("received")

				switch DF {
				case df.DF0:
					msg, err := df.DecodeDF0(data)
					if err != nil {
						log.Err(err).Hex("data", data).Msg("error decoding DF0")
					} else {
						vdb.UpdateFromDF0(msg)
					}

				case df.DF4:
					msg, err := df.DecodeDF4(data)
					if err != nil {
						log.Err(err).Hex("data", data).Msg("error decoding DF4")
					} else {
						vdb.UpdateFromDF4(msg)
					}

				case df.DF5:
					msg, err := df.DecodeDF5(data)
					if err != nil {
						log.Err(err).Hex("data", data).Msg("error decoding DF5")
					} else {
						vdb.UpdateFromDF5(msg)
					}

				case df.DF11:
					msg := df.DecodeDF11(data)
					vdb.UpdateFromDF11(msg)

				case df.DF16:
					msg, err := df.DecodeDF16(data)
					if err != nil {
						log.Err(err).Hex("data", data).Msg("error decoding DF16")
					} else {
						vdb.UpdateFromDF16(msg)
					}

				case df.DF17:
					msg := df.DecodeDF17(data)
					vdb.UpdateFromDF17(msg, data)

				case df.DF18:
					msg := df.DecodeDF18(data)
					vdb.UpdateFromDF18(msg, data)

				case df.DF19:
					// military stuff, can't decode
					break

				case df.DF20:
					msg, err := df.DecodeDF20(data)
					if err != nil {
						log.Err(err).Hex("data", data).Msg("error decoding DF20")
					} else {
						vdb.UpdateFromDF20(msg, data)
					}

				case df.DF21:
					msg, err := df.DecodeDF21(data)
					if err != nil {
						log.Err(err).Hex("data", data).Msg("error decoding DF21")
					} else {
						vdb.UpdateFromDF21(msg, data)
					}

				case df.DF24:
					// extended length messages
					// TODO: decode these
					break

				default:
					log.Error().Hex("data", data).Msg("unsupported data")
					os.Exit(1)
				}

				log.Debug().Msg("END OF FRAME")
			}
		}
		time.Sleep(time.Second * 30)
	}
}

func beastSync(conn *net.TCPConn) (frameType byte, err error) {

	state := syncStateUnknown

	discardedByteCount := 0
	// discardedBytes := make([]byte, 1024)

	buf := make([]byte, 1)

	// synchronisation
	for state != syncStateFrameType {

		// read a byte from the network
		_, err := conn.Read(buf)
		if err != nil {
			return frameType, err
		}

		// determine sync status
		switch {
		case buf[0] != 0x1a && state == syncStateUnknown:
			state = syncStateNot1a
		case buf[0] == 0x1a && state == syncStateNot1a:
			state = syncState1a
		case buf[0] == 0x31 && state == syncState1a:
			state = syncStateFrameType
			frameType = 0x31
		case buf[0] == 0x32 && state == syncState1a:
			state = syncStateFrameType
			frameType = 0x32
		case buf[0] == 0x33 && state == syncState1a:
			state = syncStateFrameType
			frameType = 0x33
		case buf[0] == 0x34 && state == syncState1a:
			state = syncStateFrameType
			frameType = 0x34
		default:
			// discardedBytes[discardedByteCount] = buf[0]
			discardedByteCount++
			state = syncStateUnknown
		}
	}
	if discardedByteCount > 0 {
		// .Hex("discardedBytes", discardedBytes[:discardedByteCount])
		log.Debug().Int("discardedByteCount", discardedByteCount).Msg("discarded data due to synchronisation")
	}
	return frameType, err
}

func beastReadFrame(conn *net.TCPConn, frameType byte) (frameData []byte, err error) {

	buf := make([]byte, 1)
	bytesToRead := 0

	switch frameType {

	// Mode-AC frame
	//  - 6 byte MLAT timestamp
	//  - 1 byte RSSI
	//  - 2 byte Mode-AC data
	case 0x31:
		bytesToRead = 9

	// Mode-S short frame
	//  - 6 byte MLAT timestamp
	//  - 1 byte RSSI
	//  - 7 byte Mode-S short data
	case 0x32:
		bytesToRead = 14

	// Mode-S long frame
	//  - 6 byte MLAT timestamp
	//  - 1 byte RSSI
	//  - 14 byte Mode-S long data
	case 0x33:
		bytesToRead = 21

	case 0x34:
		log.Warn().Msg("frame type 0x34 not implemented")
	}

	frameData = make([]byte, bytesToRead)
	n := 0

	for n < bytesToRead {
		_, err := conn.Read(buf)
		if err != nil {
			return frameData, err
		}

		if buf[0] == 0x1a {
			_, err := conn.Read(buf)
			if err != nil {
				return frameData, err
			}
			if buf[0] == 0x1a {
				frameData[n] = 0x1a
			} else {
				return frameData, errors.New("escape issue")
			}
		} else {
			frameData[n] = buf[0]
		}

		n++
	}
	return frameData, err
}

func handleFrameData(frameType byte, frameData []byte) (mlatTimestamp []byte, rssi byte, data []byte) {

	mlatTimestamp = frameData[0:5]
	rssi = frameData[6]
	data = frameData[7:]

	return mlatTimestamp, rssi, data

}
