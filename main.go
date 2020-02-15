package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func dumpPage(page string) *goquery.Document {
	resp, err := http.Get(page)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func getUrls(doc *goquery.Document) []string {
	var foundUrls []string
	if doc != nil {
		doc.Find("input").Each(func(i int, s *goquery.Selection) {
			res, exist := s.Attr("value")
			if exist == true {
				foundUrls = append(foundUrls, res)
			}
		})
	}
	return (foundUrls)
}

func request(url string) int {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return -1
	}

	if resp.StatusCode == 200 {
		htmlData, err := ioutil.ReadAll(resp.Body) //<--- here!
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !strings.HasPrefix(string(htmlData), "[core]\n") {
			return -1
		}
	}
	return resp.StatusCode
}

func getVulnServer(server string) {
	tab := [...]string{"http://", "https://", "http://www.", "https://www."}

	for _, elem := range tab {
		endpoint := elem + server + "/.git/config"
		if request(endpoint) != 200 {
			l := log.New(os.Stderr, "", 0)
			l.Printf("\033[31m%s\033[0m", endpoint)
		} else {
			log.Printf("\033[32m%s\033[0m", endpoint)
			fmt.Printf("%s\n", endpoint)
			break
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage : ./prog file_list")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		getVulnServer(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
