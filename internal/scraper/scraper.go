package scraper

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type SpiderHandlerFunc[T any] func(*goquery.Document) (*T, error)

func Spider[T any](s http.Client, url string, cb SpiderHandlerFunc[T]) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for URL: %s: %v", url, err)
	}

	// set custom user agent
	// req.Header.Set("User-Agent", "xx")

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse response body from URL %s: %v", url, err)
	}

	result, err := cb(doc)
	if err != nil {
		return nil, fmt.Errorf("callback function for URL %s returned nil, err: %v", url, err)
	}

	return result, nil
}
