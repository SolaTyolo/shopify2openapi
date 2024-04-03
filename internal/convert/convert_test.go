package convert

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

func TestToJson(t *testing.T) {
	data := openapi3.T{
		OpenAPI: "3.0.0",
	}

	t.Run("EmptyPath", func(t *testing.T) {
		err := ToJson(data, "")
		assert.EqualError(t, err, "path cannot be empty")
	})

	t.Run("RelativePath", func(t *testing.T) {
		err := ToJson(data, "data.json")
		assert.EqualError(t, err, "path must be an absolute path")
	})

	t.Run("NonExistingDirectory", func(t *testing.T) {
		err := ToJson(data, "/path/to/non/existing/directory/data.json")
		assert.EqualError(t, err, "directory does not exist")
	})

	t.Run("SuccessfulConversion", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "data.json")

		err := ToJson(data, filePath)
		assert.NoError(t, err)

		// Verify file existence
		_, err = os.Stat(filePath)
		assert.NoError(t, err)

		// Verify file content
		file, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Equal(t, "{\n    \"info\": null,\n    \"openapi\": \"3.0.0\",\n    \"paths\": null\n}", string(file))
	})
}
