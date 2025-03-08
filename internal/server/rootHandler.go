package main

import (
	"net/http"
	"strconv"

	"github.com/AwesomeLogos/bimi-explorer/generated"
)

func listHandler(w http.ResponseWriter, r *http.Request) {

	pageSize, psErr := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if psErr != nil || pageSize < 1 {
		pageSize = 2500
	}

	count, _ := countDomains()
	maxPages := int(count / int64(pageSize))
	if count%int64(pageSize) > 0 {
		maxPages++
	}

	currentPage, cpErr := strconv.Atoi(r.URL.Query().Get("page"))
	if cpErr != nil || currentPage < 1 || currentPage > maxPages {
		currentPage = 1
	}
	domains, dbErr := listDomains(int32(pageSize), int32((currentPage-1)*pageSize))

	runTemplate(w, r, "bimi/list.tmpl", map[string]any{
		"Count":       count,
		"CurrentPage": currentPage,
		"Domains":     domains,
		"Err":         dbErr,
		"MaxPage":     maxPages,
		"PageSize":    pageSize,
		"Title":       "List Domains",
	})

}

func listInvalidHandler(w http.ResponseWriter, r *http.Request) {

	pageSize, psErr := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if psErr != nil || pageSize < 1 {
		pageSize = 100
	}

	count, _ := countInvalidDomains()
	maxPages := int(count / int64(pageSize))
	if count%int64(pageSize) > 0 {
		maxPages++
	}

	currentPage, cpErr := strconv.Atoi(r.URL.Query().Get("page"))
	if cpErr != nil || currentPage < 1 || currentPage > maxPages {
		currentPage = 1
	}
	domains, dbErr := listInvalidDomains(int32(pageSize), int32((currentPage-1)*pageSize))

	runTemplate(w, r, "bimi/invalid.tmpl", map[string]any{
		"Count":       count,
		"CurrentPage": currentPage,
		"Domains":     domains,
		"Err":         dbErr,
		"MaxPage":     maxPages,
		"PageSize":    pageSize,
		"Title":       "Invalid BIMI Logos",
	})

}

func rootHandlerGet(w http.ResponseWriter, r *http.Request) {

	var domains []generated.Domain
	var dbErr error

	query := r.URL.Query().Get("q")
	if query != "" {
		domains, dbErr = searchDomains(query)
	} else {
		domains, dbErr = listRandomDomains(50)
	}

	runTemplate(w, r, "bimi/index.tmpl", map[string]any{
		"Domains": domains,
		"Err":     dbErr,
		"H1":      "Welcome",
		"Query":   query,
		"Title":   "BIMI Explorer",
	})

}

func rootHandlerPost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/bimi/", http.StatusSeeOther)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	pageSize, psErr := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if psErr != nil || pageSize < 1 {
		pageSize = 144
	}

	count, _ := countDomains()
	maxPages := int(count / int64(pageSize))
	if count%int64(pageSize) > 0 {
		maxPages++
	}

	currentPage, cpErr := strconv.Atoi(r.URL.Query().Get("page"))
	if cpErr != nil || currentPage < 1 || currentPage > maxPages {
		currentPage = 1
	}

	domains, dbErr := listDomains(int32(pageSize), int32((currentPage-1)*pageSize))

	runTemplate(w, r, "bimi/view.tmpl", map[string]any{
		"CurrentPage": currentPage,
		"Domains":     domains,
		"Err":         dbErr,
		"MaxPage":     maxPages,
		"PageSize":    pageSize,
		"Title":       "View BIMI Logos",
	})

}
