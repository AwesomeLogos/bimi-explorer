package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/aymerick/raymond"
)

//go:embed partials
var partialFiles embed.FS

//go:embed all:views
var viewFiles embed.FS

type TemplateFunc func(data any) (string, error)

type TemplateData map[string]any

var templateCache = initTemplates()

func initTemplates() map[string]TemplateFunc {

	theCache := make(map[string]TemplateFunc)

	partialErr := fs.WalkDir(partialFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			logger.Info("registering partial", "filename", path)
			content, readErr := fs.ReadFile(partialFiles, path)
			if readErr != nil {
				logger.Error("unable to read partials file", "err", readErr, "filename", path)
				return err
			}
			name := path[len("partials/") : len(path)-len(".hbs")]
			raymond.RegisterPartial(name, string(content))
		}
		return nil
	})
	if partialErr != nil {
		logger.Error("unable to register partials", "err", partialErr)
	}

	viewErr := fs.WalkDir(viewFiles, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("%v\n", d)
		if err != nil {
			return err
		}
		if !d.IsDir() {
			logger.Info("registering view", "filename", path)
			content, readErr := fs.ReadFile(viewFiles, path)
			if readErr != nil {
				logger.Error("unable to read view file", "err", readErr, "filename", path)
				return err
			}
			name := path[len("views/"):]
			template, parseErr := raymond.Parse(string(content))
			if parseErr != nil {
				logger.Error("unable to parse template", "err", parseErr, "filename", path, "content", string(content))
				return parseErr
			}
			theCache[name] = template.Exec
		}
		return nil
	})
	if viewErr != nil {
		logger.Error("unable to register views", "err", viewErr)
	}

	return theCache
}

type CrumbtrailEntry struct {
	Text string
	URL  string
}

func makeCrumbtrail(r *http.Request) []CrumbtrailEntry {
	crumbtrail := []CrumbtrailEntry{}

	// Add additional entries based on the request path
	path := r.URL.Path
	segments := strings.Split(path, "/")
	for i := 1; i < len(segments); i++ {
		entry := CrumbtrailEntry{
			Text: segments[i],
			URL:  strings.Join(segments[0:i+1], "/"),
		}
		crumbtrail = append(crumbtrail, entry)
	}

	return crumbtrail
}

func runTemplate(w http.ResponseWriter, r *http.Request, templateName string, data TemplateData) {

	fn := templateCache[templateName]
	if fn == nil {
		http.NotFound(w, r)
		return
	}
	if data == nil {
		data = make(map[string]any)
	}
	data["crumbtrail"] = makeCrumbtrail(r)

	result, execErr := fn(data)
	if execErr != nil {
		logger.Error("template failed", "err", execErr, "template", templateName)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(result))
}
