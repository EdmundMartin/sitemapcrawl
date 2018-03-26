package sitemapSeo

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SeoData struct {
	URL             string
	Title           string
	H1              string
	MetaDescription string
	StatusCode      int
}

type Parser interface {
	GetSeoData(resp *http.Response) (SeoData, error)
}

func extractUrls(response *http.Response) ([]string, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []string{}
	sel := doc.Find("loc")
	for i := range sel.Nodes {
		loc := sel.Eq(i)
		result := loc.Text()
		results = append(results, result)
	}
	return results, nil
}

func isSitemap(urls []string) ([]string, []string) {
	sitemapFiles := []string{}
	pages := []string{}
	for _, page := range urls {
		foundSitemap := strings.Contains(page, "xml")
		if foundSitemap == true {
			fmt.Println("Found Sitemap", page)
			sitemapFiles = append(sitemapFiles, page)
		} else {
			pages = append(pages, page)
		}
	}
	return sitemapFiles, pages
}

func extractSitemapURLs(startUrl string) []string {
	worklist := make(chan []string)
	toCrawl := []string{}
	var n int
	n++
	go func() { worklist <- []string{startUrl} }()
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			n++
			go func(link string) {
				response, err := makeRequest(link)
				if err != nil {
					log.Printf("Error retrieving URL: %s", link)
				}
				urls, _ := extractUrls(response)
				if err != nil {
					log.Printf("Error extracting document from response, URL: %s", link)
				}
				sitemapFiles, pages := isSitemap(urls)
				if sitemapFiles != nil {
					worklist <- sitemapFiles
				}
				for _, page := range pages {
					toCrawl = append(toCrawl, page)
				}
			}(link)
		}
	}
	return toCrawl
}

func scrapeUrls(urls []string, parser Parser, concurrency int) []SeoData {
	tokens := make(chan struct{}, concurrency)
	var n int
	n++
	worklist := make(chan []string)
	results := []SeoData{}
	go func() { worklist <- urls }()
	for ; n > 0; n-- {
		list := <-worklist
		for _, url := range list {
			if url != "" {
				n++
				go func(url string, token chan struct{}) {
					log.Printf("Requesting URL: %s", url)
					res, err := scrapePage(url, tokens, parser)
					if err != nil {
						log.Printf("Encountered error, URL: %s", url)
					} else {
						results = append(results, res)
					}
					worklist <- []string{}
				}(url, tokens)
			}
		}
	}
	return results
}

func ScrapeSitemap(url string, parser Parser, concurrency int) []SeoData {
	results := extractSitemapURLs(url)
	res := scrapeUrls(results, parser, concurrency)
	return res
}