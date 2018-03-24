package sitemapSeo

func scrapePage(url string, token chan struct{}, parser Parser) (SeoData, error) {
	res, err := crawlPage(url, token)
	if err != nil {
		return SeoData{}, err
	}
	data, err := parser.GetSeoData(res)
	if err != nil {
		return SeoData{}, err
	}
	return data, nil
}
