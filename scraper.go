package shopify2openapi

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SolaTyolo/shopify2openapi/internal/convert"
	"github.com/SolaTyolo/shopify2openapi/internal/scraper"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
)

const (
	REDOC_EXTENSION_TAG_GROUP    = "x-tagGroups"
	REDOC_EXTENSION_CODE_SAMPLES = "x-codeSamples"

	DEFAULT_FILE_NAME = "openapi.json"
)

func getPath() string {
	dir, _ := os.Getwd()
	return fmt.Sprintf("%s/%s", dir, DEFAULT_FILE_NAME)
}

// scrapper admin-rest docs
func Scrapper(version string) error {
	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}
	apiMeta, err := scraper.Spider(httpClient, GetAdminRestAddr(), scraper.ExtraShopifyOpenApiMeta(version))
	if err != nil {
		return err
	}

	subApis := []string{}
	xTagGroups := make(map[string][]string, 0)

	lo.ForEach(apiMeta.RestSideNav, func(v *scraper.ShopifyRestSideNav, index int) {
		if v == nil {
			return
		}
		if _, exist := xTagGroups[v.Key]; !exist {
			xTagGroups[v.Key] = make([]string, 0)
		}
		for _, child := range v.Children {
			subApis = append(subApis, child.Key)
			xTagGroups[v.Key] = append(xTagGroups[v.Key], child.Label)
		}
	})

	if len(subApis) == 0 {
		return fmt.Errorf("no sub apis found")
	}

	docs := lop.Map(subApis, func(sub string, _ int) *scraper.ShopifyOpenAPISpec {
		api_docs, err := scraper.Spider(httpClient, GetRestfulAPI(version, sub), scraper.ExtraShopifyOpenApiMeta(version))
		if err != nil {
			return nil
		}
		return api_docs.RestResource
	})

	landingDescription, _ := convert.HTML2Markdown(apiMeta.LandingPageData.Description)
	shopifyDoc := NewShopifyOpenApiDoc(
		apiMeta.LandingPageData.Title,
		landingDescription,
		version,
		map[string]interface{}{
			REDOC_EXTENSION_TAG_GROUP: lo.MapToSlice(
				xTagGroups,
				func(name string, tags []string) map[string]interface{} {
					return map[string]interface{}{
						"name": name,
						"tags": tags,
					}
				},
			),
		},
	)

	lo.ForEach(docs, func(doc *scraper.ShopifyOpenAPISpec, _ int) {
		if doc == nil {
			return
		}
		shopifyDoc.AddShopifyRestApiDoc(doc)
	})

	shopifyDoc.GenerateModelTag()
	if err := shopifyDoc.OutputToJson(getPath()); err != nil {
		return err
	}

	return nil
}
