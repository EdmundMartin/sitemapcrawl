# sitemapSeo
Highly concurrent sitemap crawling library written in Golang.
## Example Usage
```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/EdmundMartin/sitemapSeo"
	"github.com/PuerkitoBio/goquery"
)

type ExampleParser struct {
}

func (d ExampleParser) GetSeoData(resp *http.Response) (sitemapSeo.SeoData, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return sitemapSeo.SeoData{}, err
	}
	result := sitemapSeo.SeoData{}
	result.URL = resp.Request.URL.String()
	result.StatusCode = resp.StatusCode
	result.Title = doc.Find("title").First().Text()
	result.H1 = doc.Find("h1").First().Text()
	return result, nil
}

func main() {
	p := ExampleParser{}
	results := sitemapSeo.ScrapeSitemap("https://cartoonhd.biz/sitemap.xml", p, 10)
	for _, res := range results {
		fmt.Println(res)
	}
}
```
## TODO
* Some sitemaps don't have sub sitemaps named as XML files. Type check response for more accurate stuff
