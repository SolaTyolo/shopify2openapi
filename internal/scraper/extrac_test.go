package scraper

import (
	"reflect"
	"testing"
)

func TestExtraMetaFromHtml(t *testing.T) {
	htmlBytes := []byte(`<!DOCTYPE html>
<html>
<head>
	<title>Test Page</title>
</head>
<body>
	//<![CDATA[window.RailsData = 
		{
			"features": ["feat-1"],
			"env": "prod"
		}
	//]]>
</body>
</html>`)

	expected := &ShopifyDocsMeta{
		Features: []string{"feat-1"},
		Env:      "prod",
	}

	result, err := extraMetaFromHtml(htmlBytes)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result, got: %v, want: %v", result, expected)
	}
}
