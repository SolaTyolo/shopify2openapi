package scraper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
)

const (
	DOCS_INDEX = "window.RailsData"
)

func ExtraShopifyOpenApiMeta(version string) SpiderHandlerFunc[ShopifyAPIMeta] {
	return func(tree *goquery.Document) (meta *ShopifyAPIMeta, err error) {
		selection := tree.Find("script")
		if selection == nil {
			return nil, fmt.Errorf("ExtraShopifyOpenApiMeta() return None")
		}
		selection.Each(func(i int, e *goquery.Selection) {
			if e == nil {
				return
			}

			if !strings.Contains(e.Text(), DOCS_INDEX) {
				return
			}
			cdata := strings.ReplaceAll(e.Text(), "\n", "")
			docsMeta, err := extraMetaFromHtml([]byte(cdata))
			if err != nil {
				return
			}

			if docsMeta.API != nil && lo.Contains(docsMeta.API.SelectableVersions, version) {
				meta = docsMeta.API
			}

		})
		if meta == nil {
			err = fmt.Errorf("ExtraShopifyOpenApiMeta() return None")
		}
		return
	}
}

func extraMetaFromHtml(htmlBytes []byte) (*ShopifyDocsMeta, error) {
	cdataStart := fmt.Sprintf("//<![CDATA[%s = ", DOCS_INDEX)
	cdataEnd := "//]]>"

	start := strings.Index(string(htmlBytes), cdataStart)
	if start < 0 {
		return nil, fmt.Errorf("ExtraMetaFromHtml(): Start marker '%s' not found", cdataStart)
	}

	start += len(cdataStart) // Adjust start position by the length of the start marker

	end := strings.Index(string(htmlBytes[start:]), cdataEnd)
	if end < 0 {
		return nil, fmt.Errorf("ExtraMetaFromHtml(): End marker '%s' not found", cdataEnd)
	}

	data := htmlBytes[start : start+end] // Extract the desired content

	var result ShopifyDocsMeta
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("extraMetaFromHtml() err= %s", err)
	}
	return &result, nil
}
