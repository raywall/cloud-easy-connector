package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/raywall/cloud-easy-connector/internal/aws/s3"
	"github.com/raywall/cloud-easy-connector/internal/aws/secretsmanager"
	"github.com/raywall/cloud-easy-connector/internal/aws/ssm"
)

type CloudContextFactory struct {
	awsSession *session.Session
}

// CloudContext é a interface principal para interação com recursos AWS
type CloudContext interface {
	GetValue() (interface{}, error)
}

type ContextType int

const (
	S3Context ContextType = iota
	SSMContext
	SecretsManagerContext
)

// NewCloudContextFactory cria uma nova fábrica de contextos
func NewCloudContextFactory(region, endpoint string) (*CloudContextFactory, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao criar sessão AWS: %w", err)
	}
	if endpoint != "" {
		sess.Config.Endpoint = aws.String(endpoint)
	}
	return &CloudContextFactory{
		awsSession: sess,
	}, nil
}

// CreateContext cria um CloudContext específico com base no tipo e parâmetros
func (f *CloudContextFactory) CreateContext(contextType ContextType, params map[string]interface{}) (CloudContext, error) {
	switch contextType {
	case S3Context:
		return f.createS3Context(params)
	case SSMContext:
		return f.createSSMContext(params)
	case SecretsManagerContext:
		return f.createSecretsManagerContext(params)
	default:
		return nil, errors.New("tipo de contexto não suportado")
	}
}

func (f *CloudContextFactory) createS3Context(params map[string]interface{}) (CloudContext, error) {
	bucket, ok := params["bucket"].(string)
	if !ok || bucket == "" {
		return nil, errors.New("parâmetro 'bucket' é obrigatório para S3Context")
	}

	key, ok := params["key"].(string)
	if !ok || key == "" {
		return nil, errors.New("parâmetro 'key' é obrigatório para S3Context")
	}

	return s3.NewS3Context(f.awsSession, bucket, key), nil
}

func (f *CloudContextFactory) createSSMContext(params map[string]interface{}) (CloudContext, error) {
	paramName, ok := params["parameter_name"].(string)
	if !ok || paramName == "" {
		return nil, errors.New("parâmetro 'parameter_name' é obrigatório para SSMContext")
	}

	return ssm.NewSSMContext(f.awsSession, paramName, true), nil
}

func (f *CloudContextFactory) createSecretsManagerContext(params map[string]interface{}) (CloudContext, error) {
	secretName, ok := params["secret_name"].(string)
	if !ok || secretName == "" {
		return nil, errors.New("parâmetro 'secret_name' é obrigatório para SecretsManagerContext")
	}

	secretType := secretsmanager.TextSecret
	if typeStr, ok := params["secret_type"].(string); ok && typeStr != "" {
		secretType = secretsmanager.SecretType(typeStr)
	}

	return secretsmanager.NewSecretsManagerContext(f.awsSession, secretName, secretType), nil
}
