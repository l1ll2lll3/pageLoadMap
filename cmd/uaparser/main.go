package main

import (
	"fmt"
	"log"

	"github.com/ua-parser/uap-go/uaparser"
)

func main() {
	// uagent := "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_3; en-us; Silk/1.1.0-80) AppleWebKit/533.16 (KHTML, like Gecko) Version/5.0 Safari/533.16 Silk-Accelerated=true"

	parser, err := uaparser.New("./regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var userAgents = []string{
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36 OPR/38.0.2220.41",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.134 Safari/537.36",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36 OPR/38.0.2220.41",
		"Mozilla/5.0 (Linux; Android 10; SAMSUNG SM-N960N/KSU3FVG4) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/18.0 Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SAMSUNG SM-F936N) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/18.0 Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SAMSUNG SM-F936N/KSU1AVIG) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/18.0 Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SAMSUNG SM-N976N) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/18.0 Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SAMSUNG SM-S908N) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/18.0 Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/18.0 Chrome/99.0.4844.88 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/605.1 NAVER(inapp; search; 1010; 11.16.6; 12MINI)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Whale/2.3.8.3769 Mobile/15E148 Safari/604.1"}

	for i, uagent := range userAgents {

		client := parser.Parse(uagent)

		fmt.Println(i, client.UserAgent.Family) // "Amazon Silk"
		// fmt.Println(client.UserAgent.Major)     // "1"
		// fmt.Println(client.UserAgent.Minor)     // "1"
		// fmt.Println(client.UserAgent.Patch)     // "0-80"
		fmt.Println("os", client.Os.Family) // "Android"
		// fmt.Println(client.Os.Major)            // ""
		// fmt.Println(client.Os.Minor)            // ""
		// fmt.Println(client.Os.Patch)            // ""
		// fmt.Println(client.Os.PatchMinor)       // ""
		fmt.Println("device", client.Device.Family) // "Kindle Fire"
	}
}
