package webview

import (
	"beastdecoder/vesselstate"
	_ "embed"
	"fmt"
	"html/template"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

//go:embed webview.gtpl
var webviewTemplate string

func httpRenderWebview(w http.ResponseWriter, r *http.Request, vdb *vesselstate.Vessels) {
	log := log.With().Str("component", "webview").Logger()

	// Make and parse the HTML template
	t, err := template.New("stats").Parse(webviewTemplate)
	if err != nil {
		log.Err(err).Str("func", "httpRenderWebview").Str("reqURI", r.RequestURI).Msg("could not render webviewTemplate")
	}

	vdb.RLock()
	defer vdb.RUnlock()

	err = t.Execute(w, vdb.Vessels)
	if err != nil {
		fmt.Println(err)
		log.Panic().AnErr("err", err).Str("func", "httpRenderWebview").Str("reqURI", r.RequestURI).Msg("could not execute webviewTemplate")
	}
}

func Init(addr net.Addr, vdb *vesselstate.Vessels) {
	log := log.With().Str("component", "webview").Logger()

	// stats http server routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpRenderWebview(w, r, vdb)
	})

	// start stats http server
	log.Info().Str("ip", "0.0.0.0").Int("port", 6969).Msg("starting webview listener")
	err := http.ListenAndServe(addr.String(), nil)
	if err != nil {
		log.Panic().AnErr("err", err).Msg("stats server stopped")
	}

}
