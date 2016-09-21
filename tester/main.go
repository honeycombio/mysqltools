package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/honeycombio/mysqltools/query/normalizer"
)

var astNormalizerTime time.Duration
var scanNormalizerTime time.Duration

var astNormalizerSuccess int64
var scanNormalizerSuccess int64

var astNormalizerFailure int64
var scanNormalizerFailure int64

var (
	scanNormalizer = &normalizer.Scanner{}
	astNormalizer  = &normalizer.Parser{}
)

func testQuery(input string) {
	var now time.Time

	now = time.Now()
	normalized := scanNormalizer.NormalizeQuery(input)
	scanNormalizerTime += time.Since(now)
	if normalized == "" {
		scanNormalizerFailure++
	} else {
		scanNormalizerSuccess++
	}

	now = time.Now()
	normalized = astNormalizer.NormalizeQuery(input)
	astNormalizerTime += time.Since(now)
	if normalized == "" {
		//fmt.Println("ast normalizer failed", input)
		astNormalizerFailure++
	} else {
		astNormalizerSuccess++
	}
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	queryText := ""
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		if len(text) == 0 || strings.HasPrefix(text, "#") {
			continue
		}

		if strings.HasPrefix(text, "--") {
			text = strings.TrimPrefix(text, "--")
			c := strings.Split(text, " ")
			command, args := c[0], c[1:]
			switch command {
			case "echo":
				fmt.Println(strings.Join(args, " "))
				continue
			default:
				fmt.Println("unhandled command: " + text)
				continue
			}
		}

		queryText = queryText + " " + text
		if strings.HasSuffix(queryText, ";") {
			queryText = strings.TrimSuffix(queryText, ";")
			testQuery(queryText)
			queryText = ""
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ast normalizer : %dms for %d queries (%d queries/minute). %d failures\n", astNormalizerTime.Nanoseconds()/1e6, astNormalizerSuccess, int64(float64(astNormalizerSuccess)/astNormalizerTime.Minutes()), astNormalizerFailure)
	fmt.Printf("scan normalizer: %dms for %d queries (%d queries/minute). %d failures\n", scanNormalizerTime.Nanoseconds()/1e6, scanNormalizerSuccess, int64(float64(scanNormalizerSuccess)/scanNormalizerTime.Minutes()), scanNormalizerFailure)
}
