package main

import (
	"fmt"
	"net/http"
)

func bimiHandler(w http.ResponseWriter, r *http.Request) {

	domain := r.PathValue("domain")
	if domain == "" {
		http.Redirect(w, r, "/?err=You+must+enter+a+domain", http.StatusSeeOther)
		return
	}

	bimi, bimiErr := lookupBimi(domain)

	runTemplate(w, r, "_bimi/index.tmpl", map[string]any{
		"title":  fmt.Sprintf("BIMI for %s", domain),
		"domain": domain,
		"err":    bimiErr,
		"bimi":   bimi,
	})
}
