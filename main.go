package main

import (
	"beastdecoder/vesselstate"
	"beastdecoder/webview"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"
)

var vdb vesselstate.Vessels

func main() {
	app := &cli.App{
		Name:  "beastdecoder",
		Usage: "Decodes BEAST data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Category:  "Web Interface",
				Name:      "webview",
				Usage:     "ip:port to listen on for HTTP connections",
				TakesFile: false,
			},
			&cli.StringSliceFlag{
				Category:  "BEAST Data Input",
				Name:      "connect",
				Usage:     "ip:port of host to receive BEAST data from",
				TakesFile: false,
				KeepSpace: false,
			},
			&cli.Float64Flag{
				Category: "Receiver Location",
				Name:     "lat",
				Usage:    "latitude of receiver",
			},
			&cli.Float64Flag{
				Category: "Receiver Location",
				Name:     "lon",
				Usage:    "longitude of receiver",
			},
			&cli.BoolFlag{
				Category: "Logging",
				Name:     "debug",
				Usage:    "enable debug logging",
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Msg("finished with error")
	} else {
		log.Info().Msg("finished")
	}
}

func run(ctx *cli.Context) error {

	// init logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.UnixDate})
	if !ctx.Bool("debug") {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Info().Msg(fmt.Sprintf("starting %s", ctx.App.Name))

	// init vessel database
	vdb.Init()

	// set refLat/refLon if given
	if ctx.IsSet("lat") && ctx.IsSet("lon") {
		vdb.SetRefLatLon(ctx.Float64("lat"), ctx.Float64("lon"))
	}

	// enable web interface
	if ctx.IsSet("webview") {
		addrSplit := strings.Split(ctx.String("webview"), ":")
		if len(addrSplit) != 2 {
			log.Error().Str("addr", ctx.String("webview")).Msg("connect format must be host:port")
			return errors.New("webview format must be host:port")
		}
		ip := net.ParseIP(addrSplit[0])
		port, err := strconv.ParseInt(addrSplit[1], 10, 0)
		if err != nil {
			log.Err(err).Any("port", addrSplit[1]).Msg("could not parse port")
			return errors.New("could not parse port")
		}
		addr := &net.TCPAddr{
			IP:   ip,
			Port: int(port),
		}
		go webview.Init(addr, &vdb)
	}

	// outgoing connections
	for _, addr := range ctx.StringSlice("connect") {
		addrSplit := strings.Split(addr, ":")
		if len(addrSplit) != 2 {
			log.Error().Str("addr", addr).Msg("connect format must be host:port")
			continue
		}
		ip := net.ParseIP(addrSplit[0])
		port, err := strconv.ParseInt(addrSplit[1], 10, 0)
		if err != nil {
			log.Err(err).Any("port", addrSplit[1]).Msg("could not parse port")
		}
		addr := &net.TCPAddr{
			IP:   ip,
			Port: int(port),
		}
		beastDial(addr)
	}
	return nil
}
