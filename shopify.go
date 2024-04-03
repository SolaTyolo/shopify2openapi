package shopify2openapi

import (
	"fmt"
	"strings"

	"github.com/SolaTyolo/shopify2openapi/internal/convert"
	"github.com/SolaTyolo/shopify2openapi/internal/scraper"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/samber/lo"
)

const (
	AdminRestAddr = "https://shopify.dev/docs/api"

	VERSION_202401 = "2024-01"
	VERSION_202404 = "2024-04"
)

func GetAdminRestAddr() string {
	return AdminRestAddr + "/admin-rest"
}

func GetRestfulAPI(version string, resource string) string {
	return fmt.Sprintf("%s/%s/resources/%s", GetAdminRestAddr(), version, resource)
}

type ShopifyOpenApiDoc struct {
	openapi3.T
}

func NewShopifyOpenApiDoc(title, description, version string, extensions map[string]interface{}) ShopifyOpenApiDoc {
	return ShopifyOpenApiDoc{
		openapi3.T{
			OpenAPI: "3.0.0",
			Info: &openapi3.Info{
				Title:       title,
				Description: description,
				Version:     version,
			},
			Servers: []*openapi3.Server{
				{
					URL:         "https://{shop}.myshopify.com/",
					Description: "Default server",
				},
			},
			Tags:  make([]*openapi3.Tag, 0),
			Paths: openapi3.NewPathsWithCapacity(0),
			Components: &openapi3.Components{
				Schemas: openapi3.Schemas{},
			},
			Extensions: extensions,
		},
	}
}

func (shopifyDoc *ShopifyOpenApiDoc) AddShopifyRestApiDoc(doc *scraper.ShopifyOpenAPISpec) {
	if doc == nil {
		return
	}
	// append tag to path
	tag := doc.Info.Title
	desc, _ := convert.HTML2Markdown(doc.Info.Description)

	shopifyDoc.Tags = append(shopifyDoc.Tags, &openapi3.Tag{
		Name:        tag,
		Description: desc,
	})

	lo.ForEach(doc.Paths, func(op *openapi3.Operation, _ int) {
		_ = shopifyDoc.AppendPath(op, tag)
	})

	lo.ForEach(doc.Components, func(c *scraper.ShopifyOpenapiComponent, _ int) {
		_ = shopifyDoc.AppendComponent(c)
	})

	return
}

func (shopifyDoc *ShopifyOpenApiDoc) AppendPath(op *openapi3.Operation, tag string) error {
	if op == nil {
		return nil
	}

	op.Tags = []string{tag}
	codeSamples := make([]map[string]interface{}, 0)

	if examples, ok := op.Extensions[scraper.PATH_EXTENSION_EXAMPLES].([]interface{}); ok {
		if example, ok := examples[0].(map[string]interface{})["codeSamples"].([]interface{}); ok {
			for _, e := range example {
				if e, ok := e.(map[string]interface{}); ok {
					codeSamples = append(codeSamples, map[string]interface{}{
						"lang":   e["language"],
						"label":  e["language"],
						"source": e["example_code"],
					})
				}
			}
		}
	}

	if len(codeSamples) > 0 {
		op.Extensions[REDOC_EXTENSION_CODE_SAMPLES] = codeSamples
	}

	path, ok := op.Extensions[scraper.PATH_EXTENSION_URL].(string)
	if !ok {
		return fmt.Errorf("path is not string")
	}
	action, ok := op.Extensions[scraper.PATH_EXTENSION_ACTION].(string)
	if !ok {
		return fmt.Errorf("action is not string")
	}

	pathItem := new(openapi3.PathItem)
	pathItem.SetOperation(strings.ToUpper(action), op)

	openapi3.WithPath(path, pathItem)(shopifyDoc.Paths)
	return nil
}

func (shopifyDoc *ShopifyOpenApiDoc) AppendComponent(c *scraper.ShopifyOpenapiComponent) error {
	if c == nil {
		return nil
	}

	// refactor properties
	properties := openapi3.Schemas{}
	for _, prop := range c.Properties {
		if prop == nil {
			continue
		}
		name, ok := prop.Extensions[scraper.COMPONENT_EXTENSION_NAME].(string)
		if !ok {
			continue
		}
		prop.Description, _ = convert.HTML2Markdown(prop.Description)
		properties[name] = openapi3.NewSchemaRef("", prop)
	}

	// apiDocs.Components.
	shopifyDoc.Components.Schemas[c.Name] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:       c.Type,
		Title:      c.Title,
		Properties: properties,
		Required:   c.Required,
	})
	return nil

}

// add mode tags and x-taggroup from component
func (shopifyDoc *ShopifyOpenApiDoc) GenerateModelTag() {
	xtagGroup := shopifyDoc.Extensions[REDOC_EXTENSION_TAG_GROUP].([]map[string]interface{})
	tags := shopifyDoc.Tags

	components := lo.Keys(shopifyDoc.Components.Schemas)
	tags = append(tags, lo.Map(components, func(name string, _ int) *openapi3.Tag {
		return &openapi3.Tag{
			Name: ModelName(name),
			Extensions: map[string]interface{}{
				"x-displayName": "The " + name + " Model",
			},
			Description: fmt.Sprintf(`<SchemaDefinition schemaRef="#/components/schemas/%s" />`, name),
		}
	})...)

	xtagGroup = append(xtagGroup, map[string]interface{}{
		"name": "Models",
		"tags": lo.Map(components, func(name string, _ int) string {
			return ModelName(name)
		}),
	})

	shopifyDoc.Extensions[REDOC_EXTENSION_TAG_GROUP] = xtagGroup
	shopifyDoc.Tags = tags
}

func ModelName(name string) string {
	return fmt.Sprintf("%s_model", name)
}

func (shopifyDoc *ShopifyOpenApiDoc) OutputToJson(filepath string) error {
	return convert.ToJson(shopifyDoc.T, filepath)
}
