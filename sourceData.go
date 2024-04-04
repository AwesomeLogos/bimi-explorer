package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"net/http"

	"github.com/AwesomeLogos/bimi-explorer/logosearch"
)

func buildSourceData() ([]byte, error) {
	count, _ := countDomains()

	domains, dbErr := listDomains(int32(count), 0)
	if dbErr != nil {
		return nil, dbErr
	}

	searchIndex := logosearch.GenerateIndex(domains)

	jsonData, jsonErr := json.MarshalIndent(searchIndex, "", "  ")
	if jsonErr != nil {
		return nil, jsonErr
	}
	return jsonData, nil
}

func sourceDataJson(w http.ResponseWriter, r *http.Request) {

	sourceData, sourceErr := buildSourceData()
	if sourceErr != nil {
		http.Error(w, "Unable to generate sourceData", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(sourceData)
}

func sourceDataTgz(w http.ResponseWriter, r *http.Request) {

	sourceData, sourceErr := buildSourceData()
	if sourceErr != nil {
		http.Error(w, "Unable to generate sourceData", http.StatusInternalServerError)
		return
	}

	hdr := &tar.Header{
		Name: "sourceData.json",
		Mode: 0600,
		Size: int64(len(sourceData)),
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=sourceData.tgz")
	gw := gzip.NewWriter(w)

	tw := tar.NewWriter(gw)
	tw.WriteHeader(hdr)
	tw.Write(sourceData)
	tw.Close()
	gw.Close()

}
