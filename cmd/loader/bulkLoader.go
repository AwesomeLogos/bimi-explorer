package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/AwesomeLogos/bimi-explorer/lib/bimi"
)

var splitter = regexp.MustCompile("[\r\n, ]+")

func bulkWorker(domainChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for domain := range domainChan {
		bimi.LookupBimi(domain)
	}
}

func bulkLoader(filenames []string) {

	domainChan := make(chan string, 50*1024*1024)
	wg := &sync.WaitGroup{}

	maxWorkers, maxErr := strconv.Atoi(os.Getenv("MAXWORKERS"))
	if maxErr != nil {
		maxWorkers = 10 * runtime.NumCPU()
	}
	total := 0
	for _, filename := range filenames {
		fmt.Fprintf(os.Stderr, "INFO: loading %s\n", filename)
		count := 0
		file, openErr := os.Open(filename)
		if openErr != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to open file %s: %v", filename, openErr)
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(bufio.NewReaderSize(file, 100*1024))
		for scanner.Scan() {
			domainChan <- scanner.Text()
			count++
			if count%1000 == 0 {
				fmt.Fprintf(os.Stderr, ".")
			}
		}
		fmt.Fprintf(os.Stderr, "\nINFO: %d loaded from %s\n", count, filename)
		total += count
	}
	fmt.Fprintf(os.Stderr, "INFO: %d total domains\n", total)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go bulkWorker(domainChan, wg)
	}

	close(domainChan)

	previousLen := len(domainChan)
	for len(domainChan) > 500 {
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

func main() {

	if len(os.Args) <= 1 || os.Args[1] == "--help" {
		fmt.Fprintf(os.Stderr, "Usage: %s <file1> [<file2> ...]\n", os.Args[0])
		os.Exit(1)
	} else {
		bulkLoader(os.Args[1:])
	}
}
