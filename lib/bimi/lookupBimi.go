package bimi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"

	"github.com/AwesomeLogos/bimi-explorer/internal/common"
	"github.com/AwesomeLogos/bimi-explorer/internal/db"
)

type DnsResponse struct {
	Status int         `json:"Status"`
	Answer []DnsAnswer `json:"Answer"`
}

type DnsAnswer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

func LookupBimi(rawDomain string) (string, error) {

	domain, domainErr := PurifyDomain(rawDomain)
	if domainErr != nil {
		common.Logger.Error("invalid domain", "domain", rawDomain, "err", domainErr)
		return "", domainErr
	}

	requestURL := fmt.Sprintf("https://cloudflare-dns.com/dns-query?name=default._bimi.%s&type=TXT", domain)
	req, newErr := http.NewRequest(http.MethodGet, requestURL, nil)
	if newErr != nil {
		common.Logger.Error("client could not create request", "err", newErr)
		return "", newErr
	}
	req.Header.Set("accept", "application/dns-json")
	req.Header.Set("user-agent", "bimi-logos/1.0")

	res, httpErr := http.DefaultClient.Do(req)
	if httpErr != nil {
		common.Logger.Error("client error making http request", "err", httpErr)
		return "", httpErr
	}
	common.Logger.Info("lookup status code", "domain", domain, "statuscode", res.StatusCode)

	resBody, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		common.Logger.Error("client could not read response body", "err", readErr, "domain", domain)
		return "", readErr
	}
	common.Logger.Debug("dns lookup success", "result", resBody)

	data := &DnsResponse{}
	jsonErr := json.Unmarshal(resBody, data)
	if jsonErr != nil {
		common.Logger.Error("unable to parse json dns results", "err", jsonErr, "data", string(resBody), "domain", domain)
		return "", jsonErr
	}
	common.Logger.Info("dns response", "data", data)

	if data.Status != 0 {
		common.Logger.Error("dns error status", "statuscode", data.Status, "domain", domain)
		return "", fmt.Errorf("DNS_ERROR: %d", data.Status)
	}
	if data.Answer == nil {
		common.Logger.Error("dns error no answer", "domain", domain)
		return "", fmt.Errorf("DNS_ERROR: %s", "NILDATA")
	}
	if len(data.Answer) == 0 {
		common.Logger.Error("dns error empty answer", "domain", domain)
		return "", fmt.Errorf("DNS_ERROR: %s", "LENZERO")
	}

	for _, answerEntry := range data.Answer {
		answer := removeQuotes(answerEntry.Data)
		bimi := findBimi(answer)
		if bimi != "" {
			db.UpsertDomain(domain, bimi)
			return bimi, nil
		}
	}

	return "", fmt.Errorf("no Bimi in answers")
}

var bimiSplitter = regexp.MustCompile("; *")

func findBimi(data string) string {
	fields := bimiSplitter.Split(data, -1)
	if len(fields) < 2 {
		return ""
	}

	if !slices.Contains(fields, "v=BIMI1") {
		return ""
	}

	for _, field := range fields {
		if len(field) > 10 && field[0:2] == "l=" {
			return field[2:]
		}
	}
	return ""
}

func removeQuotes(original string) string {
	if len(original) < 2 {
		return original
	}
	if original[0] == '"' && original[len(original)-1] == '"' {
		return original[1 : len(original)-1]
	}
	return original
}
