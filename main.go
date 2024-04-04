package main

import (
	"net/http"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) > 1 {
		if os.Args[1] == "--validate" {
			bulkValidate()
		} else {
			bulkLoader(os.Args[1:])
		}
		return
	}

	var listenPort, portErr = strconv.Atoi(os.Getenv("PORT"))
	if portErr != nil {
		listenPort = 4000
	}
	var listenAddress = os.Getenv("ADDRESS")

	http.HandleFunc("/status.json", statusHandler)
	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/bimi/", http.StatusSeeOther) })
	http.HandleFunc("GET /bimi/{$}", rootHandlerGet)
	http.HandleFunc("POST /bimi/{$}", rootHandlerPost)
	http.HandleFunc("/robots.txt", staticHandler.ServeHTTP)
	http.HandleFunc("/favicon.ico", staticHandler.ServeHTTP)
	http.HandleFunc("/favicon.svg", staticHandler.ServeHTTP)
	http.HandleFunc("/images/", staticHandler.ServeHTTP)
	http.HandleFunc("/bimi/{domain}/{$}", bimiHandler)
	http.HandleFunc("/bimi/invalid.html", listInvalidHandler)
	http.HandleFunc("/bimi/list.html", listHandler)
	http.HandleFunc("/bimi/view.html", viewHandler)
	http.HandleFunc("/bimi/sourceData.json", sourceDataJson)
	http.HandleFunc("/bimi/sourceData.tgz", sourceDataTgz)

	err := http.ListenAndServe(listenAddress+":"+strconv.Itoa(listenPort), nil)
	if err != nil {
		logger.Error("unable to listen", "address", listenAddress, "port", listenPort, "error", err)
	}
}
