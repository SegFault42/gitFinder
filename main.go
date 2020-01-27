package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func main() {
	url := "https://serveur-prive.net/minecraft/page/"
	var serverList []string

	for i := 1; ; i++ {
		newUrl := url + strconv.Itoa(i)
		doc := dumpPage(newUrl)
		lstUrl := getUrls(doc)
		if len(lstUrl) == 0 {
			break
		}
		serverList = append(serverList, lstUrl...)
	}

	for nb, elem := range serverList {
		fmt.Printf("%d : %s\n", nb, elem)
	}
}
