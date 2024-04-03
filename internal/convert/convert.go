package convert

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/getkin/kin-openapi/openapi3"
)

// ToJson converts the given data of type openapi3.T to JSON format and writes it to the specified file path.
func ToJson(data openapi3.T, path string) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}

	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		return errors.New("path must be an absolute path")
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(cleanPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return errors.New("directory does not exist")
	}
	if err = os.WriteFile(path, jsonData, 0600); err != nil {
		return err
	}
	return nil
}

// converts HTML content into Markdown format.
func HTML2Markdown(htmlContent string) (string, error) {
	converter := md.NewConverter("", true, nil)
	return converter.ConvertString(htmlContent)
}
