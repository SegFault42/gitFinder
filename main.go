package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

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

func getVulnServer(serverList []string) {
	for _, link := range serverList {
		link = "http://" + link + "/.git/index"

		client := http.Client{
			Timeout: 2 * time.Second,
		}
		resp, _ := client.Get(link)

		if resp != nil && resp.StatusCode == 200 {
			log.Printf("\033[32m%s\033[0m", link)
		} else {
			log.Printf("\033[31m%s\033[0m", link)
		}
	}
}

func main() {
	url := "https://serveur-prive.net/ark-survival-evolved/page/"
	//var serverList []string

	// get all server list
	for i := 1; ; i++ {
		newUrl := url + strconv.Itoa(i)
		doc := dumpPage(newUrl)
		lstUrl := getUrls(doc)
		if len(lstUrl) == 0 {
			break
		}
		getVulnServer(lstUrl)
		//serverList = append(serverList, lstUrl...)
	}

}
