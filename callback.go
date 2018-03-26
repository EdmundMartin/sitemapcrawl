package crawlhub

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func postCallback(callbackUrl string, result ScrapeResult) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(result)
	http.Post(callbackUrl, "application/json; charset=utf-8", b)
}
