package ssm

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ssmParam "github.com/aws/aws-sdk-go/service/ssm"
)

type SSMResource interface {
	GetParameter(input *ssmParam.GetParameterInput) (*ssmParam.GetParameterOutput, error)
}

// SSMCloudContext implementa CloudContext para SSM Parameter Store
type SSMCloudContext struct {
	svc            SSMResource
	paramName      string
	withDecryption bool
}

func NewSSMContext(sess *session.Session, paramName string, withDecryption bool) *SSMCloudContext {
	return &SSMCloudContext{
		svc:            ssmParam.New(sess),
		paramName:      paramName,
		withDecryption: withDecryption,
	}
}

// GetValue obtém o valor do parâmetro SSM
func (ctx *SSMCloudContext) GetValue() (interface{}, error) {
	input := &ssmParam.GetParameterInput{
		Name:           aws.String(ctx.paramName),
		WithDecryption: aws.Bool(ctx.withDecryption),
	}

	result, err := ctx.svc.GetParameter(input)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter parâmetros SSM: %w", err)
	}
	return *result.Parameter.Value, nil
}
