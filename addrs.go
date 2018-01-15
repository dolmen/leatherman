package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// allAddrs returns all email addresses an email was sent to (To, Cc, and Bcc)
func allAddrs(email *mail.Message, path string) []*mail.Address {
	addrs := []*mail.Address{}
	headers := []string{"To", "Cc", "Bcc"}
	for _, x := range headers {
		if email.Header.Get(x) != "" {
			iAddrs, err := email.Header.AddressList(x)
			// Emails tend to be messy, just move on.
			if err != nil {
				continue
			}
			addrs = append(addrs, iAddrs...)
		}
	}

	return addrs
}

// buildFrecencyMap parses all emails found in the passed glob and returns a map
// of addresses, scored based on how recently they were mailed to.
func buildFrecencyMap(glob string) map[string]float64 {
	matches, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal("couldn't get glob", err)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(matches)))
	now := time.Now()
	lambda := math.Exp(2) / 30

	score := map[string]float64{}

	for _, path := range matches {
		file, err := os.Open(path)
		if err != nil {
			log.Println("Coudln't open email", path, err)
		}
		email, err := mail.ReadMessage(file)
		if err != nil {
			log.Println("Coudln't parse email", path, err)
			continue
		}
		for _, addr := range allAddrs(email, path) {
			if _, ok := score[addr.Address]; ok {
				continue
			}
			time, err := email.Header.Date()
			if err != nil {
				log.Println("Couldn't read date for", file, err)
				continue
			}
			age := now.Sub(time).Hours() / 24

			score[addr.Address] = math.Exp(lambda * age)
		}
	}
	return score
}

// buildAddrMap returns a map of address and content, based on os.Stdin
func buildAddrMap() map[string]string {
	scanner := bufio.NewScanner(os.Stdin)

	ret := map[string]string{}
	for scanner.Scan() {
		z := strings.SplitN(scanner.Text(), "\t", 2)
		if len(z) < 2 {
			continue
		}
		if _, ok := ret[z[0]]; ok {
			continue
		}
		ret[z[0]] = z[1]
	}

	return ret
}

// Addrs sorts the addresses passed on stdin based on how recently they were
// used, based on the glob passed on the arguments.
func Addrs(args []string) {
	if len(args) != 2 {
		log.Fatal("Please pass a glob")
	}
	score := buildFrecencyMap(args[1])
	addrs := buildAddrMap()

	// map of addresses that have been scored
	scored := map[string]string{}
	// keys list, for sorting based on score later
	scoredKeys := []string{}
	for key := range score {
		scored[key] = addrs[key]
		delete(addrs, key)
		scoredKeys = append(scoredKeys, key)
	}

	// sort keys based on score
	sort.Slice(
		sort.StringSlice(scoredKeys),
		func(i, j int) bool { return score[scoredKeys[i]] < score[scoredKeys[j]] },
	)

	// first line is blank
	fmt.Println()
	for _, key := range scoredKeys {
		fmt.Println(key + "\t" + scored[key])
	}

	// sort remaining addrs based on keys
	addrKeys := []string{}
	for key := range addrs {
		addrKeys = append(addrKeys, key)
	}
	sort.Sort(sort.StringSlice(addrKeys))
	for _, key := range addrKeys {
		fmt.Println(key + "\t" + addrs[key])
	}
}