package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func bulkHandlerGet(w http.ResponseWriter, r *http.Request) {

	domains := os.Getenv("BULKDOMAINS")

	runTemplate(w, r, "bimi/bulk.tmpl", map[string]any{
		"title":   "Bulk Add",
		"domains": domains,
	})
}

func bulkHandlerPost(w http.ResponseWriter, r *http.Request) {

	formErr := r.ParseForm()
	if formErr != nil {
		logger.Error("failed to read request data", "err", formErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	domains := splitter.Split(r.Form.Get("domains"), -1)
	if len(domains) == 0 {
		http.Redirect(w, r, "bulk.html?err=input+required", http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf8")
	w.Header().Set("X-Content-Type-Options", "nosniff") // FU Chrome

	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}

	errCount := 0

	for index, domain := range domains {
		w.Write([]byte(fmt.Sprintf("Processing %s (%d of %d)...\n", domain, index, len(domains))))

		bimi, lookupErr := lookupBimi(domain)
		if lookupErr != nil {
			errCount++
		} else {
			w.Write([]byte(fmt.Sprintf("BIMI LOGO: %s\n", bimi)))
		}
		flusher.Flush()
		time.Sleep(250 * time.Millisecond)
	}
	w.Write([]byte("Complete!"))
}
