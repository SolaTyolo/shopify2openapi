package main

import (
	"fmt"
	"os"

	"github.com/SolaTyolo/shopify2openapi"
)

func main() {
	version := shopify2openapi.VERSION_202404
	if err := shopify2openapi.Scrapper(version); err != nil {
		fmt.Fprintf(os.Stdout, "Failed to scrape shopify docs err: %v", err)
	}
}
