package ssm

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para SSM
type mockSSMClient struct {
	mock.Mock
}

func (m *mockSSMClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*ssm.GetParameterOutput), args.Error(1)
}

var (
	mockSSM = new(mockSSMClient)
	ctx     = &SSMCloudContext{
		svc:            mockSSM,
		paramName:      "/test/param",
		withDecryption: true,
	}
)

func loadDefaultVariables() {
	mockSSM = new(mockSSMClient)
}

func TestSSMCloudContext_GetValue(t *testing.T) {
	t.Run("Get parameter content", func(t *testing.T) {
		loadDefaultVariables()

		// Preparar mock para SSM
		paramValue := "test-parameter-value"
		mockSSM.On("GetParameter", mock.Anything).Return(&ssm.GetParameterOutput{
			Parameter: &ssm.Parameter{
				Value: aws.String(paramValue),
			},
		}, nil)

		// executar GetValue
		result, err := ctx.GetValue()

		// Verificar resultados
		assert.NoError(t, err)
		assert.Equal(t, paramValue, result)
	})
}
