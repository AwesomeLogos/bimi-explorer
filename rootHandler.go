package main

import (
	"fmt"
	"net/http"
)

func rootHandlerGet(w http.ResponseWriter, r *http.Request) {

	domains, dbErr := listSampleDomains(50)

	runTemplate(w, r, "index.hbs", map[string]any{
		"domains": domains,
		"err":     dbErr,
		"h1":      "Welcome",
		"title":   "BIMI Explorer",
	})

}

func rootHandlerPost(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	domain := r.Form.Get("domain")
	if domain == "" {
		http.Redirect(w, r, "/?err=You+must+enter+a+domain", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/bimi/%s/", domain), http.StatusSeeOther)
}
