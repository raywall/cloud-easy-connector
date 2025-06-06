package main

import (
	"encoding/json"
	"time"

	"github.com/raywall/cloud-easy-connector/pkg/auth"
	"github.com/raywall/cloud-easy-connector/pkg/cloud"
	"github.com/raywall/cloud-easy-connector/pkg/local"
)

var (
	cloudContext cloud.CloudContext
	err          error
)

func init() {
	// initializes a cloud context
	availableResources := &cloud.CloudContextList{
		cloud.S3Context,
		cloud.SSMContext,
		cloud.SecretsManagerContext,
	}

	cloudContext, err = cloud.NewAwsCloudContext(
		"us-east-1",
		"http://localhost:4566",
		availableResources)

	if err != nil {
		panic(err)
	}

	// recupera o valor de um secrets manager
	jsonSecretsValue, err := cloudContext.GetSecretValue(
		"my-secrets-manager",
		cloud.JSONSecret)

	if err != nil {
		panic(err)
	}

	authRequest := auth.AuthRequest{}
	err = json.Unmarshal(jsonSecretsValue.([]byte), &authRequest)
	if err != nil {
		panic(err)
	}

	// inicializa um token client auto gerenciado
	cloudContext.NewAutoManagedToken(
		local.New().GetEnvOrDefault("AUTH_BASE_URL", "https://sts.teste.net/api/oauth/token"),
		authRequest.ClientID,
		authRequest.ClientSecret,
		false)

	if err = cloudContext.GetAutoManagedToken().Start(); err != nil {
		panic(err)
	}
}

func main() {
	for i := 0; i < 310; i++ {
		// recupera um token
		token, err := cloudContext.GetAutoManagedToken().GetToken()
		if err != nil {
			panic(err)
		}

		// imprime o token
		println(token)
		time.Sleep(10 * time.Second)
	}
}
