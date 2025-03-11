package bimi

import (
	"strings"

	"golang.org/x/net/idna"
)

func PurifyDomain(rawDomain string) (string, error) {

	if strings.HasPrefix(rawDomain, "http://") || strings.HasPrefix(rawDomain, "https://") {
		rawDomain = strings.TrimPrefix(rawDomain, "http://")
		rawDomain = strings.TrimPrefix(rawDomain, "https://")
	}

	if strings.HasPrefix("www.", rawDomain) {
		rawDomain = strings.TrimPrefix(rawDomain, "www.")
	}

	if !isASCII(rawDomain) {
		// Punycode the domain
		punycode, err := idna.ToASCII(rawDomain)
		if err != nil {
			return "", err
		}
		rawDomain = punycode
	}

	domain := strings.ToLower(rawDomain)

	//LATER: check against public suffix list
	return domain, nil
}

func isASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}
