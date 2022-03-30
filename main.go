package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type initConfig struct {
	DB_HOST string
	DB_PASS string
}

func main() {
	// Environment variables to be exported - ORG_ID, ENVIRON, PROJECT_NAME, AWS_DEFAULT_REGION, RUNNING_ON_LOCAL(yes, no)
	secretName := "/" + os.Getenv("ORG_ID") + "/" + os.Getenv("ENVIRON") + "/" + os.Getenv("PROJECT_NAME") + "-secrets" // Ex - /myorg/stg/testproject-secrets
	secretRegion := os.Getenv("AWS_DEFAULT_REGION")
	getInitConfig := GetSecret(secretName, secretRegion)

	fmt.Println("Getting secrets from config struct")
	fmt.Println(getInitConfig.DB_HOST)

	fmt.Println("Getting secrets from environment variables")
	fmt.Println(os.Getenv("DB_HOST"))

}

func GetSecret(secretName string, secretRegion string) (secretsMapInitConfig initConfig) {

	versionStage := "AWSCURRENT"

	switch os.Getenv("RUNNING_ON_LOCAL") {

	case "no":
		fmt.Println("Looing for secret - ", secretName)
		secretAwsSession := session.Must(session.NewSession())
		svc := secretsmanager.New(
			secretAwsSession,
			aws.NewConfig().WithRegion(secretRegion),
		)
		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(secretName),
			VersionStage: aws.String(versionStage),
		}
		result, err := svc.GetSecretValue(input)
		if err != nil {
			fmt.Println("Error: Unable to fetch secrets")
			panic(err.Error())
		}
		var secretString string
		if result.SecretString != nil {
			secretString = *result.SecretString

			//Use secrets manager directly
			var secretsMapEnv map[string]interface{}
			json.Unmarshal([]byte(secretString), &secretsMapEnv)

			// Return the struct
			json.Unmarshal([]byte(secretString), &secretsMapInitConfig)

			// pass secret manager key values as environment variable
			for k, v := range secretsMapEnv {
				os.Setenv(k, fmt.Sprint(v))
			}
			return secretsMapInitConfig

		} else {
			fmt.Println("Secret is empty.")
		}

	case "yes":
		fmt.Println("Application running on local. Feed an ENV file using `docker run` or `docker-compose`")

	default:
		fmt.Println("Environment variable RUNNING_ON_LOCAL not set OR it has a wrong value")
	}

	return secretsMapInitConfig
}
