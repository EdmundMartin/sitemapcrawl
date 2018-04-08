# sitemapcrawl
Highly concurrent sitemap crawling library written in Golang.
## Example Usage
```golang
package main

import (
	"fmt"

	"github.com/EdmundMartin/sitemapcrawl"
)

func main() {
	p := sitemapcrawl.DefaultParser{}
	results := sitemapcrawl.ScrapeSitemap("http://edmundmartin.com/sitemap.xml", p, 10)
	for _, res := range results {
		fmt.Println(res)
	}
}

```
## TODO
* Some sitemaps don't have sub sitemaps named as XML files. Type check response for more accurate stuff
