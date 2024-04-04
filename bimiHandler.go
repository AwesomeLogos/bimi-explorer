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
		"Bimi":   bimi,
		"Domain": domain,
		"Err":    bimiErr,
		"Title":  fmt.Sprintf("BIMI Logo for %s", domain),
	})
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {

	requestedDomain := r.PathValue("domain")
	if requestedDomain == "" {
		http.Error(w, "You must enter a domain", http.StatusBadRequest)
		return
	}

	domain, domainErr := getDomain(requestedDomain)
	if domainErr != nil {
		http.Error(w, "Unable to get domain", http.StatusNotFound)
		return
	}

	bimi, bimiErr := lookupBimi(requestedDomain)

	runTemplate(w, r, "_bimi/index.tmpl", map[string]any{
		"Bimi":   bimi,
		"Domain": domain,
		"Err":    bimiErr,
		"Title":  fmt.Sprintf("BIMI Logo for %s", domain),
	})
}

func sourceHandler(w http.ResponseWriter, r *http.Request) {

	requestedDomain := r.PathValue("domain")
	if requestedDomain == "" {
		http.Redirect(w, r, "/?err=You+must+enter+a+domain", http.StatusSeeOther)
		return
	}

	domain, domainErr := getDomain(requestedDomain)
	if domainErr != nil {
		http.Error(w, "Unable to get domain", http.StatusNotFound)
		return
	}

	contentType, body, fetchErr := fetchImgURL(domain.Imgurl.String)

	runTemplate(w, r, "_bimi/source.tmpl", map[string]any{
		"ContentType": contentType,
		"Domain":      requestedDomain,
		"Err":         fetchErr,
		"Formatted":   body, //LATER: pretty print the XML
		"Raw":         body,
		"Title":       fmt.Sprintf("BIMI Logo for %s", domain),
	})
}
