package s3

import (
	"bytes"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockS3Client struct {
	mock.Mock
}

func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

var (
	mockS3 *mockS3Client
	ctx    *S3CloudContext
)

func loadDefaultVariables() {
	mockS3 = new(mockS3Client)
	ctx = &S3CloudContext{
		svc: mockS3,
	}
}

func TestS3CloudContext_GetValue(t *testing.T) {
	t.Run("Get JSON object", func(t *testing.T) {
		loadDefaultVariables()

		// Preparar mock para S3 com arquivo JSON
		jsonContent := `{"name": "test", "value": 123}`
		mockS3.On("GetObject", mock.Anything).Return(&s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewReader([]byte(jsonContent))),
		}, nil)

		// Executar GetValue
		result, err := ctx.GetValue("test-bucket", "test-file.json")

		// Verificar resultados
		assert.NoError(t, err)

		jsonResult, ok := result.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "test", jsonResult["name"])
		assert.Equal(t, float64(123), jsonResult["value"])
	})

	t.Run("Get YAML object", func(t *testing.T) {
		loadDefaultVariables()

		// Preparar mock para S3 com arquivo YAML
		yamlContent := `name: test
value: 123
list:
  - item1
  - item2`

		mockS3.On("GetObject", mock.Anything).Return(&s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewReader([]byte(yamlContent))),
		}, nil)

		// Executar GetValue
		result, err := ctx.GetValue("test-bucket", "test-file.yaml")

		// Verificar resultados
		assert.NoError(t, err)

		yamlResult, ok := result.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "test", yamlResult["name"])
		assert.Equal(t, 123, yamlResult["value"])
	})

	t.Run("Get CSV object", func(t *testing.T) {
		loadDefaultVariables()

		// Preparar mock para S3 com arquivo CSV
		csvContent := `id,name,age
1,Alice,30
2,Bob,25`

		mockS3.On("GetObject", mock.Anything).Return(&s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewReader([]byte(csvContent))),
		}, nil)

		// Executar GetValue
		result, err := ctx.GetValue("test-bucket", "test-file.csv")

		// Verificar resultados
		assert.NoError(t, err)

		csvResult, ok := result.([]map[string]string)
		assert.True(t, ok)
		assert.Len(t, csvResult, 2)
		assert.Equal(t, "Alice", csvResult[0]["name"])
		assert.Equal(t, "30", csvResult[0]["age"])
		assert.Equal(t, "Bob", csvResult[1]["name"])
		assert.Equal(t, "25", csvResult[1]["age"])
	})

	t.Run("Get Text object", func(t *testing.T) {
		loadDefaultVariables()

		// Preparar mock para S3 com arquivo de texto
		textContent := "This is a plain text file."
		mockS3.On("GetObject", mock.Anything).Return(&s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewReader([]byte(textContent))),
		}, nil)

		// Executar GetValue
		result, err := ctx.GetValue("test-bucket", "test-file.txt")

		// Verificar resultados
		assert.NoError(t, err)

		textResult, ok := result.(string)
		assert.True(t, ok)
		assert.Equal(t, textContent, textResult)
	})
}
