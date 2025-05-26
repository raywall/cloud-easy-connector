package secretsmanager

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	sm "github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretType string

const (
	TextSecret SecretType = "text"
	JSONSecret SecretType = "json"
)

type SecretsManagerResource interface {
	GetSecretValue(input *sm.GetSecretValueInput) (*sm.GetSecretValueOutput, error)
}

// SecretsManagerCloudContext implementa CloudContext para Secrets Manager
type SecretsManagerCloudContext struct {
	svc        SecretsManagerResource
	secretName string
	secretType SecretType
}

func NewSecretsManagerContext(sess *session.Session, secretName string, secretType SecretType) *SecretsManagerCloudContext {
	return &SecretsManagerCloudContext{
		svc:        sm.New(sess),
		secretName: secretName,
		secretType: secretType,
	}
}

// GetValue obtém e processa o segredo do Secrets Manager
func (ctx *SecretsManagerCloudContext) GetValue() (interface{}, error) {
	input := &sm.GetSecretValueInput{
		SecretId: aws.String(ctx.secretName),
	}

	result, err := ctx.svc.GetSecretValue(input)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter segredo: %w", err)
	}

	var secretValue string
	if result.SecretString != nil {
		secretValue = *result.SecretString
	} else {
		return nil, errors.New("segredo binário não é suportado")
	}

	switch ctx.secretType {
	case TextSecret:
		return secretValue, nil
	case JSONSecret:
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(secretValue), &jsonData); err != nil {
			return nil, fmt.Errorf("erro ao analisar JSON do segredo: %w", err)
		}
		return []byte(secretValue), nil
	default:
		return secretValue, nil
	}
}
