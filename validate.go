package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AwesomeLogos/bimi-explorer/generated"
)

func fetchImgURL(imgurl string) (string, []byte, error) {

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(imgurl)
	if err != nil {
		return err.Error(), nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	return resp.Header.Get("Content-Type"), body, nil
}

func validate(imgurl string) (bool, string, error) {

	if imgurl == "" {
		return false, "imgurl is empty", nil
	}

	contentType, _, fetchErr := fetchImgURL(imgurl)
	if fetchErr != nil {
		return false, "Unable to fetch image", fetchErr
	}

	if !strings.HasPrefix(contentType, "image/") {
		return false, "Not an image", nil
	}

	//LATER: validate image content

	return true, "OK", nil
}

func bulkValidate() {

	count, countErr := countUnvalidatedDomains()
	if countErr != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to count domains: %v\n", countErr)
		return
	}

	domains, domainErr := listUnvalidatedDomains(int32(count), 0)
	if domainErr != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to list domains: %v\n", domainErr)
		return
	}

	validateChanneler(domains)
}

func validateWorker(domainChan chan generated.Domain, wg *sync.WaitGroup) {
	defer wg.Done()
	for domain := range domainChan {
		valid, msg, _ := validate(domain.Imgurl.String)
		updateValidation(domain.Domain, valid, msg)
	}
}

func validateChanneler(domains []generated.Domain) {

	domainChan := make(chan generated.Domain, len(domains))
	wg := &sync.WaitGroup{}

	maxWorkers, maxErr := strconv.Atoi(os.Getenv("MAXWORKERS"))
	if maxErr != nil {
		maxWorkers = 10 * runtime.NumCPU()
	}
	total := 0
	for _, domain := range domains {
		domainChan <- domain
		total++
	}
	fmt.Fprintf(os.Stderr, "INFO: %d total domains (using %d workers)", total, maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go validateWorker(domainChan, wg)
	}

	close(domainChan)

	previousLen := len(domainChan)
	for len(domainChan) > 0 {
		if previousLen-len(domainChan) > 1000 {
			fmt.Fprintf(os.Stderr, "\nINFO: Domains remaining: %d..", len(domainChan))
			previousLen = len(domainChan)
		}
		fmt.Fprintf(os.Stderr, ".")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Fprintf(os.Stderr, "\nINFO: Total domains: %d\n", total)

	wg.Wait()
}
