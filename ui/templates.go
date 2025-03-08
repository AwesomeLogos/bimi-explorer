package ui

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/AwesomeLogos/bimi-explorer/internal/common"
)

//go:embed partials
var partialsFS embed.FS

//go:embed all:views
var viewsFS embed.FS

type TemplateFunc func(data any) (string, error)

type TemplateData map[string]any

var templateCache = initTemplates()

func initTemplates() map[string]TemplateFunc {

	funcMap := template.FuncMap{
		"dec": func(i int) int {
			return i - 1
		},
		"inc": func(i int) int {
			return i + 1
		},
		"loop": func(from, to int) []int {
			result := []int{}
			for i := from; i < to; i++ {
				result = append(result, i)
			}
			return result
		},
	}

	theCache := make(map[string]TemplateFunc)

	var partials bytes.Buffer

	partialErr := fs.WalkDir(partialsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// read the partials file and append to partials
			content, readErr := fs.ReadFile(partialsFS, path)
			if readErr != nil {
				common.Logger.Error("unable to read partials file", "err", readErr, "filename", path)
				return err
			}
			partials.Write(content)
		}
		return nil
	})
	if partialErr != nil {
		common.Logger.Error("unable to register partials", "err", partialErr)
	}

	viewErr := fs.WalkDir(viewsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			common.Logger.Error("walkdir error", "err", err)
			return err
		}
		if !d.IsDir() {
			common.Logger.Debug("registering view", "filename", path)
			content, readErr := fs.ReadFile(viewsFS, path)
			if readErr != nil {
				common.Logger.Error("unable to read view file", "err", readErr, "filename", path)
				return err
			}
			name := path[len("views/"):]

			var templateBuffer bytes.Buffer
			templateBuffer.Write(content)
			templateBuffer.Write(partials.Bytes())
			t := template.New(name).Funcs(funcMap)
			template, parseErr := t.Parse(templateBuffer.String())
			if parseErr != nil {
				common.Logger.Error("unable to parse template", "err", parseErr, "filename", path, "content", string(content))
				return parseErr
			}
			theCache[name] = func(data any) (string, error) {
				var buf bytes.Buffer
				err := template.Execute(&buf, data)
				if err != nil {
					common.Logger.Error("unable to execute template", "err", err, "filename", path, "content", string(content))
					return "", err
				}
				return buf.String(), nil
			}
		}
		return nil
	})
	if viewErr != nil {
		common.Logger.Error("unable to register views", "err", viewErr)
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

func RunTemplate(w http.ResponseWriter, r *http.Request, templateName string, data TemplateData) {

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
		common.Logger.Error("template failed", "err", execErr, "template", templateName)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(result))
}
