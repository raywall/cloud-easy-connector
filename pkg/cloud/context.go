package cloud

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/raywall/cloud-easy-connector/internal/aws/s3"
	"github.com/raywall/cloud-easy-connector/internal/aws/secretsmanager"
	"github.com/raywall/cloud-easy-connector/internal/aws/ssm"
	"github.com/raywall/cloud-easy-connector/pkg/auth"
)

type (
	ContextType      int
	CloudContextType int
	CloudContextList []ContextType
	SecretType       string
)

const (
	AwsCloud CloudContextType = iota
	Azure
	GoogleCloud

	S3Context ContextType = iota
	SSMContext
	SecretsManagerContext

	TextSecret SecretType = "text"
	JSONSecret SecretType = "json"
)

type CloudContextObject struct {
	awsSession           *session.Session
	awsContextCollection map[ContextType]interface{}
	managedToken         auth.AutoManagedToken
}

// CloudContext é a interface principal para interação com recursos AWS
type CloudContext interface {
	GetS3ObjectValue(bucketName, keyName string) (interface{}, error)
	GetParameterValue(parameterName string, withDecryption bool) (interface{}, error)
	GetSecretValue(secretName string, secretType SecretType) (interface{}, error)
	NewAutoManagedToken(url, clientId, clientSecret string, certSkipVerify bool)
	GetAutoManagedToken() auth.AutoManagedToken
}

// NewAwsCloudContext cria um novo contexto de cloud para interação com recursos AWS
func NewAwsCloudContext(region, endpoint string, availableResources *CloudContextList) (CloudContext, error) {
	if region == "" {
		return nil, fmt.Errorf("unsupported region: %s", region)
	}
	if len(*availableResources) == 0 {
		return nil, errors.New("you need to identify the resources that will be used")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao criar sessão AWS: %w", err)
	}
	if endpoint != "" {
		sess.Config.Endpoint = aws.String(endpoint)
	}

	cloudContext := CloudContextObject{
		awsSession:           sess,
		awsContextCollection: make(map[ContextType]interface{}, 0),
	}
	for _, res := range *availableResources {
		switch res {
		case S3Context:
			cloudContext.awsContextCollection[res] = s3.NewS3Context(cloudContext.awsSession)
			continue

		case SSMContext:
			cloudContext.awsContextCollection[res] = ssm.NewSSMContext(cloudContext.awsSession)
			continue

		case SecretsManagerContext:
			cloudContext.awsContextCollection[res] = secretsmanager.NewSecretsManagerContext(cloudContext.awsSession)
			continue

		default:
			return nil, fmt.Errorf("the ContextType was not identified: %v", res)
		}
	}

	return &cloudContext, nil
}

func (c *CloudContextObject) GetS3ObjectValue(bucketName, keyName string) (interface{}, error) {
	if ctx, ok := c.awsContextCollection[S3Context]; ok {
		return (ctx.(*s3.S3CloudContext)).GetValue(bucketName, keyName)
	}
	return nil, errors.New("can't find the available context to s3 resource")
}

func (c *CloudContextObject) GetParameterValue(parameterName string, withDecryption bool) (interface{}, error) {
	if ctx, ok := c.awsContextCollection[SSMContext]; ok {
		return (ctx.(*ssm.SSMCloudContext)).GetValue(parameterName, withDecryption)
	}
	return nil, errors.New("can't find the available secrets manager resource")
}

func (c *CloudContextObject) GetSecretValue(secretName string, secretType SecretType) (interface{}, error) {
	if ctx, ok := c.awsContextCollection[SecretsManagerContext]; ok {
		return (ctx.(*secretsmanager.SecretsManagerCloudContext)).GetValue(secretName, string(secretType))
	}
	return nil, errors.New("can't find the available context to secrets manager resource")
}

func (c *CloudContextObject) NewAutoManagedToken(url, clientId, clientSecret string, certSkipVerify bool) {
	c.managedToken = auth.NewAutoManagedToken(
		url,
		auth.AuthRequest{
			ClientID:     clientId,
			ClientSecret: clientSecret,
		},
		certSkipVerify,
	)
}

func (c *CloudContextObject) GetAutoManagedToken() auth.AutoManagedToken {
	return c.managedToken
}
