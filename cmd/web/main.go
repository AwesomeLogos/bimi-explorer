package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/AwesomeLogos/bimi-explorer/internal/common"
	"github.com/AwesomeLogos/bimi-explorer/internal/server"
	"github.com/AwesomeLogos/bimi-explorer/ui"
)

func main() {

	var listenPort, portErr = strconv.Atoi(os.Getenv("PORT"))
	if portErr != nil {
		listenPort = 4000
	}
	var listenAddress = os.Getenv("ADDRESS")

	http.HandleFunc("/status.json", server.StatusHandler)
	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/bimi/", http.StatusSeeOther) })
	http.HandleFunc("GET /bimi/{$}", server.RootHandlerGet)
	http.HandleFunc("POST /bimi/{$}", server.RootHandlerPost)
	http.HandleFunc("/robots.txt", ui.StaticHandler.ServeHTTP)
	http.HandleFunc("/favicon.ico", ui.StaticHandler.ServeHTTP)
	http.HandleFunc("/favicon.svg", ui.StaticHandler.ServeHTTP)
	http.HandleFunc("/images/", ui.StaticHandler.ServeHTTP)
	http.HandleFunc("/bimi/{domain}/{$}", server.BimiHandler)
	http.HandleFunc("/bimi/{domain}/refresh.html", server.RefreshHandler)
	http.HandleFunc("/bimi/invalid.html", server.ListInvalidHandler)
	http.HandleFunc("/bimi/list.html", server.ListHandler)
	http.HandleFunc("/bimi/view.html", server.ViewHandler)
	http.HandleFunc("/bimi/sourceData.json", server.SourceDataJson)
	http.HandleFunc("/bimi/sourceData.tgz", server.SourceDataTgz)

	err := http.ListenAndServe(listenAddress+":"+strconv.Itoa(listenPort), nil)
	if err != nil {
		common.Logger.Error("unable to listen", "address", listenAddress, "port", listenPort, "error", err)
	}
}
