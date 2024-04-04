package main

import (
	"fmt"
	"net/http"
	"strconv"
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

	runTemplate(w, r, "bimi/index.tmpl", map[string]any{
		"H1":    "Welcome",
		"Title": "BIMI Explorer",
	})

}

func rootHandlerPost(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	domain := r.Form.Get("domain")
	if domain == "" {
		http.Redirect(w, r, "/bimi/?err=You+must+enter+a+domain", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/bimi/%s/", domain), http.StatusSeeOther)
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
