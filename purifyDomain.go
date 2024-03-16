package main

import "strings"

func purifyDomain(domain string) (string, error) {
	//LATER: check against public suffix list
	return strings.ToLower(domain), nil
}
