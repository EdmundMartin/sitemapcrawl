package crawlhub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getRequest(targetUrl string) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", targetUrl, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	return res, nil
}

func discoverLinks(doc *goquery.Document, baseURL string) []string {
	foundUrls := []string{}
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			res, _ := s.Attr("href")
			foundUrls = append(foundUrls, res)
		})
		return foundUrls
	} else {
		return []string{}
	}
}

func checkRelative(href string, baseUrl string) string {
	if strings.HasPrefix(href, "/") {
		return fmt.Sprintf("%s%s", baseUrl, href)
	}
	return href
}

func resolveRelativeLinks(href string, baseUrl string) (bool, string) {
	resultHref := checkRelative(href, baseUrl)
	baseParse, _ := url.Parse(baseUrl)
	resultParse, _ := url.Parse(resultHref)
	if baseParse != nil && resultParse != nil {
		if baseParse.Host == resultParse.Host {
			return true, resultHref
		} else {
			return false, ""
		}
	}
	return false, ""
}

func crawlPage(targetUrl string, baseUrl string, parser Parser, token chan struct{}) ([]string, ScrapeResult) {
	fmt.Println(targetUrl)
	token <- struct{}{}
	resp, _ := getRequest(targetUrl)
	<-token
	doc, _ := goquery.NewDocumentFromResponse(resp)
	pageResults := parser.ParsePage(doc)
	links := discoverLinks(doc, baseUrl)
	foundUrls := []string{}
	for _, link := range links {
		ok, correctLink := resolveRelativeLinks(link, baseUrl)
		if ok {
			if correctLink != "" {
				foundUrls = append(foundUrls, correctLink)
			}
		}
	}
	return foundUrls, pageResults
}

func StandardCrawl(baseDomain, startUrl, callback string, parser Parser, concurrency int) {
	worklist := make(chan []string)
	var n int
	n++
	var tokens = make(chan struct{}, concurrency)
	go func() { worklist <- []string{baseDomain} }()
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string, baseDomain, callback string, parser Parser, token chan struct{}) {
					foundLinks, pageResults := crawlPage(link, baseDomain, parser, token)
					fmt.Println(pageResults)
					go func(callbackURL string, results ScrapeResult) {
						postCallback(callbackURL, results)
					}(callback, pageResults)
					if foundLinks != nil {
						worklist <- foundLinks
					}
				}(link, baseDomain, callback, parser, tokens)
			}
		}
	}
}
