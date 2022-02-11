package main

import (
	"fmt"
	"os"
)

func main() {
	// Environment variables to be exported - ORG_ID, ENVIRON, PROJECT_NAME, AWS_DEFAULT_REGION, RUNNING_ON_LOCAL(yes, no)
        secretName := "/" + os.Getenv("ORG_ID") + "/" + os.Getenv("ENVIRON") + "/" + os.Getenv("PROJECT_NAME") + "-secrets" // Ex - /myorg/stg/testsecret-secrets
	secretRegion := os.Getenv("AWS_DEFAULT_REGION")
	GetSecret(secretName, secretRegion)

	fmt.Println(os.Getenv("TEST"))
}

func GetSecret(secretName string, secretRegion string) {

	versionStage := "AWSCURRENT"

	switch os.Getenv("RUNNING_ON_LOCAL") {

	case "no":
		fmt.Println("Looing for secret - ", secretName)
		svc := secretsmanager.New(
			session.New(),
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
			var secretsMap map[string]interface{}
			json.Unmarshal([]byte(secretString), &secretsMap)

			// pass secret manager key values as environment variable
			for k, v := range secretsMap {
				os.Setenv(k, fmt.Sprint(v))
			}
		} else {
			fmt.Println("Secret is empty.")
		}

	case "yes":
		fmt.Println("Application running on local. Feed an ENV file using `docker run` or `docker-compose`")

	default:
		fmt.Println("Environment variable RUNNING_ON_LOCAL not set OR it has a wrong value")
	}
}
