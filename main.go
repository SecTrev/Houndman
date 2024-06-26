package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func main() {
	fmt.Println("01001000 01101111 01110101 01101110 01100100 01101101 01100001 01101110")
	//Created with love -Trevor

	wordlist := flag.String("w", "", "Path to the wordlist")
	url := flag.String("u", "", "URL to scan")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	usage := func() {
		fmt.Println("houndman -w wordlist.txt -u http://example.com \nFor further help visit https://github.com/SecTrev/Houndman")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	if *wordlist == "" || *url == "" {
		fmt.Println("houndman -w wordlist.txt -u http://example.com \nFor further help visit https://github.com/SecTrev/Houndman")
		os.Exit(1)
	}

	fmt.Printf("Wordlist: %s\n", *wordlist)
	fmt.Printf("URL: %s\n", *url)
	// Further processing here

	file, err := os.Open(*wordlist)
	if err != nil {
		fmt.Printf("Failed to open wordlist: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := scanner.Text()
		fullUrl := fmt.Sprintf("%s/%s", *url, path)
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error requesting %s: %v\n", url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				fmt.Printf("Found: %s\n", url)
			}
		}(fullUrl)
	}
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading wordlist: %v\n", err)
	}

	if *help || *wordlist == "" || *url == "" {
		usage()
		os.Exit(1)
	}
}
