package crawlhub

import (
	"github.com/PuerkitoBio/goquery"
)

type ScrapeResult struct {
	PageTitle     string   `json:"page_title"`
	PrimaryH1     string   `json:"primary_h1"`
	ExtractedInfo []string `json:"extracted_info"`
}

type Parser interface {
	ParsePage(*goquery.Document) ScrapeResult
}
