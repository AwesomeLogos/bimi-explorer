package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var embeddedFiles embed.FS
var staticHandler = initStaticHandler()

func initStaticHandler() http.Handler {

	fsys, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		logger.Error("unable to create static file system", "error", err)
		panic(err)
	}

	return http.FileServer(http.FS(fsys))
}
