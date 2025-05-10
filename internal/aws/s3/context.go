package s3

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3bucket "github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v3"
)

type S3Resource interface {
	GetObject(input *s3bucket.GetObjectInput) (*s3bucket.GetObjectOutput, error)
}

// S3CloudContext implements CloudContext para S3
type S3CloudContext struct {
	svc    S3Resource
	bucket string
	key    string
}

func NewS3Context(sess *session.Session, bucket, key string) *S3CloudContext {
	return &S3CloudContext{
		svc:    s3bucket.New(sess),
		bucket: bucket,
		key:    key,
	}
}

// GetValue obtém o conteúdo do arquivo S3 e o converte para o formato apropriado
func (ctx *S3CloudContext) GetValue() (interface{}, error) {
	input := &s3bucket.GetObjectInput{
		Bucket: aws.String(ctx.bucket),
		Key:    aws.String(ctx.key),
	}

	result, err := ctx.svc.GetObject(input)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter objeto do S3: %w", err)
	}
	defer result.Body.Close()

	// Ler o conteúdo do arquivo
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler conteúdo do arquivo: %w", err)
	}
	content := string(bodyBytes)

	// Determinar o tipo de arquivo a processar de acordo
	if strings.HasSuffix(ctx.key, ".json") {
		var jsonData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
			return nil, fmt.Errorf("erro ao analisar JSON: %w", err)
		}
		return jsonData, nil

	} else if strings.HasSuffix(ctx.key, ".yaml") || strings.HasSuffix(ctx.key, ".yml") {
		var yamlData map[string]interface{}
		if err := yaml.Unmarshal(bodyBytes, &yamlData); err != nil {
			return nil, fmt.Errorf("erro ao analisar YAML: %w", err)
		}
		return yamlData, nil

	} else if strings.HasSuffix(ctx.key, ".csv") {
		reader := csv.NewReader(strings.NewReader(content))
		records, err := reader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("erro ao analisar CSV: %w", err)
		}

		// Converter CSV para mapa
		if len(records) < 2 {
			return records, nil // Retorna registros crus se não houver cabeçalho
		}

		headers := records[0]
		result := make([]map[string]string, 0, len(records)-1)

		for i := 1; i < len(records); i++ {
			row := make(map[string]string)
			for j := 0; j < len(headers) && j < len(records[i]); j++ {
				row[headers[j]] = records[i][j]
			}
			result = append(result, row)
		}
		return result, nil

	} else {
		// Assumir que é um arquivo de teste simples
		return content, nil
	}
}
