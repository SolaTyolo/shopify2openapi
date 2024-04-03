package scraper

import "github.com/getkin/kin-openapi/openapi3"

const (
	// info extension key
	INFO_EXTENSION_OWNER = "x-owner"

	// component extension key
	COMPONENT_EXTENSION_NAME = "name"

	// path extension key
	PATH_EXTENSION_URL      = "url"
	PATH_EXTENSION_ACTION   = "action"
	PATH_EXTENSION_EXAMPLES = "x-examples"
)

type ShopifyDocsMeta struct {
	Features []string        `json:"features"`
	Env      string          `json:"env"`
	API      *ShopifyAPIMeta `json:"api"`
}

type ShopifyAPIMeta struct {
	BasePath             string                `json:"base_path"`
	RestSideNav          []*ShopifyRestSideNav `json:"rest_sidenav,omitempty"`
	RestResource         *ShopifyOpenAPISpec   `json:"rest_resource,omitempty"`
	SelectableVersions   []string              `json:"selectable_versions"`
	CurrentStableVersion string                `json:"current_stable_version"`
	LandingPage          string                `json:"landing_page"`
	LandingPageId        string                `json:"landing_page_id"`
	LandingPageData      LandingPageData       `json:"landing_page_data"`
}

type ShopifyRestSideNav struct {
	Key      string               `json:"key"`
	Label    string               `json:"label"`
	Children []ShopifyRestSideNav `json:"children,omitempty"`
}

type LandingPageData struct {
	Title               string `json:"title"`
	Description         string `json:"description"`
	ResourceUnsupported string `json:"resource_unsupported"`
}

type ShopifyOpenAPISpec struct {
	Info         openapi3.Info              `json:"info"`
	XShopifyMeta XShopifyMeta               `json:"x-shopify-meta"`
	OpenAPI      string                     `json:"openapi" yaml:"openapi"` // Required
	Components   []*ShopifyOpenapiComponent `json:"components,omitempty" yaml:"components,omitempty"`
	Paths        []*openapi3.Operation      `json:"paths" yaml:"paths"` // Required
}

type ShopifyOpenapiComponent struct {
	Name       string             `json:"name"`
	Properties []*openapi3.Schema `json:"properties,omitempty"`
	Required   []string           `json:"required,omitempty" yaml:"required,omitempty"`
	Type       string             `json:"type,omitempty" yaml:"type,omitempty"`
	Title      string             `json:"title,omitempty" yaml:"title,omitempty"`
}

type XShopifyMeta struct {
	ApiVersioning   bool   `json:"api_versioning"`
	Filename        string `json:"filename"`
	Gid             string `json:"gid"`
	Glossary        bool   `json:"glossary"`
	Hidden          bool   `json:"hidden"`
	MetaDescription string `json:"meta_description"`
	PostmanGroup    string `json:"postman_group"`
}
