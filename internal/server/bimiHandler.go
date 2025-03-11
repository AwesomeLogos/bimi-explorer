package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/AwesomeLogos/bimi-explorer/internal/db"
	"github.com/AwesomeLogos/bimi-explorer/lib/bimi"
	"github.com/AwesomeLogos/bimi-explorer/ui"
)

func BimiHandler(w http.ResponseWriter, r *http.Request) {

	rawDomain := r.PathValue("domain")
	if rawDomain == "" {
		http.Redirect(w, r, "/?err=You+must+enter+a+domain", http.StatusSeeOther)
		return
	}

	domain, domainErr := bimi.PurifyDomain(rawDomain)
	if domainErr != nil {
		http.Redirect(w, r, "/?err=Invalid+domain", http.StatusSeeOther)
		return
	}

	if domain != rawDomain {
		http.Redirect(w, r, fmt.Sprintf("/bimi/%s/", domain), http.StatusSeeOther)
		return
	}

	bimi, bimiErr := bimi.LookupBimi(domain)

	var proto = "https"
	if r.Host == "localhost" || r.Host[:10] == "localhost:" {
		proto = "http"
	}

	ui.RunTemplate(w, r, "_bimi/index.tmpl", map[string]any{
		"Bimi":       bimi,
		"CurrentURL": fmt.Sprintf("%s://%s%s", proto, r.Host, r.URL.RequestURI()),
		"Domain":     domain,
		"H1":         template.HTML(fmt.Sprintf("BIMI Logo for <a href=\"https://www.%s/\">%s</a>", domain, domain)),
		"Err":        bimiErr,
		"Title":      fmt.Sprintf("BIMI Logo for %s", domain),
	})
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {

	requestedDomain := r.PathValue("domain")
	if requestedDomain == "" {
		http.Error(w, "You must enter a domain", http.StatusBadRequest)
		return
	}

	domain, domainErr := db.GetDomain(requestedDomain)
	if domainErr != nil {
		http.Error(w, "Unable to get domain", http.StatusNotFound)
		return
	}

	bimi, bimiErr := bimi.LookupBimi(requestedDomain)

	ui.RunTemplate(w, r, "_bimi/index.tmpl", map[string]any{
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

	domain, domainErr := db.GetDomain(requestedDomain)
	if domainErr != nil {
		http.Error(w, "Unable to get domain", http.StatusNotFound)
		return
	}

	contentType, body, fetchErr := bimi.FetchImgURL(domain.Imgurl.String)

	ui.RunTemplate(w, r, "_bimi/source.tmpl", map[string]any{
		"ContentType": contentType,
		"Domain":      requestedDomain,
		"Err":         fetchErr,
		"Formatted":   body, //LATER: pretty print the XML
		"Raw":         body,
		"Title":       fmt.Sprintf("BIMI Logo for %s", domain),
	})
}
