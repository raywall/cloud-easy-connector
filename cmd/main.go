package main

import (
	"encoding/json"
	"time"

	"github.com/raywall/cloud-easy-connector/pkg/auth"
	"github.com/raywall/cloud-easy-connector/pkg/aws"
	"github.com/raywall/cloud-easy-connector/pkg/local"
)

var (
	cloudContextFactory *aws.CloudContextFactory
	tokenManager        auth.AutoManagedToken
	err                 error
)

func init() {
	// inicializa um contexto cloud
	cloudContextFactory, err = aws.NewCloudContextFactory("sa-east-1", "http://localhost:4566")
	if err != nil {
		panic(err)
	}

	// inicializa um contexto de secrets manager
	secretsContext, err := cloudContextFactory.CreateContext(
		aws.SecretsManagerContext,
		map[string]interface{}{
			"secret_name": "my-secrets-manager",
			"secret_type": "json",
		})

	if err != nil {
		panic(err)
	}

	// recupera o valor de um secrets manager
	jsonSecretsValue, err := secretsContext.GetValue()
	if err != nil {
		panic(err)
	}

	authRequest := auth.AuthRequest{}
	err = json.Unmarshal(jsonSecretsValue.([]byte), &authRequest)
	if err != nil {
		panic(err)
	}

	// inicializa um token client auto gerenciado
	tokenManager = auth.NewAutoManagedToken(local.New().GetEnvOrDefault("AUTH_BASE_URL", "https://sts.teste.net/api/oauth/token"), authRequest, false)
	if err = tokenManager.Start(); err != nil {
		panic(err)
	}
}

func main() {
	for i := 0; i < 310; i++ {
		// recupera um token
		token, err := tokenManager.GetToken()
		if err != nil {
			panic(err)
		}

		// imprime o token
		println(token)
		time.Sleep(10 * time.Second)
	}
}
