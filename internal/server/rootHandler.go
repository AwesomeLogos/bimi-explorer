package server

import (
	"net/http"
	"strconv"

	"github.com/AwesomeLogos/bimi-explorer/internal/db"
	"github.com/AwesomeLogos/bimi-explorer/internal/db/generated"
	"github.com/AwesomeLogos/bimi-explorer/ui"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {

	pageSize, psErr := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if psErr != nil || pageSize < 1 {
		pageSize = 2500
	}

	count, _ := db.CountDomains()
	maxPages := int(count / int64(pageSize))
	if count%int64(pageSize) > 0 {
		maxPages++
	}

	currentPage, cpErr := strconv.Atoi(r.URL.Query().Get("page"))
	if cpErr != nil || currentPage < 1 || currentPage > maxPages {
		currentPage = 1
	}
	domains, dbErr := db.ListDomains(int32(pageSize), int32((currentPage-1)*pageSize))

	ui.RunTemplate(w, r, "bimi/list.tmpl", map[string]any{
		"Count":       count,
		"CurrentPage": currentPage,
		"Domains":     domains,
		"Err":         dbErr,
		"MaxPage":     maxPages,
		"PageSize":    pageSize,
		"Title":       "List Domains",
	})

}

func ListInvalidHandler(w http.ResponseWriter, r *http.Request) {

	pageSize, psErr := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if psErr != nil || pageSize < 1 {
		pageSize = 100
	}

	count, _ := db.CountInvalidDomains()
	maxPages := int(count / int64(pageSize))
	if count%int64(pageSize) > 0 {
		maxPages++
	}

	currentPage, cpErr := strconv.Atoi(r.URL.Query().Get("page"))
	if cpErr != nil || currentPage < 1 || currentPage > maxPages {
		currentPage = 1
	}
	domains, dbErr := db.ListInvalidDomains(int32(pageSize), int32((currentPage-1)*pageSize))

	ui.RunTemplate(w, r, "bimi/invalid.tmpl", map[string]any{
		"Count":       count,
		"CurrentPage": currentPage,
		"Domains":     domains,
		"Err":         dbErr,
		"MaxPage":     maxPages,
		"PageSize":    pageSize,
		"Title":       "Invalid BIMI Logos",
	})

}

func RootHandlerGet(w http.ResponseWriter, r *http.Request) {

	var domains []generated.Domain
	var dbErr error

	query := r.URL.Query().Get("q")
	if query != "" {
		domains, dbErr = db.SearchDomains(query)
	} else {
		domains, dbErr = db.ListRandomDomains(50)
	}

	ui.RunTemplate(w, r, "bimi/index.tmpl", map[string]any{
		"Domains": domains,
		"Err":     dbErr,
		"H1":      "Welcome",
		"Query":   query,
		"Title":   "BIMI Explorer",
	})

}

func RootHandlerPost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/bimi/", http.StatusSeeOther)
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {

	pageSize, psErr := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if psErr != nil || pageSize < 1 {
		pageSize = 144
	}

	count, _ := db.CountDomains()
	maxPages := int(count / int64(pageSize))
	if count%int64(pageSize) > 0 {
		maxPages++
	}

	currentPage, cpErr := strconv.Atoi(r.URL.Query().Get("page"))
	if cpErr != nil || currentPage < 1 || currentPage > maxPages {
		currentPage = 1
	}

	domains, dbErr := db.ListDomains(int32(pageSize), int32((currentPage-1)*pageSize))

	ui.RunTemplate(w, r, "bimi/view.tmpl", map[string]any{
		"CurrentPage": currentPage,
		"Domains":     domains,
		"Err":         dbErr,
		"MaxPage":     maxPages,
		"PageSize":    pageSize,
		"Title":       "View BIMI Logos",
	})

}
